package model

import (
	"context"
	"time"
)

type AwsEventStore interface {
	CreateControlRecord(ctx context.Context, record *AwsControlRecord)
	CreateUncontrolledEventRecord(ctx context.Context, record *AwsEventRecord)
	CreateControlledEventRecord(ctx context.Context, event *AwsEventRecord, control *AwsControlRecord)

	GetControlRecords(ctx context.Context, group string) []*AwsControlRecord
	GetEventRecordByID(ctx context.Context, source string, id string) *AwsEventRecord
	GetLastEventRecord(ctx context.Context, source string) *AwsEventRecord

	GetEventHeaders(ctx context.Context, source string, after *string, before *string, nextToken *string, limit *int64) *AwsEventHeaderRecordPage
	GetSourceEvents(ctx context.Context, source string, after *time.Time, before *time.Time, nextToken *string, limit *int64) *AwsEventRecordPage
	GetDayEvents(ctx context.Context, day string, after *time.Time, before *time.Time, nextToken *string, limit *int64) *AwsEventRecordPage
	GetTypeEvents(ctx context.Context, typ string, after *time.Time, before *time.Time, nextToken *string, limit *int64) *AwsEventRecordPage
	GetAuthorEvents(ctx context.Context, author string, after *time.Time, before *time.Time, nextToken *string, limit *int64) *AwsEventRecordPage
	GetTransactionEvents(ctx context.Context, transaction string, after *time.Time, before *time.Time, nextToken *string, limit *int64) *AwsEventRecordPage
}
