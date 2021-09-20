package impl

import (
	"context"
	"fmt"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/cevixe/aws-sdk-go/aws/model"
	"github.com/cevixe/core-sdk-go/core"
	"github.com/stoewer/go-strcase"
	"time"
)

type stateStoreImpl struct {
	eventStore model.AwsEventStore
	stateStore model.AwsStateStore
}

func (s stateStoreImpl) GetLastVersion(ctx context.Context, typ string, id string) core.Entity {
	source := fmt.Sprintf("/domain/%s/%s", strcase.KebabCase(typ), id)
	event := s.eventStore.GetLastEventRecord(ctx, source)
	if event == nil {
		return nil
	}
	return NewEvent(ctx, event).Entity()
}

func (s stateStoreImpl) GetByVersion(ctx context.Context, typ string, id string, version uint64) core.Entity {
	source := fmt.Sprintf("/domain/%s/%s", strcase.KebabCase(typ), id)
	event := s.eventStore.GetEventRecordByID(ctx, source, fmt.Sprintf("%020d", version))
	if event == nil {
		return nil
	}
	return NewEvent(ctx, event).Entity()
}

func (s stateStoreImpl) GetVersions(ctx context.Context, typ string, id string,
	after *uint64, before *uint64, nextToken *string, limit *int64) core.VersionPage {

	source := fmt.Sprintf("/domain/%s/%s", strcase.KebabCase(typ), id)

	var afterToken *string
	if after != nil {
		afterToken = aws.String(fmt.Sprintf("%020d", *after))
	}

	var beforeToken *string
	if before != nil {
		beforeToken = aws.String(fmt.Sprintf("%020d", *before))
	}

	page := s.eventStore.GetEventHeaders(ctx, source, afterToken, beforeToken, nextToken, limit)
	return NewVersionPage(page)
}

func (s stateStoreImpl) GetByType(ctx context.Context, typ string,
	after *time.Time, nextToken *string, limit *int64) core.EntityPage {

	page := s.stateStore.GetStates(ctx, typ, after, nextToken, limit)
	return NewEntityPage(ctx, page)
}

func NewStateStore(eventStore model.AwsEventStore, stateStore model.AwsStateStore) core.StateStore {
	return &stateStoreImpl{
		eventStore: eventStore,
		stateStore: stateStore,
	}
}
