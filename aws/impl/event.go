package impl

import (
	"context"
	"github.com/cevixe/aws-sdk-go/aws/model"
	"github.com/cevixe/aws-sdk-go/aws/util"
	"github.com/cevixe/core-sdk-go/core"
	"time"
)

type EventImpl struct {
	Ctx    context.Context
	Record *model.AwsEventRecord
}

func NewEvent(ctx context.Context, record *model.AwsEventRecord) core.Event {
	return &EventImpl{
		Ctx:    ctx,
		Record: record,
	}
}

func (e EventImpl) ID() string {
	return *e.Record.EventID
}

func (e EventImpl) Source() string {
	return *e.Record.EventSource
}

func (e EventImpl) Class() core.EventClass {
	return core.EventClass(*e.Record.EventClass)
}

func (e EventImpl) Type() string {
	return *e.Record.EventType
}

func (e EventImpl) Time() time.Time {
	nanoseconds := *e.Record.EventTime * int64(time.Millisecond)
	return time.Unix(0, nanoseconds)
}

func (e EventImpl) Author() string {
	return *e.Record.EventAuthor
}

func (e EventImpl) Data(v interface{}) {
	json := util.MarshalJsonString(e.Record.EventData)
	util.UnmarshalJsonString(json, v)
}

func (e EventImpl) Entity() core.Entity {
	if core.EventClass(*e.Record.EventClass) == core.CommandEvent ||
		core.EventClass(*e.Record.EventClass) == core.BusinessEvent {
		return nil
	}
	return &EntityImpl{LastEvent: e.Record}
}

func (e EventImpl) Transaction() string {
	return *e.Record.Transaction
}

func (e EventImpl) Trigger() core.Event {
	awsContext := e.Ctx.Value(AwsContext).(*Context)
	triggerValue := awsContext.AwsEventStore.GetEventRecordByID(e.Ctx, *e.Record.TriggerSource, *e.Record.TriggerID)
	return NewEvent(e.Ctx, triggerValue)
}
