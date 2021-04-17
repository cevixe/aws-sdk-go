package sns

import (
	"github.com/aws/aws-sdk-go/service/sns"
	"github.com/aws/aws-sdk-go/service/sns/snsiface"
	"github.com/cevixe/aws-sdk-go/aws/integration/session"
)

const AwsSnsFactory string = "AwsSnsFactory"

type ClientFactory interface {
	GetClient(region string) snsiface.SNSAPI
}

type clientFactoryImpl struct {
	ClientFactory
	clients        map[string]snsiface.SNSAPI
	sessionFactory session.Factory
}

func (f clientFactoryImpl) GetClient(region string) snsiface.SNSAPI {
	if f.clients[region] != nil {
		return f.clients[region]
	}

	sess := f.sessionFactory.GetSession(region)
	newClient := sns.New(sess)
	f.clients[region] = newClient

	return newClient
}

func NewClientFactory(sessionFactory session.Factory) ClientFactory {
	return &clientFactoryImpl{
		clients:        make(map[string]snsiface.SNSAPI),
		sessionFactory: sessionFactory,
	}
}
