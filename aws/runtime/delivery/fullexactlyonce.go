package delivery

import (
	"context"
	"fmt"
	"github.com/cevixe/aws-sdk-go/aws/impl"
	"github.com/cevixe/aws-sdk-go/aws/model"
	"github.com/cevixe/aws-sdk-go/aws/runtime/control"
	"github.com/cevixe/core-sdk-go/cevixe"
	"github.com/cevixe/core-sdk-go/core"
	"time"
)

func FullExactlyOnce(ctx context.Context, delegate core.EventHandler) core.EventHandler {

	awsContext := ctx.Value(impl.AwsContext).(*impl.Context)
	return func(ctx context.Context, event core.Event) core.Event {

		controlGroup := control.GetControlGroup(ctx, event)

		controlRecords := awsContext.AwsEventStore.GetControlRecords(ctx, controlGroup)
		lastControlRecord := controlRecords[0]

		if lastControlRecord.ControlType == model.ConfirmControl {
			return nil
		}

		currentTime := time.Now().UnixNano() / int64(time.Millisecond)
		if (currentTime - lastControlRecord.ControlTime) < int64(lastControlRecord.HandlerTimeout) {
			panic(fmt.Errorf("invalid event handling reintent"))
		}

		nextIntent := uint64(int64(lastControlRecord.ControlIntent) + 1)
		blockControlRecord := control.NewControlRecord(ctx, event, model.BlockControl, nextIntent)
		awsContext.AwsEventStore.CreateControlRecord(ctx, blockControlRecord)

		confirmControlRecord := control.NewControlRecord(ctx, event, model.ConfirmControl, nextIntent)
		defer func() {
			HandleRecovery(ctx, recover(), confirmControlRecord)
		}()

		newCevixeEvent := delegate(ctx, event)

		if newCevixeEvent == nil {
			factory := ctx.Value(cevixe.CevixeEventFactory).(core.EventFactory)
			newCevixeEvent = factory.NewSystemEvent(ctx, core.NoEventGenerated{Handler: awsContext.AwsHandlerID})
		}

		eventRecord := newCevixeEvent.(*impl.EventImpl).Record
		awsContext.AwsEventStore.CreateControlledEventRecord(ctx, eventRecord, confirmControlRecord)

		return newCevixeEvent
	}
}
