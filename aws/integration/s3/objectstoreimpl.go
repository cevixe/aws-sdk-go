package s3

import (
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
	"strings"
)

type objectStoreImpl struct {
	objectStoreBucket string
	objectStoreRegion string
	s3ClientsCache    map[string]s3iface.S3API
	sessionFactory    session.Factory
}

func NewS3ObjectStore(
	objectStoreBucket string,
	objectStoreRegion string,
	sessionFactory session.Factory) model.AwsObjectStore {

	return &objectStoreImpl{
		objectStoreBucket: objectStoreBucket,
		objectStoreRegion: objectStoreRegion,
		s3ClientsCache:    map[string]s3iface.S3API{},
		sessionFactory:    sessionFactory,
	}
}

func NewDefaultS3ObjectStore(sessionFactory session.Factory) model.AwsObjectStore {

	region := os.Getenv(env.AwsRegion)
	objectStoreBucketName := os.Getenv(env.CevixeObjectStoreBucketName)

	return NewS3ObjectStore(objectStoreBucketName, region, sessionFactory)
}

func (o *objectStoreImpl) getClient(region string) s3iface.S3API {
	if o.s3ClientsCache[region] != nil {
		return o.s3ClientsCache[region]
	}

	sess := o.sessionFactory.GetSession(region)
	newClient := s3.New(sess)
	o.s3ClientsCache[region] = newClient

	return newClient
}

func (o objectStoreImpl) GetObject(ctx context.Context, reference *model.AwsObjectStoreReference, v interface{}) {

	input := &s3.GetObjectInput{
		Bucket:    aws.String(reference.Bucket),
		Key:       aws.String(reference.Key),
		VersionId: reference.Version,
	}

	client := o.getClient(reference.Region)
	output, err := client.GetObjectWithContext(ctx, input)

	if err != nil {
		panic(fmt.Errorf("cannot get reference(%v) from S3\n%v", reference, err))
	}

	buffer, err := ioutil.ReadAll(output.Body)
	if err != nil {
		panic(fmt.Errorf("cannot unmarshal object(%v) from S3\n%v", reference, err))
	}

	util2.UnmarshalJson(buffer, v)
}

func (o objectStoreImpl) SaveObject(ctx context.Context, key string, v interface{}) *model.AwsObjectStoreReference {

	json := util2.MarshalJsonString(v)
	reader := strings.NewReader(json)

	input := &s3.PutObjectInput{
		Bucket:      aws.String(o.objectStoreBucket),
		Key:         aws.String(key),
		ContentType: aws.String("application/json"),
		Body:        aws.ReadSeekCloser(reader),
	}

	client := o.getClient(o.objectStoreRegion)
	output, err := client.PutObjectWithContext(ctx, input)

	if err != nil {
		panic(fmt.Errorf("cannot put object(%s) to S3\n%v", key, err))
	}

	return &model.AwsObjectStoreReference{
		Region:  o.objectStoreRegion,
		Bucket:  o.objectStoreBucket,
		Key:     key,
		Version: output.VersionId,
	}
}
