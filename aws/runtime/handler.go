package runtime

import (
	"context"
	"github.com/aws/aws-lambda-go/events"
	"github.com/cevixe/aws-sdk-go/aws/impl"
	"github.com/cevixe/aws-sdk-go/aws/integration/sqs"
	"github.com/cevixe/aws-sdk-go/aws/model"
	"github.com/cevixe/core-sdk-go/cevixe"
	"github.com/cevixe/core-sdk-go/core"
)

func readCevixeEvent(ctx context.Context, event events.SQSEvent) core.Event {

	eventRecord := &model.AwsEventRecord{}
	sqs.UnmarshallSQSEvent(event, eventRecord)
	return impl.NewEvent(ctx, eventRecord)
}

func NewHandler(delegate core.EventHandler) func(ctx context.Context, event events.SQSEvent) error {

	return func(ctx context.Context, event events.SQSEvent) (err error) {

		cevixeEvent := readCevixeEvent(ctx, event)
		ctx = context.WithValue(ctx, cevixe.CevixeEventTrigger, cevixeEvent)
		delegate(ctx, cevixeEvent)

		return nil
	}
}
