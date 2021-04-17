package runtime

import (
	"context"
	runtime "github.com/aws/aws-lambda-go/lambda"
	"github.com/cevixe/aws-sdk-go/aws/impl"
	"github.com/cevixe/aws-sdk-go/aws/integration/appsync"
	"github.com/cevixe/aws-sdk-go/aws/integration/dynamo"
	"github.com/cevixe/aws-sdk-go/aws/integration/lambda"
	"github.com/cevixe/aws-sdk-go/aws/integration/s3"
	"github.com/cevixe/aws-sdk-go/aws/integration/session"
	"github.com/cevixe/aws-sdk-go/aws/integration/sns"
	"github.com/cevixe/aws-sdk-go/aws/model"
	"github.com/cevixe/core-sdk-go/cevixe"
	"github.com/cevixe/core-sdk-go/core"
	"os"
)

func Start(fn core.EventHandler) {
	ctx := NewContext()
	decorator := lambda.NewDecorator(ctx, fn)
	runtime.StartWithContext(ctx, decorator.Handle)
}

func NewContext() context.Context {

	ctx := context.Background()
	region := os.Getenv("AWS_REGION")
	sessionFactory := session.NewSessionFactory()
	sessionFactory.GetSession(region)

	dynamoClientFactory := dynamo.NewClientFactory(sessionFactory)
	s3ClientFactory := s3.NewClientFactory(sessionFactory)
	snsClientFactory := sns.NewClientFactory(sessionFactory)

	ctx = context.WithValue(ctx, session.AwsSessionFactory, sessionFactory)
	ctx = context.WithValue(ctx, dynamo.AwsDynamoFactory, dynamoClientFactory)
	ctx = context.WithValue(ctx, s3.AwsS3Factory, s3ClientFactory)
	ctx = context.WithValue(ctx, sns.AwsSnsFactory, snsClientFactory)

	enableObjectStore := os.Getenv("CVX_ENABLE_OBJECT_STORE") == "1"
	eventStoreTable := os.Getenv("CVX_EVENT_STORE")
	stateStoreTable := os.Getenv("CVX_STATE_STORE")
	eventBusTopic := os.Getenv("CVX_EVENT_BUS")
	objectStoreBucket := os.Getenv("CVX_OBJECT_STORE")
	graphqlGatewayUrl := os.Getenv("CVX_GRAPHQL_GATEWAY")

	awsObjectStore := s3.NewObjectStore(region, objectStoreBucket, s3ClientFactory)
	awsEventStore := dynamo.NewEventStore(region, eventStoreTable, dynamoClientFactory)
	awsEventBus := sns.NewEventBus(region, eventBusTopic, snsClientFactory)
	awsStateStore := dynamo.NewStateStore(region, stateStoreTable, dynamoClientFactory)
	awsGraphqlGateway := appsync.NewGraphqlGateway(region, graphqlGatewayUrl, sessionFactory)

	awsContext := &model.Context{
		AwsObjectStore:  awsObjectStore,
		AwsEventStore:   awsEventStore,
		AwsEventBus:     awsEventBus,
		AwsStateStore:   awsStateStore,
		AwsGraphGateway: awsGraphqlGateway,
	}

	eventBus := impl.NewEventBus(awsEventBus)
	stateStore := impl.NewStateStore(awsStateStore)
	eventStore := impl.NewEventStore(awsEventStore, awsObjectStore, enableObjectStore)
	eventFactory := impl.NewEventFactory()

	ctx = context.WithValue(ctx, cevixe.CevixeEventFactory, eventFactory)
	ctx = context.WithValue(ctx, cevixe.CevixeEventStore, eventStore)
	ctx = context.WithValue(ctx, cevixe.CevixeEventBus, eventBus)
	ctx = context.WithValue(ctx, cevixe.CevixeStateStore, stateStore)
	ctx = context.WithValue(ctx, model.AwsContext, awsContext)
	return ctx
}
