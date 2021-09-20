package impl

import (
	"context"
	"fmt"
	"github.com/cevixe/aws-sdk-go/aws/model"
	"github.com/cevixe/core-sdk-go/core"
	"github.com/stoewer/go-strcase"
	"time"
)

type eventStoreImpl struct {
	eventStore model.AwsEventStore
}

func NewEventStore(eventStore model.AwsEventStore) core.EventStore {
	return &eventStoreImpl{eventStore: eventStore}
}

func (e eventStoreImpl) GetLastEvent(ctx context.Context, source string) core.Event {
	eventValue := e.eventStore.GetLastEventRecord(ctx, source)
	return NewEvent(ctx, eventValue)
}

func (e eventStoreImpl) GetEventByID(ctx context.Context, source string, id string) core.Event {
	eventValue := e.eventStore.GetEventRecordByID(ctx, source, id)
	return NewEvent(ctx, eventValue)
}

func (e eventStoreImpl) GetEntityEvents(ctx context.Context, typ string, id string,
	after *time.Time, before *time.Time, nextToken *string, limit *int64) core.EventPage {

	source := fmt.Sprintf("/domain/%s/%s", strcase.KebabCase(typ), id)
	page := e.eventStore.GetSourceEvents(ctx, source, after, before, nextToken, limit)
	return NewEventPage(ctx, page)
}

func (e eventStoreImpl) GetDayEvents(ctx context.Context, day string,
	after *time.Time, before *time.Time, nextToken *string, limit *int64) core.EventPage {

	page := e.eventStore.GetDayEvents(ctx, day, after, before, nextToken, limit)
	return NewEventPage(ctx, page)
}

func (e eventStoreImpl) GetTypeEvents(ctx context.Context, typ string,
	after *time.Time, before *time.Time, nextToken *string, limit *int64) core.EventPage {

	page := e.eventStore.GetTypeEvents(ctx, typ, after, before, nextToken, limit)
	return NewEventPage(ctx, page)
}

func (e eventStoreImpl) GetAuthorEvents(ctx context.Context, author string,
	after *time.Time, before *time.Time, nextToken *string, limit *int64) core.EventPage {

	page := e.eventStore.GetAuthorEvents(ctx, author, after, before, nextToken, limit)
	return NewEventPage(ctx, page)
}

func (e eventStoreImpl) GetTransactionEvents(ctx context.Context, transaction string,
	after *time.Time, before *time.Time, nextToken *string, limit *int64) core.EventPage {

	page := e.eventStore.GetTransactionEvents(ctx, transaction, after, before, nextToken, limit)
	return NewEventPage(ctx, page)
}
