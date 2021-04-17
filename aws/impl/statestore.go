package impl

import (
	"context"
	"github.com/cevixe/aws-sdk-go/aws/model"
	"github.com/cevixe/core-sdk-go/core"
)

type stateStoreImpl struct {
	core.StateStore
	stateStore model.StateStore
}

func NewStateStore(stateStore model.StateStore) core.StateStore {
	return &stateStoreImpl{
		stateStore: stateStore,
	}
}

func (s stateStoreImpl) UpdateState(ctx context.Context, entities []core.Entity) {

	events := make([]*model.EventObject, len(entities))
	for idx, entity := range entities {
		payload := &map[string]interface{}{}
		entity.State(payload)
		events[idx] = entity.(*entityImpl).lastEvent
	}
	s.stateStore.UpdateState(ctx, events)
}
