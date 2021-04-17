package impl

import (
	"context"
	"github.com/cevixe/aws-sdk-go/aws/model"
	"github.com/cevixe/core-sdk-go/core"
	"strconv"
)

type eventStoreImpl struct {
	core.EventStore
	enableObjectStore bool
	eventStore        model.EventStore
	objectStore       model.ObjectStore
}

func NewEventStore(eventStore model.EventStore, objectStore model.ObjectStore, enableObjectStore bool) core.EventStore {
	return &eventStoreImpl{
		eventStore:        eventStore,
		objectStore:       objectStore,
		enableObjectStore: enableObjectStore,
	}
}

func (e eventStoreImpl) GetEvent(ctx context.Context, source string, id *uint64) core.Event {
	if id == nil {
		eventValue := e.eventStore.GetLatestEvent(ctx, source)
		return NewEvent(ctx, eventValue)
	} else {
		eventValue := e.eventStore.GetEventById(ctx, source, *id)
		return NewEvent(ctx, eventValue)
	}
}

func (e eventStoreImpl) SaveEvent(ctx context.Context, event core.Event) {
	eventObject := event.(*eventImpl).value

	if e.enableObjectStore {
		eventObject.SourceState = nil
		eventObject.EventPayload = nil
		key := generateEventKey(event.Source().Type(), event.Source().ID(), event.ID())
		eventObject.Reference = e.objectStore.SaveObject(ctx, key, eventObject)
	}
	e.eventStore.SaveEvent(ctx, eventObject)
}

func generateEventKey(sourceType string, sourceId string, id uint64) string {
	return sourceType + "/" + sourceId + "/" + strconv.FormatUint(id, 10)
}
