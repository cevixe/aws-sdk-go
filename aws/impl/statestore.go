package impl

import (
	"context"
	"github.com/cevixe/core-sdk-go/core"
	"github.com/stoewer/go-strcase"
	"strconv"
)

type stateStoreImpl struct {
	eventStore core.EventStore
}

func (s stateStoreImpl) GetLastVersion(ctx context.Context, typ string, id string) core.Entity {
	source := generateEventSource(typ, id)
	lastEvent := s.eventStore.GetLastEvent(ctx, source)
	if lastEvent == nil {
		return nil
	}
	return lastEvent.Entity()
}

func (s stateStoreImpl) GetByVersion(ctx context.Context, typ string, id string, version uint64) core.Entity {
	source := generateEventSource(typ, id)
	event := s.eventStore.GetEventByID(ctx, source, strconv.FormatUint(version, 10))
	if event == nil {
		return nil
	}
	return event.Entity()
}

func generateEventSource(typ string, id string) string {
	entityTypeName := strcase.KebabCase(typ)
	return "/domain/" + entityTypeName + "/" + id
}

func NewStateStore(eventStore core.EventStore) core.StateStore {
	return &stateStoreImpl{
		eventStore: eventStore,
	}
}
