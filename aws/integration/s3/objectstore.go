package s3

import (
	"context"
	"fmt"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/cevixe/aws-sdk-go/aws/model"
	"github.com/cevixe/aws-sdk-go/util"
	"io/ioutil"
	"strings"
)

type objectStoreImpl struct {
	objectStoreRegion string
	objectStoreBucket string
	clientFactory     ClientFactory
}

func (o objectStoreImpl) SaveObject(ctx context.Context, key string, v interface{}) *model.ObjectRef {

	json := util.MarshalJsonString(v)
	reader := strings.NewReader(json)

	input := &s3.PutObjectInput{
		Bucket:      aws.String(o.objectStoreBucket),
		Key:         aws.String(key),
		ContentType: aws.String("application/json"),
		Body:        aws.ReadSeekCloser(reader),
	}

	client := o.clientFactory.GetClient(o.objectStoreRegion)
	output, err := client.PutObjectWithContext(ctx, input)

	if err != nil {
		panic(fmt.Errorf("cannot put object(%s) to S3\n%v", key, err))
	}

	return &model.ObjectRef{
		Region:  o.objectStoreRegion,
		Bucket:  o.objectStoreBucket,
		Key:     key,
		Version: output.VersionId,
	}
}

func (o objectStoreImpl) GetObject(ctx context.Context, ref *model.ObjectRef, v interface{}) {

	input := &s3.GetObjectInput{
		Bucket:    aws.String(ref.Bucket),
		Key:       aws.String(ref.Key),
		VersionId: ref.Version,
	}

	client := o.clientFactory.GetClient(ref.Region)
	output, err := client.GetObjectWithContext(ctx, input)

	if err != nil {
		panic(fmt.Errorf("cannot get reference(%v) from S3\n%v", ref, err))
	}

	buffer, err := ioutil.ReadAll(output.Body)
	if err != nil {
		panic(fmt.Errorf("cannot unmarshal object(%v) from S3\n%v", ref, err))
	}

	util.UnmarshalJson(buffer, v)
}

func NewObjectStore(objectStoreRegion string, objectStoreBucket string, clientFactory ClientFactory) model.ObjectStore {
	return &objectStoreImpl{
		objectStoreRegion: objectStoreRegion,
		objectStoreBucket: objectStoreBucket,
		clientFactory:     clientFactory,
	}
}
