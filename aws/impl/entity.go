package impl

import (
	"context"
	"fmt"
	"github.com/cevixe/aws-sdk-go/aws/model"
	"github.com/cevixe/aws-sdk-go/aws/util"
	"github.com/cevixe/core-sdk-go/core"
	"github.com/pkg/errors"
	"strconv"
	"time"
)

type EntityImpl struct {
	Context     context.Context
	EventRecord *model.AwsEventRecord
	StateRecord *model.AwsStateRecord
}

func NewEntity(ctx context.Context,
	stateRecord *model.AwsStateRecord,
	eventRecord *model.AwsEventRecord) core.Entity {

	if eventRecord != nil {
		if eventRecord.EntityDeleted {
			return nil
		}
		entityVersion, err := strconv.ParseUint(*eventRecord.EventID, 10, 64)
		if err != nil {
			panic(errors.Wrap(err, "cannot get entity version"))
		}
		entityState := make(map[string]interface{})
		if eventRecord.EntityState != nil {
			entityState = *eventRecord.EntityState
		}
		return &EntityImpl{
			Context: ctx,
			StateRecord: &model.AwsStateRecord{
				Type:            *eventRecord.EntityType,
				ID:              *eventRecord.EntityID,
				Version:         entityVersion,
				Deleted:         eventRecord.EntityDeleted,
				State:           entityState,
				UpdatedAt:       *eventRecord.EventTime,
				UpdatedBy:       *eventRecord.EventAuthor,
				CreatedAt:       *eventRecord.EntityCreatedAt,
				CreatedBy:       *eventRecord.EntityCreatedBy,
				ContentLocation: eventRecord.ContentLocation,
				ContentType:     eventRecord.ContentType,
				ContentEncoding: eventRecord.ContentEncoding,
				Content:         eventRecord.Content,
			},
			EventRecord: nil,
		}
	}

	if stateRecord != nil {
		return &EntityImpl{
			Context:     ctx,
			StateRecord: stateRecord,
			EventRecord: nil,
		}
	}

	return nil
}

func (e EntityImpl) ID() string {
	return e.StateRecord.ID
}

func (e EntityImpl) Type() string {
	return e.StateRecord.Type
}

func (e EntityImpl) Version() uint64 {
	return e.StateRecord.Version
}

func (e *EntityImpl) State(v interface{}) {
	if e.StateRecord.State != nil {
		fmt.Println("State case")
		json := util.MarshalJson(e.StateRecord.State)
		util.UnmarshalJson(json, v)
	} else {
		if e.EventRecord == nil {
			fmt.Println("State store case")
			e.EventRecord = GetEventContent(
				e.Context,
				e.StateRecord.ContentLocation,
				e.StateRecord.ContentEncoding,
				e.StateRecord.ContentType,
				e.StateRecord.Content)
		} else if e.EventRecord.EventData == nil {
			fmt.Println("Event record case")
			e.EventRecord = GetEventContent(
				e.Context,
				e.EventRecord.ContentLocation,
				e.EventRecord.ContentEncoding,
				e.EventRecord.ContentType,
				e.EventRecord.Content)
		}
		json := util.MarshalJson(e.EventRecord.EntityState)
		util.UnmarshalJson(json, v)
	}
}

func (e EntityImpl) UpdatedAt() time.Time {
	nanoseconds := e.StateRecord.UpdatedAt * int64(time.Millisecond)
	return time.Unix(0, nanoseconds)
}

func (e EntityImpl) UpdatedBy() string {
	return e.StateRecord.UpdatedBy
}

func (e EntityImpl) CreatedAt() time.Time {
	nanoseconds := e.StateRecord.CreatedAt * int64(time.Millisecond)
	return time.Unix(0, nanoseconds)
}

func (e EntityImpl) CreatedBy() string {
	return e.StateRecord.CreatedBy
}
