package model

import "context"

type AwsEventBus interface {
	PublishEvent(ctx context.Context, event *AwsEventRecord)
}
