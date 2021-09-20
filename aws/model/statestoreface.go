package model

import (
	"context"
	"time"
)

type AwsStateStore interface {
	UpdateState(ctx context.Context, state *AwsStateRecord)
	UpdateStates(ctx context.Context, state []*AwsStateRecord)
	GetStates(ctx context.Context, typ string, after *time.Time, nextToken *string, limit *int64) *AwsStateRecordPage
}
