package impl

import (
	"context"
	"github.com/cevixe/aws-sdk-go/aws/model"
	"github.com/cevixe/core-sdk-go/core"
)

type eventStoreImpl struct {
	eventStore  model.AwsEventStore
	objectStore model.AwsObjectStore
}

func NewEventStore(eventStore model.AwsEventStore, objectStore model.AwsObjectStore) core.EventStore {
	return &eventStoreImpl{
		eventStore:  eventStore,
		objectStore: objectStore,
	}
}

func (e eventStoreImpl) GetLastEvent(ctx context.Context, source string) core.Event {
	eventValue := e.eventStore.GetLastEventRecord(ctx, source)
	if eventValue == nil {
		return nil
	}
	if eventValue.Reference != nil && eventValue.EventData == nil {
		e.objectStore.GetObject(ctx, eventValue.Reference, eventValue)
	}
	return NewEvent(ctx, eventValue)
}

func (e eventStoreImpl) GetEventByID(ctx context.Context, source string, id string) core.Event {
	eventValue := e.eventStore.GetEventRecordByID(ctx, source, id)
	if eventValue == nil {
		return nil
	}
	if eventValue.Reference != nil && eventValue.EventData == nil {
		e.objectStore.GetObject(ctx, eventValue.Reference, eventValue)
	}
	return NewEvent(ctx, eventValue)
}
