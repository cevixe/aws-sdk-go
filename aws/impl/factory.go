package impl

import (
	"context"
	"github.com/cevixe/aws-sdk-go/aws/model"
	"github.com/cevixe/aws-sdk-go/util"
	"github.com/cevixe/core-sdk-go/cevixe"
	"github.com/cevixe/core-sdk-go/core"
	"github.com/google/uuid"
	"reflect"
	"time"
)

type eventFactoryImpl struct {
	core.EventFactory
}

func NewEventFactory() core.EventFactory {
	return &eventFactoryImpl{}
}

func (f eventFactoryImpl) NewEvent(
	ctx context.Context,
	entity core.Entity,
	payload interface{},
	state interface{}) core.Event {

	if entity == nil {
		return f.newFirstEvent(ctx, payload, state)
	} else {
		return f.newEvent(ctx, entity, payload, state)
	}
}

func (f eventFactoryImpl) newFirstEvent(
	ctx context.Context,
	payload interface{},
	state interface{}) core.Event {

	trigger := ctx.Value(cevixe.CevixeEventTrigger).(core.Event)

	eventTime := time.Now().UnixNano() / int64(time.Millisecond)
	eventPayload := &map[string]interface{}{}
	eventPayloadJson := util.MarshalJsonString(payload)
	util.UnmarshalJsonString(eventPayloadJson, eventPayload)
	eventType := reflect.TypeOf(payload).Name()
	if eventType == "" {
		eventType = reflect.ValueOf(payload).Type().Name()
	}

	entityID := uuid.New().String()
	entityState := &map[string]interface{}{}
	entityStateJson := util.MarshalJsonString(state)
	util.UnmarshalJsonString(entityStateJson, entityState)
	entityType := util.GetTypeName(state)

	eventObject := &model.EventObject{
		SourceID:      entityID,
		SourceType:    entityType,
		SourceTime:    eventTime,
		SourceOwner:   trigger.Author(),
		SourceState:   entityState,
		EventID:       1,
		EventType:     eventType,
		EventTime:     eventTime,
		EventAuthor:   trigger.Author(),
		EventPayload:  eventPayload,
		Transaction:   trigger.Transaction(),
		TriggerSource: "/" + trigger.Source().Type() + "/" + trigger.Source().ID(),
		TriggerID:     trigger.ID(),
	}

	return NewEvent(ctx, eventObject)
}

func (f eventFactoryImpl) newEvent(
	ctx context.Context,
	entity core.Entity,
	payload interface{},
	state interface{}) core.Event {

	trigger := ctx.Value(cevixe.CevixeEventTrigger).(core.Event)

	eventTime := time.Now().UnixNano() / int64(time.Millisecond)
	eventPayload := &map[string]interface{}{}
	eventPayloadJson := util.MarshalJsonString(payload)
	util.UnmarshalJsonString(eventPayloadJson, eventPayload)
	eventType := util.GetTypeName(payload)

	entityTime := entity.Time().UnixNano() / int64(time.Millisecond)
	entityState := &map[string]interface{}{}
	entityStateJson := util.MarshalJsonString(state)
	util.UnmarshalJsonString(entityStateJson, entityState)

	eventObject := &model.EventObject{
		SourceID:      entity.ID(),
		SourceType:    entity.Type(),
		SourceTime:    entityTime,
		SourceOwner:   entity.Owner(),
		SourceState:   entityState,
		EventID:       entity.Version() + 1,
		EventType:     eventType,
		EventTime:     eventTime,
		EventAuthor:   trigger.Author(),
		EventPayload:  eventPayload,
		Transaction:   trigger.Transaction(),
		TriggerSource: "/" + trigger.Source().Type() + "/" + trigger.Source().ID(),
		TriggerID:     trigger.ID(),
	}

	return NewEvent(ctx, eventObject)
}
