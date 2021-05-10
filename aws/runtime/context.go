package runtime

import (
	"context"
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
	awsObjectStore := s3.NewDefaultS3ObjectStore(sessionFactory)
	awsEventStore := dynamodb.NewDefaultDynamodbEventStore(sessionFactory)
	awsStateStore := dynamodb.NewDefaultDynamodbStateStore(sessionFactory)
	awsEventBus := sns.NewDefaultSnsEventBus(sessionFactory)
	awsGraphqlGateway := appsync.NewDefaultAppsyncGateway(sessionFactory)

	awsContext := &impl.Context{
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
	eventStore := impl.NewEventStore(awsContext.AwsEventStore, awsContext.AwsObjectStore)
	stateStore := impl.NewStateStore(eventStore)
	eventFactory := impl.NewEventFactory()

	ctx = context.WithValue(ctx, cevixe.CevixeEventFactory, eventFactory)
	ctx = context.WithValue(ctx, cevixe.CevixeEventStore, eventStore)
	ctx = context.WithValue(ctx, cevixe.CevixeStateStore, stateStore)

	return ctx
}
