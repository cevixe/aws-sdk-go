package lambda

import (
	"context"
	"github.com/aws/aws-lambda-go/events"
)

type AWSLambdaHandlerFn func(ctx context.Context, event events.SQSEvent) error
