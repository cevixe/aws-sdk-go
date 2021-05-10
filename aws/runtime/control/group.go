package control

import (
	"context"
	"github.com/cevixe/aws-sdk-go/aws/impl"
	"github.com/cevixe/core-sdk-go/core"
)

func GetControlGroup(ctx context.Context, trigger core.Event) string {

	cvxContext := ctx.Value(impl.AwsContext).(*impl.Context)

	handlerUniqueID := "/" + cvxContext.AwsHandlerID + "/" + cvxContext.AwsHandlerVersion
	eventUniqueID := trigger.Source() + "/" + trigger.ID()

	return handlerUniqueID + eventUniqueID
}
