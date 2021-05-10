package control

import (
	"context"
	"github.com/cevixe/aws-sdk-go/aws/impl"
	"github.com/cevixe/aws-sdk-go/aws/model"
	"github.com/cevixe/core-sdk-go/core"
	"strconv"
	"time"
)

func NewControlRecord(ctx context.Context, trigger core.Event, controlType model.AwsControlType, controlIntent uint64) *model.AwsControlRecord {

	cvxContext := ctx.Value(impl.AwsContext).(*impl.Context)

	controlGroup := GetControlGroup(ctx, trigger)
	controlID := string(controlType) + strconv.FormatUint(controlIntent, 10)
	controlTime := time.Now().UnixNano() / int64(time.Millisecond)

	return &model.AwsControlRecord{
		ControlGroup:   controlGroup,
		ControlID:      controlID,
		ControlIntent:  controlIntent,
		ControlType:    controlType,
		ControlTime:    controlTime,
		HandlerID:      cvxContext.AwsHandlerID,
		HandlerVersion: cvxContext.AwsHandlerVersion,
		HandlerTimeout: cvxContext.AwsHandlerTimeout,
		EventSource:    trigger.Source(),
		EventID:        trigger.ID(),
		Transaction:    trigger.Transaction(),
	}
}
