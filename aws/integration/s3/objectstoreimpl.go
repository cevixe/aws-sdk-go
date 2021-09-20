package s3

import (
	"bytes"
	"context"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/aws/aws-sdk-go/service/s3/s3iface"
	"github.com/cevixe/aws-sdk-go/aws/env"
	"github.com/cevixe/aws-sdk-go/aws/factory"
	"github.com/cevixe/aws-sdk-go/aws/model"
	util2 "github.com/cevixe/aws-sdk-go/aws/util"
	"github.com/pkg/errors"
	"io/ioutil"
	"os"
)

type objectStoreImpl struct {
	objectStoreBucket string
	s3Client          s3iface.S3API
}

func NewS3ObjectStore(
	objectStoreBucket string,
	s3Client s3iface.S3API) model.AwsObjectStore {

	return &objectStoreImpl{
		objectStoreBucket: objectStoreBucket,
		s3Client:          s3Client,
	}
}

func NewDefaultS3ObjectStore(awsFactory factory.AwsFactory) model.AwsObjectStore {
	objectStoreBucketName := os.Getenv(env.CevixeObjectStoreBucketName)
	s3Client := awsFactory.S3Client()
	return NewS3ObjectStore(objectStoreBucketName, s3Client)
}

func (o objectStoreImpl) GetRawObject(ctx context.Context, id string) []byte {
	return o.getObject(ctx, id)
}

func (o objectStoreImpl) GetJsonObject(ctx context.Context, id string, v interface{}) {
	buffer := o.getObject(ctx, id)
	util2.UnmarshalJson(buffer, v)
}

func (o objectStoreImpl) getObject(ctx context.Context, id string) []byte {

	input := &s3.GetObjectInput{
		Bucket: aws.String(o.objectStoreBucket),
		Key:    aws.String("object/" + id),
	}

	output, err := o.s3Client.GetObjectWithContext(ctx, input)

	if err != nil {
		panic(errors.Wrapf(err, "cannot get reference(%v) from S3", id))
	}

	buffer, err := ioutil.ReadAll(output.Body)
	if err != nil {
		panic(errors.Wrapf(err, "cannot unmarshal object(%v) from S3", id))
	}

	return buffer
}

func (o objectStoreImpl) SaveRawObject(ctx context.Context, id string, buffer []byte) {
	o.saveObject(ctx, id, "binary/octet-stream", buffer)
}

func (o objectStoreImpl) SaveJsonObject(ctx context.Context, id string, v interface{}) {
	buffer := util2.MarshalJson(v)
	o.saveObject(ctx, id, "application/json", buffer)
}

func (o objectStoreImpl) saveObject(ctx context.Context, id string, contentType string, buffer []byte) {

	reader := bytes.NewReader(buffer)

	input := &s3.PutObjectInput{
		Bucket:      aws.String(o.objectStoreBucket),
		Key:         aws.String("object/" + id),
		ContentType: aws.String(contentType),
		Body:        aws.ReadSeekCloser(reader),
	}

	_, err := o.s3Client.PutObjectWithContext(ctx, input)

	if err != nil {
		panic(errors.Wrapf(err, "cannot put object(%s) to S3", id))
	}
}
