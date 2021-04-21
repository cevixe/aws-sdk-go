package lambda

import (
	"context"
	"fmt"
	"github.com/aws/aws-lambda-go/events"
	"github.com/cevixe/aws-sdk-go/aws/integration/sqs"
	"github.com/cevixe/core-sdk-go/cevixe"
	"github.com/cevixe/core-sdk-go/core"
	"runtime/debug"
)

type Decorator interface {
	Handle(ctx context.Context, event events.SQSEvent) error
}

type decoratorImpl struct {
	eventStore core.EventStore
	delegate   core.EventHandler
}

func (h decoratorImpl) Handle(ctx context.Context, event events.SQSEvent) (err error) {

	defer func() {
		if r := recover(); r != nil {
			fmt.Println("An unexpected error has occurred: ", r)
			fmt.Printf("Stack trace: \n%s\n", string(debug.Stack()))
			err = fmt.Errorf("An unexpected error has occurred: \n%v\n", r)
		}
	}()

	cevixeEvent := sqs.MapSQSEventToSingleCevixeEvent(ctx, event)
	requestCtx := context.WithValue(ctx, cevixe.CevixeEventTrigger, cevixeEvent)
	newCevixeEvent := h.delegate(requestCtx, cevixeEvent)

	if newCevixeEvent != nil {
		h.eventStore.SaveEvent(ctx, newCevixeEvent)
	}

	return nil
}

func NewDecorator(ctx context.Context, delegate core.EventHandler) Decorator {
	eventStore := ctx.Value(cevixe.CevixeEventStore).(core.EventStore)
	return &decoratorImpl{
		eventStore: eventStore,
		delegate:   delegate,
	}
}
