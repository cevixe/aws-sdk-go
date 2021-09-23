package impl

import (
	"github.com/cevixe/aws-sdk-go/aws/model"
	"github.com/cevixe/core-sdk-go/core"
	"github.com/pkg/errors"
	"strconv"
	"time"
)

type VersionImpl struct {
	Record *model.AwsEventHeaderRecord
}

func (v VersionImpl) ID() uint64 {
	entityVersion, err := strconv.ParseUint(*v.Record.EventID, 10, 64)
	if err != nil {
		panic(errors.Wrap(err, "cannot get entity version"))
	}
	return entityVersion
}

func (v VersionImpl) Time() time.Time {
	timeStamp := *v.Record.EventTime * int64(time.Millisecond)
	return time.Unix(0, timeStamp)
}

func (v VersionImpl) Author() string {
	return *v.Record.EventAuthor
}

func NewVersion(record *model.AwsEventHeaderRecord) core.Version {
	if record.EntityDeleted {
		return nil
	}
	return &VersionImpl{Record: record}
}
