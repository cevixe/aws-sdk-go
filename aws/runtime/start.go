package runtime

import (
	"context"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/cevixe/aws-sdk-go/aws/runtime/delivery"
	"github.com/cevixe/core-sdk-go/core"
)

func Start(fn core.EventHandler) {
	ctx := NewContext()
	StartAtLeastOnce(ctx, fn)
}

func StartWithContext(ctx context.Context, fn core.EventHandler) {
	StartAtLeastOnce(ctx, fn)
}

func StartAtMostOnce(ctx context.Context, fn core.EventHandler) {
	fn = delivery.AtMostOnce(ctx, fn)
	lambda.StartWithContext(ctx, NewHandler(fn))
}

func StartAtLeastOnce(ctx context.Context, fn core.EventHandler) {
	fn = delivery.AtLeastOnce(ctx, fn)
	lambda.StartWithContext(ctx, NewHandler(fn))
}

func StartExactlyOnce(ctx context.Context, fn core.EventHandler) {
	fn = delivery.ExactlyOnce(ctx, fn)
	lambda.StartWithContext(ctx, NewHandler(fn))
}

func StartFullExactlyOnce(ctx context.Context, fn core.EventHandler) {
	fn = delivery.FullExactlyOnce(ctx, fn)
	lambda.StartWithContext(ctx, NewHandler(fn))
}
