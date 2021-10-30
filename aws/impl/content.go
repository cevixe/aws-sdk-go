package impl

import (
	"context"
	"github.com/cevixe/aws-sdk-go/aws/model"
	"github.com/cevixe/aws-sdk-go/aws/serdes/gzip"
	"github.com/cevixe/aws-sdk-go/aws/serdes/json"
	"github.com/pkg/errors"
)

func GetContent(
	ctx context.Context,
	contentLocation string,
	contentEncoding string,
	content []byte) []byte {

	data := content
	if contentLocation != "" {
		awsContext := ctx.Value(AwsContext).(*Context)
		store := awsContext.AwsObjectStore
		data = store.GetRawObject(ctx, contentLocation)
	}

	if contentEncoding == "gzip" {
		return gzip.Decompress(data)
	}
	if contentEncoding == "identity" || contentEncoding == "" {
		return data
	}

	panic(errors.Errorf("content encoding not supported: %s", contentEncoding))
}

func GetEventContent(
	ctx context.Context,
	contentLocation string,
	contentEncoding string,
	contentType string,
	content []byte) *model.AwsEventRecord {

	data := GetContent(ctx, contentLocation, contentEncoding, content)

	if contentType == "application/json" || contentType == "" {
		record := &model.AwsEventRecord{}
		json.Unmarshall(data, record)
		return record
	}

	panic(errors.Errorf("content type not supported: %s", contentType))
}
