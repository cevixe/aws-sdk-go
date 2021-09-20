package runtime

import (
	"context"
	"github.com/cevixe/aws-sdk-go/aws/factory"
	"github.com/cevixe/aws-sdk-go/aws/http"
	"github.com/cevixe/aws-sdk-go/aws/impl"
	"github.com/cevixe/aws-sdk-go/aws/integration/appsync"
	"github.com/cevixe/aws-sdk-go/aws/integration/dynamodb"
	"github.com/cevixe/aws-sdk-go/aws/integration/s3"
	"github.com/cevixe/aws-sdk-go/aws/integration/session"
	"github.com/cevixe/aws-sdk-go/aws/integration/sns"
	"github.com/cevixe/core-sdk-go/cevixe"
)

func NewContext() context.Context {
	ctx := context.Background()
	ctx = configureAwsContext(ctx)
	ctx = configureCevixeContext(ctx)
	return ctx
}

func configureAwsContext(ctx context.Context) context.Context {

	client := http.NewDefaultHttpClient()
	sessionFactory := session.NewSessionFactory(client)
	awsFactory := factory.New(sessionFactory)
	awsObjectStore := s3.NewDefaultS3ObjectStore(awsFactory)
	awsEventStore := dynamodb.NewDefaultDynamodbEventStore(awsFactory)
	awsStateStore := dynamodb.NewDefaultDynamodbStateStore(awsFactory)
	awsEventBus := sns.NewDefaultSnsEventBus(awsFactory)
	awsGraphqlGateway := appsync.NewDefaultAppsyncGateway(sessionFactory)

	awsContext := &impl.Context{
		AwsFactory:        awsFactory,
		AwsObjectStore:    awsObjectStore,
		AwsEventStore:     awsEventStore,
		AwsEventBus:       awsEventBus,
		AwsStateStore:     awsStateStore,
		AwsGraphqlGateway: awsGraphqlGateway,
	}

	return context.WithValue(ctx, impl.AwsContext, awsContext)
}

func configureCevixeContext(ctx context.Context) context.Context {

	awsContext := ctx.Value(impl.AwsContext).(*impl.Context)
	eventStore := impl.NewEventStore(awsContext.AwsEventStore)
	stateStore := impl.NewStateStore(awsContext.AwsEventStore, awsContext.AwsStateStore)
	eventFactory := impl.NewEventFactory()

	ctx = context.WithValue(ctx, cevixe.CevixeEventFactory, eventFactory)
	ctx = context.WithValue(ctx, cevixe.CevixeEventStore, eventStore)
	ctx = context.WithValue(ctx, cevixe.CevixeStateStore, stateStore)

	return ctx
}
