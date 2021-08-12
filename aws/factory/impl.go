package factory

import (
	"fmt"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbiface"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/aws/aws-sdk-go/service/s3/s3iface"
	"github.com/aws/aws-sdk-go/service/sns"
	"github.com/aws/aws-sdk-go/service/sns/snsiface"
	"github.com/cevixe/aws-sdk-go/aws/env"
	"github.com/cevixe/aws-sdk-go/aws/integration/session"
	"os"
)

type factoryImpl struct {
	region string
	sess   session.Factory
	cache  map[string]interface{}
}

func (f *factoryImpl) cacheKey(service string, region string) string {
	return fmt.Sprintf("%s#%s", service, region)
}

func (f *factoryImpl) toCache(service string, client interface{}, region string) {
	f.cache[f.cacheKey(service, region)] = client
}

func (f *factoryImpl) fromCache(service string, region string) interface{} {
	return f.cache[f.cacheKey(service, region)]
}

func (f factoryImpl) DynamodbClient(region ...string) dynamodbiface.DynamoDBAPI {
	reg := f.readRegion(region)
	client := f.fromCache(DynamoDB, reg).(dynamodbiface.DynamoDBAPI)

	if client != nil {
		return client
	} else {
		client = dynamodb.New(f.sess.GetSession(reg))
		f.toCache(DynamoDB, client, reg)
		return client
	}
}

func (f factoryImpl) SnsClient(region ...string) snsiface.SNSAPI {
	reg := f.readRegion(region)
	client := f.fromCache(SNS, reg).(snsiface.SNSAPI)

	if client != nil {
		return client
	} else {
		client = sns.New(f.sess.GetSession(reg))
		f.toCache(SNS, client, reg)
		return client
	}
}

func (f factoryImpl) S3Client(region ...string) s3iface.S3API {
	reg := f.readRegion(region)
	client := f.fromCache(S3, reg).(s3iface.S3API)

	if client != nil {
		return client
	} else {
		client = s3.New(f.sess.GetSession(reg))
		f.toCache(S3, client, reg)
		return client
	}
}

func (f factoryImpl) readRegion(region []string) string {
	if len(region) == 0 {
		return f.region
	} else {
		return region[0]
	}
}

func New(sess session.Factory) AwsFactory {
	region := os.Getenv(env.AwsRegion)
	return &factoryImpl{region: region, sess: sess}
}
