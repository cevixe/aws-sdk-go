package delivery

import (
	"context"
	"github.com/cevixe/aws-sdk-go/aws/impl"
	"github.com/cevixe/core-sdk-go/core"
)

func AtLeastOnce(ctx context.Context, delegate core.EventHandler) core.EventHandler {

	awsContext := ctx.Value(impl.AwsContext).(*impl.Context)
	return func(ctx context.Context, event core.Event) core.Event {
		newCevixeEvent := delegate(ctx, event)

		if newCevixeEvent != nil {
			record := newCevixeEvent.(*impl.EventImpl).Record
			awsContext.AwsEventStore.CreateUncontrolledEventRecord(ctx, record)
		}

		return newCevixeEvent
	}
}
