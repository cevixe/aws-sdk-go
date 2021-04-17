package sqs

import (
	"context"
	"fmt"
	"github.com/aws/aws-lambda-go/events"
	"github.com/cevixe/aws-sdk-go/aws/impl"
	"github.com/cevixe/aws-sdk-go/aws/model"
	"github.com/cevixe/aws-sdk-go/util"
	"github.com/cevixe/core-sdk-go/core"
)

func MapSQSEventToSingleCevixeEvent(ctx context.Context, sqsEvent events.SQSEvent) core.Event {

	if len(sqsEvent.Records) > 1 {
		panic(fmt.Errorf("sqs event with multiple sns records not supported"))
	}

	sqsMessage := sqsEvent.Records[0]
	snsMessage := &events.SNSEntity{}
	eventValue := &model.EventObject{}
	util.UnmarshalJsonString(sqsMessage.Body, snsMessage)
	util.UnmarshalJsonString(snsMessage.Message, eventValue)

	return impl.NewEvent(ctx, eventValue)
}
