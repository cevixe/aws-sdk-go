package model

import "context"

type AwsObjectStore interface {
	GetRawObject(ctx context.Context, id string) []byte
	SaveRawObject(ctx context.Context, id string, buffer []byte)
	GetJsonObject(ctx context.Context, id string, v interface{})
	SaveJsonObject(ctx context.Context, id string, v interface{})
}
