package delivery

import (
	"context"
	"github.com/cevixe/aws-sdk-go/aws/impl"
	"github.com/cevixe/aws-sdk-go/aws/model"
	"github.com/cevixe/aws-sdk-go/aws/runtime/control"
	"github.com/cevixe/core-sdk-go/core"
)

func ExactlyOnce(ctx context.Context, delegate core.EventHandler) core.EventHandler {

	awsContext := ctx.Value(impl.AwsContext).(*impl.Context)
	return func(ctx context.Context, event core.Event) core.Event {

		newCevixeEvent := delegate(ctx, event)

		if newCevixeEvent != nil {
			eventRecord := newCevixeEvent.(*impl.EventImpl).Record
			controlRecord := control.NewControlRecord(ctx, event, model.ConfirmControl, 1)
			awsContext.AwsEventStore.CreateControlledEventRecord(ctx, eventRecord, controlRecord)
		}

		return newCevixeEvent
	}
}
