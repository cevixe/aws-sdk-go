package runtime

import (
	"context"
	"github.com/cevixe/aws-sdk-go/aws/env"
	"github.com/cevixe/aws-sdk-go/aws/factory"
	"github.com/cevixe/aws-sdk-go/aws/http"
	"github.com/cevixe/aws-sdk-go/aws/impl"
	"github.com/cevixe/aws-sdk-go/aws/integration/appsync"
	"github.com/cevixe/aws-sdk-go/aws/integration/dynamodb"
	"github.com/cevixe/aws-sdk-go/aws/integration/s3"
	"github.com/cevixe/aws-sdk-go/aws/integration/session"
	"github.com/cevixe/aws-sdk-go/aws/integration/sns"
	"github.com/cevixe/core-sdk-go/cevixe"
	"github.com/pkg/errors"
	"os"
	"strconv"
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
	awsCounterStore := dynamodb.NewDefaultDynamodbCounterStore(awsFactory)
	awsEventBus := sns.NewDefaultSnsEventBus(awsFactory)
	awsGraphqlGateway := appsync.NewDefaultAppsyncGateway(sessionFactory)

	awsHandlerTimeoutString := os.Getenv(env.AwsLambdaFunctionTimeout)
	awsHandlerTimeout, err := strconv.ParseUint(awsHandlerTimeoutString, 10, 64)
	if err != nil {
		panic(errors.Wrap(err, "cannot get handler timeout"))
	}
	awsContext := &impl.Context{
		AwsHandlerID:      os.Getenv(env.AwsLambdaFunctionName),
		AwsHandlerVersion: os.Getenv(env.AwsLambdaFunctionVersion),
		AwsHandlerTimeout: awsHandlerTimeout,
		AwsFactory:        awsFactory,
		AwsObjectStore:    awsObjectStore,
		AwsEventStore:     awsEventStore,
		AwsEventBus:       awsEventBus,
		AwsStateStore:     awsStateStore,
		AwsCounterStore:   awsCounterStore,
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
