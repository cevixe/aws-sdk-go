package delivery

import (
	"context"
	"github.com/cevixe/aws-sdk-go/aws/impl"
	"github.com/cevixe/aws-sdk-go/aws/model"
	"github.com/cevixe/aws-sdk-go/aws/runtime/control"
	"github.com/cevixe/core-sdk-go/cevixe"
	"github.com/cevixe/core-sdk-go/core"
)

func AtMostOnce(ctx context.Context, delegate core.EventHandler) core.EventHandler {

	awsContext := ctx.Value(impl.AwsContext).(*impl.Context)
	return func(ctx context.Context, event core.Event) core.Event {

		controlRecord := control.NewControlRecord(ctx, event, model.BlockControl, 1)
		awsContext.AwsEventStore.CreateControlRecord(ctx, controlRecord)

		defer func() {
			HandleRecovery(ctx, recover(), nil)
		}()

		newCevixeEvent := delegate(ctx, event)

		if newCevixeEvent == nil {
			factory := ctx.Value(cevixe.CevixeEventFactory).(core.EventFactory)
			newCevixeEvent = factory.NewSystemEvent(ctx, core.NoEventGenerated{Handler: awsContext.AwsHandlerID})
		}

		eventRecord := newCevixeEvent.(*impl.EventImpl).Record
		awsContext.AwsEventStore.CreateUncontrolledEventRecord(ctx, eventRecord)

		return newCevixeEvent
	}
}
