package runtime

import (
	"context"
	"github.com/aws/aws-lambda-go/events"
	"github.com/cevixe/aws-sdk-go/aws/env"
	"github.com/cevixe/aws-sdk-go/aws/impl"
	"github.com/cevixe/aws-sdk-go/aws/integration/sqs"
	"github.com/cevixe/aws-sdk-go/aws/model"
	"github.com/cevixe/core-sdk-go/cevixe"
	"github.com/cevixe/core-sdk-go/core"
	"github.com/pkg/errors"
	"os"
	"strconv"
)

func readCevixeEvent(ctx context.Context, event events.SQSEvent) core.Event {

	eventRecord := &model.AwsEventRecord{}
	sqs.UnmarshallSQSEvent(event, eventRecord)
	return impl.NewEvent(ctx, eventRecord)
}

func NewHandler(delegate core.EventHandler) func(ctx context.Context, event events.SQSEvent) error {

	return func(ctx context.Context, event events.SQSEvent) (err error) {

		awsHandlerTimeoutString := os.Getenv(env.AwsLambdaFunctionTimeout)
		awsHandlerTimeout, err := strconv.ParseUint(awsHandlerTimeoutString, 10, 64)
		if err != nil {
			panic(errors.Wrap(err, "cannot get handler timeout"))
		}

		awsContext := ctx.Value(impl.AwsContext).(*impl.Context)
		awsContext.AwsHandlerID = os.Getenv(env.AwsLambdaFunctionName)
		awsContext.AwsHandlerVersion = os.Getenv(env.AwsLambdaFunctionVersion)
		awsContext.AwsHandlerTimeout = awsHandlerTimeout
			
		cevixeEvent := readCevixeEvent(ctx, event)

		ctx = context.WithValue(ctx, cevixe.CevixeUserID, cevixeEvent.Author())
		ctx = context.WithValue(ctx, cevixe.CevixeTransaction, cevixeEvent.Transaction())
		ctx = context.WithValue(ctx, cevixe.CevixeEventTrigger, cevixeEvent)
		delegate(ctx, cevixeEvent)

		return nil
	}
}
