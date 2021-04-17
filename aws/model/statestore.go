package model

import "context"

type StateStore interface {
	UpdateState(ctx context.Context, events []*EventObject)
}
