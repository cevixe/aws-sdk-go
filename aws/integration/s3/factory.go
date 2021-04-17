package s3

import (
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/aws/aws-sdk-go/service/s3/s3iface"
	"github.com/cevixe/aws-sdk-go/aws/integration/session"
)

const AwsS3Factory string = "AwsS3Factory"

type ClientFactory interface {
	GetClient(region string) s3iface.S3API
}

type clientFactoryImpl struct {
	ClientFactory
	clients        map[string]s3iface.S3API
	sessionFactory session.Factory
}

func (f clientFactoryImpl) GetClient(region string) s3iface.S3API {
	if f.clients[region] != nil {
		return f.clients[region]
	}

	sess := f.sessionFactory.GetSession(region)
	newClient := s3.New(sess)
	f.clients[region] = newClient

	return newClient
}

func NewClientFactory(sessionFactory session.Factory) ClientFactory {
	return &clientFactoryImpl{
		clients:        make(map[string]s3iface.S3API),
		sessionFactory: sessionFactory,
	}
}
