package impl

import (
	"context"
	"github.com/cevixe/aws-sdk-go/aws/model"
	"github.com/cevixe/aws-sdk-go/util"
	"github.com/cevixe/core-sdk-go/core"
	"time"
)

type eventImpl struct {
	core.Event
	ctx   context.Context
	value *model.EventObject
}

func NewEvent(
	ctx context.Context,
	eventObjectValue *model.EventObject) core.Event {
	return &eventImpl{
		ctx:   ctx,
		value: eventObjectValue,
	}
}

func (e *eventImpl) ID() uint64 {
	return e.value.EventID
}

func (e *eventImpl) Type() string {
	return e.value.EventType
}

func (e *eventImpl) Time() time.Time {
	nanoseconds := e.value.EventTime * int64(time.Millisecond)
	return time.Unix(0, nanoseconds)
}

func (e *eventImpl) Author() string {
	return e.value.EventAuthor
}

func (e *eventImpl) Payload(v interface{}) {
	if e.value.EventPayload == nil && e.value.Reference != nil {
		awsContext := e.ctx.Value(model.AwsContext).(*model.Context)
		newValue := &model.EventObject{}
		awsContext.AwsObjectStore.GetObject(e.ctx, e.value.Reference, newValue)
		e.value = newValue
	}
	json := util.MarshalJsonString(e.value.EventPayload)
	util.UnmarshalJsonString(json, v)
}

func (e *eventImpl) Source() core.Entity {
	return &entityImpl{
		ctx:       e.ctx,
		lastEvent: e.value}
}

func (e *eventImpl) Transaction() string {
	return e.value.Transaction
}

func (e *eventImpl) Trigger() core.Event {
	awsContext := e.ctx.Value(model.AwsContext).(model.Context)
	triggerValue := awsContext.AwsEventStore.GetEventById(e.ctx, e.value.TriggerSource, e.value.TriggerID)
	return NewEvent(e.ctx, triggerValue)
}
