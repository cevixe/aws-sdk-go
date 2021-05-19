package s3

import (
	"bytes"
	"context"
	"fmt"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/aws/aws-sdk-go/service/s3/s3iface"
	"github.com/cevixe/aws-sdk-go/aws/env"
	"github.com/cevixe/aws-sdk-go/aws/integration/session"
	"github.com/cevixe/aws-sdk-go/aws/model"
	util2 "github.com/cevixe/aws-sdk-go/aws/util"
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

func NewDefaultS3ObjectStore(sessionFactory session.Factory) model.AwsObjectStore {

	region := os.Getenv(env.AwsRegion)
	objectStoreBucketName := os.Getenv(env.CevixeObjectStoreBucketName)
	s3Client := s3.New(sessionFactory.GetSession(region))

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
		Key:    aws.String(id),
	}

	output, err := o.s3Client.GetObjectWithContext(ctx, input)

	if err != nil {
		panic(fmt.Errorf("cannot get reference(%v) from S3\n%v", id, err))
	}

	buffer, err := ioutil.ReadAll(output.Body)
	if err != nil {
		panic(fmt.Errorf("cannot unmarshal object(%v) from S3\n%v", id, err))
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
		Key:         aws.String(id),
		ContentType: aws.String(contentType),
		Body:        aws.ReadSeekCloser(reader),
	}

	_, err := o.s3Client.PutObjectWithContext(ctx, input)

	if err != nil {
		panic(fmt.Errorf("cannot put object(%s) to S3\n%v", id, err))
	}
}