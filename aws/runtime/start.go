package runtime

import (
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/cevixe/aws-sdk-go/aws/runtime/delivery"
	"github.com/cevixe/core-sdk-go/core"
)

func Start(fn core.EventHandler) {
	StartAtLeastOnce(fn)
}

func StartAtMostOnce(fn core.EventHandler) {
	ctx := NewContext()
	fn = delivery.AtMostOnce(ctx, fn)
	lambda.StartWithContext(ctx, NewHandler(fn))
}

func StartAtLeastOnce(fn core.EventHandler) {
	ctx := NewContext()
	fn = delivery.AtLeastOnce(ctx, fn)
	lambda.StartWithContext(ctx, NewHandler(fn))
}

func StartExactlyOnce(fn core.EventHandler) {
	ctx := NewContext()
	fn = delivery.ExactlyOnce(ctx, fn)
	lambda.StartWithContext(ctx, NewHandler(fn))
}

func StartFullExactlyOnce(fn core.EventHandler) {
	ctx := NewContext()
	fn = delivery.FullExactlyOnce(ctx, fn)
	lambda.StartWithContext(ctx, NewHandler(fn))
}
