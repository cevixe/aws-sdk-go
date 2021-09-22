package impl

import (
	"context"
	"fmt"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/cevixe/aws-sdk-go/aws/model"
	"github.com/cevixe/aws-sdk-go/aws/serdes/gzip"
	"github.com/cevixe/aws-sdk-go/aws/serdes/json"
	"github.com/cevixe/aws-sdk-go/aws/util"
	"github.com/cevixe/core-sdk-go/cevixe"
	"github.com/cevixe/core-sdk-go/core"
	"github.com/google/uuid"
	"github.com/pkg/errors"
	"github.com/stoewer/go-strcase"
	"time"
)

type eventFactoryImpl struct {
}

func NewEventFactory() core.EventFactory {
	return &eventFactoryImpl{}
}

func (f eventFactoryImpl) NewCommandEvent(ctx context.Context, data interface{}) core.Event {
	return newDefaultEvent(ctx, core.CommandEvent, util.GetTypeName(data), data, "", nil, nil)
}

func (f eventFactoryImpl) NewCommandEventWithCustomType(ctx context.Context, typ string, data interface{}) core.Event {
	return newDefaultEvent(ctx, core.CommandEvent, typ, data, "", nil, nil)
}

func (f eventFactoryImpl) NewBusinessEvent(ctx context.Context, data interface{}) core.Event {
	return newDefaultEvent(ctx, core.BusinessEvent, util.GetTypeName(data), data, "", nil, nil)
}

func (f eventFactoryImpl) NewBusinessEventWithCustomType(ctx context.Context, typ string, data interface{}) core.Event {
	return newDefaultEvent(ctx, core.BusinessEvent, typ, data, "", nil, nil)
}

func (f eventFactoryImpl) NewSystemEvent(ctx context.Context, data interface{}) core.Event {
	return newDefaultEvent(ctx, core.SystemEvent, util.GetTypeName(data), data, "", nil, nil)
}

func (f eventFactoryImpl) NewSystemEventWithCustomType(ctx context.Context, typ string, data interface{}) core.Event {
	return newDefaultEvent(ctx, core.SystemEvent, typ, data, "", nil, nil)
}

func (f eventFactoryImpl) NewDomainEvent(ctx context.Context, data interface{}, entity core.Entity, state interface{}) core.Event {
	return newDefaultEvent(ctx, core.DomainEvent, util.GetTypeName(data), data, entity.ID(), entity, state)
}

func (f eventFactoryImpl) NewDomainEventWithCustomType(ctx context.Context, typ string, data interface{}, entity core.Entity, state interface{}) core.Event {
	return newDefaultEvent(ctx, core.DomainEvent, typ, data, entity.ID(), entity, state)
}

func (f eventFactoryImpl) NewFirstDomainEvent(ctx context.Context, data interface{}, state interface{}) core.Event {
	return newDefaultEvent(ctx, core.DomainEvent, util.GetTypeName(data), data, "", nil, state)
}

func (f eventFactoryImpl) NewFirstDomainEventWithCustomType(ctx context.Context, typ string, data interface{}, state interface{}) core.Event {
	return newDefaultEvent(ctx, core.DomainEvent, typ, data, "", nil, state)
}

func (f eventFactoryImpl) NewFirstDomainEventWithCustomID(ctx context.Context, data interface{}, id string, state interface{}) core.Event {
	return newDefaultEvent(ctx, core.DomainEvent, util.GetTypeName(data), data, id, nil, state)
}

func (f eventFactoryImpl) NewFirstDomainEventWithCustomIDAndCustomType(ctx context.Context, typ string, data interface{}, id string, state interface{}) core.Event {
	return newDefaultEvent(ctx, core.DomainEvent, typ, data, id, nil, state)
}

