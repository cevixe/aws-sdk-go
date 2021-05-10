package impl

import (
	"fmt"
	"github.com/cevixe/aws-sdk-go/aws/model"
	"github.com/cevixe/aws-sdk-go/aws/util"
	"strconv"
	"time"
)

type EntityImpl struct {
	LastEvent *model.AwsEventRecord
}

func (e EntityImpl) ID() string {
	return *e.LastEvent.EntityID
}

func (e EntityImpl) Type() string {
	return *e.LastEvent.EntityType
}

func (e EntityImpl) Version() uint64 {
	number, err := strconv.ParseUint(*e.LastEvent.EventID, 10, 64)
	if err != nil {
		panic(fmt.Errorf("cannot get entity version\n%v", err))
	}
	return number
}

func (e *EntityImpl) State(v interface{}) {
	json := util.MarshalJson(e.LastEvent.EntityState)
	util.UnmarshalJson(json, v)
}

func (e EntityImpl) UpdatedAt() time.Time {
	nanoseconds := *e.LastEvent.EventTime * int64(time.Millisecond)
	return time.Unix(0, nanoseconds)
}

func (e EntityImpl) UpdatedBy() string {
	return *e.LastEvent.EventAuthor
}

func (e EntityImpl) CreatedAt() time.Time {
	nanoseconds := *e.LastEvent.EntityCreatedAt * int64(time.Millisecond)
	return time.Unix(0, nanoseconds)
}

func (e EntityImpl) CreatedBy() string {
	return *e.LastEvent.EntityCreatedBy
}
