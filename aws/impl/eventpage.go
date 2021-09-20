package impl

import (
	"context"
	"github.com/cevixe/aws-sdk-go/aws/model"
	"github.com/cevixe/core-sdk-go/core"
)

type EventPageImpl struct {
	Ctx    context.Context
	Record *model.AwsEventRecordPage
}

func (e EventPageImpl) Items() []core.Event {
	items := make([]core.Event, 0)
	for idx, _ := range e.Record.Items {
		items = append(items, NewEvent(e.Ctx, e.Record.Items[idx]))
	}
	return items
}

func (e EventPageImpl) NextToken() string {
	if e.Record.NextToken == nil {
		return ""
	}
	return *e.Record.NextToken
}

func NewEventPage(ctx context.Context, record *model.AwsEventRecordPage) core.EventPage {
	return &EventPageImpl{
		Ctx:    ctx,
		Record: record,
	}
}