func newDefaultEvent(ctx context.Context, class core.EventClass, eventType string, data interface{}, id string, entity core.Entity, state interface{}) core.Event {

	trigger := ctx.Value(cevixe.CevixeEventTrigger)
	transaction := ctx.Value(cevixe.CevixeTransaction).(string)
	userID := ctx.Value(cevixe.CevixeUserID).(string)
	loc, _ := time.LoadLocation("America/Lima")
	eventTime := time.Now().In(loc)
	eventTimeStamp := eventTime.UnixNano() / int64(time.Millisecond)

	eventRecord := &model.AwsEventRecord{
		EventClass:  aws.String(string(class)),
		EventType:   aws.String(eventType),
		EventDay:    aws.String(eventTime.Format("2006-01-02")),
		EventTime:   aws.Int64(eventTimeStamp),
		EventAuthor: aws.String(userID),
		EventData:   toGenericData(data),
		Transaction: aws.String(transaction),
	}

	if trigger != nil {
		triggerEvent := trigger.(core.Event)
		eventRecord.TriggerSource = aws.String(triggerEvent.Source())
		eventRecord.TriggerID = aws.String(triggerEvent.ID())
	}

	addEntityMetadata(class, id, entity, state, eventRecord)
	addEventIdentity(class, eventType, entity, state, eventRecord)
	eventRecord = compressEventRecord(ctx, eventRecord)

	return NewEvent(ctx, eventRecord)
}

func addEntityMetadata(class core.EventClass, id string, entity core.Entity, state interface{}, record *model.AwsEventRecord) {

	if class != core.DomainEvent {
		return
	}
	if entity == nil {
		if id == "" {
			record.EntityID = aws.String(uuid.NewString())
		} else {
			record.EntityID = aws.String(id)
		}
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

func addEventIdentity(class core.EventClass, typ string, entity core.Entity, state interface{}, record *model.AwsEventRecord) {

	switch class {
	case core.CommandEvent:
		dataTypeName := strcase.KebabCase(typ)
		record.EventSource = aws.String("/command/" + dataTypeName)
		record.EventID = aws.String(uuid.NewString())
		break
	case core.DomainEvent:
		if entity == nil {
			entityTypeName := strcase.KebabCase(util.GetTypeName(state))
			record.EventSource = aws.String("/domain/" + entityTypeName + "/" + *record.EntityID)
			record.EventID = aws.String(fmt.Sprintf("%020d", 1))
		} else {
			entityTypeName := strcase.KebabCase(entity.Type())
			record.EventSource = aws.String("/domain/" + entityTypeName + "/" + *record.EntityID)
			record.EventID = aws.String(fmt.Sprintf("%020d", entity.Version()+1))
		}
		break
	case core.BusinessEvent:
		dataTypeName := strcase.KebabCase(typ)
		record.EventSource = aws.String("/business/" + dataTypeName)
		record.EventID = aws.String(uuid.NewString())
		break
	case core.SystemEvent:
		dataTypeName := strcase.KebabCase(typ)
		record.EventSource = aws.String("/system/" + dataTypeName)
		record.EventID = aws.String(uuid.NewString())
		break
	default:
		panic(errors.New("event class not supported"))
	}
}

func toGenericData(object interface{}) *map[string]interface{} {
	payload := &map[string]interface{}{}
	objectJson := util.MarshalJson(object)
	util.UnmarshalJson(objectJson, payload)
	return payload
}

const RecordSizeLimit = 960

func compressEventRecord(ctx context.Context, record *model.AwsEventRecord) *model.AwsEventRecord {

	content := json.Marshall(record)
	if len(content) <= RecordSizeLimit {
		return record
	}

	contentEncoding := "gzip"
	contentType := "application/json"
	content = gzip.Compress(content)
	if len(content) <= RecordSizeLimit {
		record.ContentEncoding = contentEncoding
		record.ContentType = contentType
		record.Content = content
		record.EventData = nil
		record.EntityState = nil
		return record
	}

	contentLocation := uuid.NewString()
	awsContext := ctx.Value(AwsContext).(*Context)
	awsContext.AwsObjectStore.SaveRawObject(ctx, contentLocation, content)

	record.ContentLocation = contentLocation
	record.ContentEncoding = contentEncoding
	record.ContentType = contentType
	record.Content = nil
	record.EventData = nil
	record.EntityState = nil

	return record
}
