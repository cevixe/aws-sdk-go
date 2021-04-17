package model

import "context"

type EventBus interface {
	PublishEvent(ctx context.Context, event *EventObject)
}
