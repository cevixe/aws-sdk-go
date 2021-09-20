package delivery

import (
	"context"
	"github.com/cevixe/aws-sdk-go/aws/impl"
	"github.com/cevixe/aws-sdk-go/aws/model"
	"github.com/cevixe/core-sdk-go/cevixe"
	"github.com/cevixe/core-sdk-go/core"
	"log"
	"runtime"
)

func HandleRecovery(ctx context.Context, err interface{}, controlRecord *model.AwsControlRecord) {
	if err != nil {
		factory := ctx.Value(cevixe.CevixeEventFactory).(core.EventFactory)
		awsContext := ctx.Value(impl.AwsContext).(*impl.Context)

		buf := make([]byte, 1<<16)
		stackSize := runtime.Stack(buf, true)
		errorStack := string(buf[0:stackSize])

		log.Printf("Internal error: %v", err)
		log.Printf("%s\n", string(buf[0:stackSize]))

		newEvent := factory.NewSystemEvent(ctx, core.EventHandlingFailed{
			Handler:    awsContext.AwsHandlerID,
			Error:      err.(error).Error(),
			StackTrace: errorStack,
		})
		record := newEvent.(*impl.EventImpl).Record
		if controlRecord != nil {
			awsContext.AwsEventStore.CreateControlledEventRecord(ctx, record, controlRecord)
		} else {
			awsContext.AwsEventStore.CreateUncontrolledEventRecord(ctx, record)
		}
	}
}
