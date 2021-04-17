package impl

import (
	"context"
	"github.com/cevixe/aws-sdk-go/aws/model"
	"github.com/cevixe/core-sdk-go/core"
)

type eventBusImpl struct {
	core.EventBus
	eventBus model.EventBus
}

func NewEventBus(eventBus model.EventBus) core.EventBus {
	return &eventBusImpl{
		eventBus: eventBus,
	}
}

func (e eventBusImpl) PublishEvent(ctx context.Context, event core.Event) {
	eventObject := event.(*eventImpl).value
	e.eventBus.PublishEvent(ctx, eventObject)
}
