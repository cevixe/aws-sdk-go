package decorator

import (
	"context"
	"github.com/cevixe/core-sdk-go/core"
)

type EventHandlerDecorator interface {
	Handle(ctx context.Context, event core.Event) core.Event
}
