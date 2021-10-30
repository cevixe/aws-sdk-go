package impl

import (
	"context"
	"fmt"
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
		fmt.Printf("gzip data:%s\n", data)
		return gzip.Decompress(data)
	}
	if contentEncoding == "identity" || contentEncoding == "" {
		fmt.Printf("identity data:%s\n", data)
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
		fmt.Println(json.Marshall(content))
		return record
	}

	panic(errors.Errorf("content type not supported: %s", contentType))
}
