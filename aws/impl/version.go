package impl

import (
	"github.com/cevixe/aws-sdk-go/aws/model"
	"github.com/cevixe/core-sdk-go/core"
	"time"
)

type VersionImpl struct {
	Record *model.AwsEventHeaderRecord
}

func (v VersionImpl) ID() uint64 {
	panic("implement me")
}

func (v VersionImpl) Time() time.Time {
	panic("implement me")
}

func NewVersion(record *model.AwsEventHeaderRecord) core.Version {
	return &VersionImpl{Record: record}
}
