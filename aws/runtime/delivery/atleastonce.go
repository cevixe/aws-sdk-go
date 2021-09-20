package delivery

import (
	"context"
	"github.com/cevixe/aws-sdk-go/aws/impl"
	"github.com/cevixe/core-sdk-go/cevixe"
	"github.com/cevixe/core-sdk-go/core"
)

func AtLeastOnce(ctx context.Context, delegate core.EventHandler) core.EventHandler {

	awsContext := ctx.Value(impl.AwsContext).(*impl.Context)
	return func(ctx context.Context, event core.Event) core.Event {

		defer func() {
			HandleRecovery(ctx, recover(), nil)
		}()

		newCevixeEvent := delegate(ctx, event)

		if newCevixeEvent == nil {
			factory := ctx.Value(cevixe.CevixeEventFactory).(core.EventFactory)
			newCevixeEvent = factory.NewSystemEvent(ctx, core.NoEventGenerated{Handler: awsContext.AwsHandlerID})
		}

		record := newCevixeEvent.(*impl.EventImpl).Record
		awsContext.AwsEventStore.CreateUncontrolledEventRecord(ctx, record)
		
		return newCevixeEvent
	}
}
