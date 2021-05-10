package model

import "context"

type AwsStateStore interface {
	UpdateState(ctx context.Context, state *AwsStateRecord)
	UpdateStates(ctx context.Context, state []*AwsStateRecord)
}
