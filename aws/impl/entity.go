package impl

import (
	"context"
	"github.com/cevixe/aws-sdk-go/aws/model"
	"github.com/cevixe/aws-sdk-go/util"
	"github.com/cevixe/core-sdk-go/core"
	"time"
)

type entityImpl struct {
	core.Entity
	ctx       context.Context
	lastEvent *model.EventObject
}

func (e *entityImpl) ID() string {
	return e.lastEvent.SourceID
}

func (e *entityImpl) Type() string {
	return e.lastEvent.SourceType
}

func (e *entityImpl) Time() time.Time {
	nanoseconds := e.lastEvent.SourceTime * int64(time.Millisecond)
	return time.Unix(0, nanoseconds)
}

func (e *entityImpl) Owner() string {
	return e.lastEvent.SourceOwner
}

func (e *entityImpl) Version() uint64 {
	return e.lastEvent.EventID
}

func (e *entityImpl) State(v interface{}) {
	if e.lastEvent.SourceState == nil && e.lastEvent.EventPayload == nil && e.lastEvent.Reference != nil {
		awsContext := e.ctx.Value(model.AwsContext).(*model.Context)
		newValue := &model.EventObject{}
		awsContext.AwsObjectStore.GetObject(e.ctx, e.lastEvent.Reference, newValue)
		e.lastEvent = newValue
	}

	var json []byte
	if e.lastEvent.SourceState == nil {
		json = util.MarshalJson(e.lastEvent.EventPayload)
	} else {
		json = util.MarshalJson(e.lastEvent.SourceState)
	}
	util.UnmarshalJson(json, v)
}
