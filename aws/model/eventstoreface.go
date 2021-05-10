package model

import "context"

type AwsEventStore interface {
	CreateControlRecord(ctx context.Context, record *AwsControlRecord)
	CreateUncontrolledEventRecord(ctx context.Context, record *AwsEventRecord)
	CreateControlledEventRecord(ctx context.Context, event *AwsEventRecord, control *AwsControlRecord)

	GetControlRecords(ctx context.Context, group string) []*AwsControlRecord
	GetEventRecordByID(ctx context.Context, source string, id string) *AwsEventRecord
	GetLastEventRecord(ctx context.Context, source string) *AwsEventRecord
}
