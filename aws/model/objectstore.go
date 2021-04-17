package model

import (
	"context"
)

type ObjectStore interface {
	GetObject(ctx context.Context, ref *ObjectRef, v interface{})
	SaveObject(ctx context.Context, key string, v interface{}) *ObjectRef
}
