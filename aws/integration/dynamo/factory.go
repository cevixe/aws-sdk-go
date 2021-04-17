package dynamo

import (
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbiface"
	"github.com/cevixe/aws-sdk-go/aws/integration/session"
)
const AwsDynamoFactory string = "AwsDynamoFactory"

type ClientFactory interface {
	GetClient(region string) dynamodbiface.DynamoDBAPI
}

type clientFactoryImpl struct {
	ClientFactory
	clients        map[string]dynamodbiface.DynamoDBAPI
	sessionFactory session.Factory
}

func (f clientFactoryImpl) GetClient(region string) dynamodbiface.DynamoDBAPI {
	if f.clients[region] != nil {
		return f.clients[region]
	}

	sess := f.sessionFactory.GetSession(region)
	newClient := dynamodb.New(sess)
	f.clients[region] = newClient

	return newClient
}

func NewClientFactory(sessionFactory session.Factory) ClientFactory {
	return &clientFactoryImpl{
		clients:        make(map[string]dynamodbiface.DynamoDBAPI),
		sessionFactory: sessionFactory,
	}
}
