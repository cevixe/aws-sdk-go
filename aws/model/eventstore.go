package model

import (
	"context"
)

type EventStore interface {
	GetEventById(ctx context.Context, source string, id uint64) *EventObject
	GetLatestEvent(ctx context.Context, source string) *EventObject
	SaveEvent(ctx context.Context, event *EventObject)
}
