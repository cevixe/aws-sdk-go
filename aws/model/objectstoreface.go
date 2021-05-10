package model

import "context"

type AwsObjectStore interface {
	GetObject(ctx context.Context, reference *AwsObjectStoreReference, v interface{})
	SaveObject(ctx context.Context, id string, v interface{}) *AwsObjectStoreReference
}
