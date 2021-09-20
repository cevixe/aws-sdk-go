package impl

import (
	"context"
	"github.com/cevixe/aws-sdk-go/aws/model"
	"github.com/cevixe/core-sdk-go/core"
)

type EntityPageImpl struct {
	Ctx    context.Context
	Record *model.AwsStateRecordPage
}

func (e EntityPageImpl) Items() []core.Entity {
	items := make([]core.Entity, 0)
	for idx, _ := range e.Record.Items {
		items = append(items, NewEntity(e.Ctx, e.Record.Items[idx], nil))
	}
	return items
}

func (e EntityPageImpl) NextToken() string {
	if e.Record.NextToken == nil {
		return ""
	}
	return *e.Record.NextToken
}

func NewEntityPage(ctx context.Context, record *model.AwsStateRecordPage) core.EntityPage {
	return &EntityPageImpl{
		Ctx:    ctx,
		Record: record,
	}
}
