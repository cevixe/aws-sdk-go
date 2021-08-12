package sns

import (
	"context"
	"fmt"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/sns"
	"github.com/aws/aws-sdk-go/service/sns/snsiface"
	"github.com/cevixe/aws-sdk-go/aws/env"
	"github.com/cevixe/aws-sdk-go/aws/factory"
	"github.com/cevixe/aws-sdk-go/aws/model"
	util2 "github.com/cevixe/aws-sdk-go/aws/util"
	"os"
)

type eventBusImpl struct {
	eventBusTopic string
	snsClient     snsiface.SNSAPI
}

func NewSnsEventBus(
	eventBusTopic string,
	snsClient snsiface.SNSAPI) model.AwsEventBus {

	return &eventBusImpl{
		eventBusTopic: eventBusTopic,
		snsClient:     snsClient,
	}
}

func NewDefaultSnsEventBus(awsFactory factory.AwsFactory) model.AwsEventBus {

	eventBusTopicArn := os.Getenv(env.CevixeEventBusTopicArn)
	snsClient := awsFactory.SnsClient()

	return NewSnsEventBus(eventBusTopicArn, snsClient)
}

func (e eventBusImpl) PublishEvent(ctx context.Context, event *model.AwsEventRecord) {

	messageJson := util2.MarshalJsonString(map[string]interface{}{
		"default": util2.MarshalJsonString(event),
	})

	var input = &sns.PublishInput{
		TopicArn:         aws.String(e.eventBusTopic),
		Message:          aws.String(messageJson),
		MessageStructure: aws.String("json"),
		MessageAttributes: map[string]*sns.MessageAttributeValue{
			"event_type": {
				DataType:    aws.String("String"),
				StringValue: event.EventType,
			},
			"event_class": {
				DataType:    aws.String("String"),
				StringValue: event.EventClass,
			},
		},
	}

	_, err := e.snsClient.PublishWithContext(ctx, input)
	if err != nil {
		panic(fmt.Errorf("cannot publish event to sns\n%v", err))
	}
}
