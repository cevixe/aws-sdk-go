package runtime

import (
	"context"
	"fmt"
	"github.com/aws/aws-lambda-go/events"
	"github.com/cevixe/aws-sdk-go/aws/impl"
	"github.com/cevixe/aws-sdk-go/aws/integration/sqs"
	"github.com/cevixe/aws-sdk-go/aws/model"
	"github.com/cevixe/core-sdk-go/cevixe"
	"github.com/cevixe/core-sdk-go/core"
	"runtime/debug"
)

func readCevixeEvent(ctx context.Context, event events.SQSEvent) core.Event {

	awsContext := ctx.Value(impl.AwsContext).(*impl.Context)
	eventRecord := &model.AwsEventRecord{}
	sqs.UnmarshallSQSEvent(event, eventRecord)
	if eventRecord.Reference != nil && eventRecord.EventData == nil {
		awsContext.AwsObjectStore.GetObject(ctx, eventRecord.Reference, eventRecord)
	}
	return impl.NewEvent(ctx, eventRecord)
}

func NewHandler(delegate core.EventHandler) func(ctx context.Context, event events.SQSEvent) error {

	return func(ctx context.Context, event events.SQSEvent) (err error) {

		defer func() {
			if r := recover(); r != nil {
				fmt.Println("An unexpected error has occurred: ", r)
				fmt.Printf("Stack trace: \n%s\n", string(debug.Stack()))
				err = fmt.Errorf("An unexpected error has occurred: \n%v\n", r)
			}
		}()

		cevixeEvent := readCevixeEvent(ctx, event)
		ctx = context.WithValue(ctx, cevixe.CevixeEventTrigger, cevixeEvent)
		delegate(ctx, cevixeEvent)
		return nil
	}
}
