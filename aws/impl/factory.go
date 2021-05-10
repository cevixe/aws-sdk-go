package impl

import (
	"context"
	"fmt"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/cevixe/aws-sdk-go/aws/model"
	"github.com/cevixe/aws-sdk-go/aws/util"
	"github.com/cevixe/core-sdk-go/cevixe"
	"github.com/cevixe/core-sdk-go/core"
	"github.com/google/uuid"
	"github.com/stoewer/go-strcase"
	"strconv"
	"time"
)

type eventFactoryImpl struct {
}

func NewEventFactory() core.EventFactory {
	return &eventFactoryImpl{}
}

func (f eventFactoryImpl) NewCommandEvent(ctx context.Context, data interface{}) core.Event {
	return newDefaultEvent(ctx, core.CommandEvent, data, nil, nil)
}

func (f eventFactoryImpl) NewBusinessEvent(ctx context.Context, data interface{}) core.Event {
	return newDefaultEvent(ctx, core.BusinessEvent, data, nil, nil)
}

func (f eventFactoryImpl) NewDomainEvent(ctx context.Context, data interface{}, entity core.Entity, state interface{}) core.Event {
	return newDefaultEvent(ctx, core.DomainEvent, data, entity, state)
}

func newDefaultEvent(ctx context.Context, class core.EventClass, data interface{}, entity core.Entity, state interface{}) core.Event {

	trigger := ctx.Value(cevixe.CevixeEventTrigger).(core.Event)
	eventType := util.GetTypeName(data)
	eventTime := time.Now().UnixNano() / int64(time.Millisecond)

	eventRecord := &model.AwsEventRecord{
		EventClass:    aws.String(string(class)),
		EventType:     aws.String(eventType),
		EventTime:     aws.Int64(eventTime),
		EventAuthor:   aws.String(trigger.Author()),
		EventData:     toGenericData(data),
		TriggerSource: aws.String(trigger.Source()),
		TriggerID:     aws.String(trigger.ID()),
		Transaction:   aws.String(trigger.Transaction()),
	}
	addEntityMetadata(class, entity, state, eventRecord)
	addEventIdentity(class, data, entity, state, eventRecord)

	return NewEvent(ctx, eventRecord)
}

func addEntityMetadata(class core.EventClass, entity core.Entity, state interface{}, record *model.AwsEventRecord) {

	if class != core.DomainEvent {
		return
	}
	if entity == nil {
		record.EntityID = aws.String(uuid.NewString())
		record.EntityType = aws.String(util.GetTypeName(state))
		record.EntityCreatedAt = record.EventTime
		record.EntityCreatedBy = record.EventAuthor
	} else {
		record.EntityID = aws.String(entity.ID())
		record.EntityType = aws.String(entity.Type())
		record.EntityCreatedAt = aws.Int64(entity.CreatedAt().UnixNano() / int64(time.Millisecond))
		record.EntityCreatedBy = aws.String(entity.CreatedBy())
	}
	record.EntityState = toGenericData(state)
}

func addEventIdentity(class core.EventClass, data interface{}, entity core.Entity, state interface{}, record *model.AwsEventRecord) {

	switch class {
	case core.CommandEvent:
		dataTypeName := strcase.KebabCase(util.GetTypeName(data))
		record.EventSource = aws.String("/command/" + dataTypeName)
		record.EventID = aws.String(uuid.NewString())
		break
	case core.DomainEvent:
		if entity == nil {
			entityTypeName := strcase.KebabCase(util.GetTypeName(state))
			record.EventSource = aws.String("/domain/" + entityTypeName + "/" + *record.EntityID)
			record.EventID = aws.String("1")
		} else {
			entityTypeName := strcase.KebabCase(entity.Type())
			record.EventSource = aws.String("/domain/" + entityTypeName + "/" + *record.EntityID)
			record.EventID = aws.String(strconv.FormatUint(entity.Version()+1, 10))
		}
		break
	case core.BusinessEvent:
		dataTypeName := strcase.KebabCase(util.GetTypeName(data))
		record.EventSource = aws.String("/business/" + dataTypeName)
		record.EventID = aws.String(uuid.NewString())
		break
	default:
		panic(fmt.Errorf("event class not supported"))
	}
}

func toGenericData(object interface{}) *map[string]interface{} {
	payload := &map[string]interface{}{}
	objectJson := util.MarshalJson(object)
	util.UnmarshalJson(objectJson, payload)
	return payload
}
