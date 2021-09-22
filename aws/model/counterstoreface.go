package model

import (
	"context"
)

type AwsCounterStore interface {
	GetValue(ctx context.Context, category string, name string) uint64
	NewValue(ctx context.Context, category string, name string) uint64
}
