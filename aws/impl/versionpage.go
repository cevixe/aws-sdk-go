package impl

import (
	"github.com/cevixe/aws-sdk-go/aws/model"
	"github.com/cevixe/core-sdk-go/core"
)

type VersionPageImpl struct {
	Record *model.AwsEventHeaderRecordPage
}

func (e VersionPageImpl) Items() []core.Version {
	items := make([]core.Version, 0)
	for idx, _ := range e.Record.Items {
		items = append(items, NewVersion(e.Record.Items[idx]))
	}
	return items
}

func (e VersionPageImpl) NextToken() string {
	if e.Record.NextToken == nil {
		return ""
	}
	return *e.Record.NextToken
}

func NewVersionPage(record *model.AwsEventHeaderRecordPage) core.VersionPage {
	return &VersionPageImpl{Record: record}
}
