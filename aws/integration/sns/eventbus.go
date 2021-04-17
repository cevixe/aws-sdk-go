package sns

import (
	"context"
	"fmt"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/sns"
	"github.com/cevixe/aws-sdk-go/aws/model"
	"github.com/cevixe/aws-sdk-go/util"
)

type eventBusImpl struct {
	model.EventBus
	eventBusRegion string
	eventBusTopic  string
	clientFactory  ClientFactory
}

func (e eventBusImpl) PublishEvent(ctx context.Context, event *model.EventObject) {

	eventJson := util.MarshalJsonString(event)
	messageJson := util.MarshalJsonString(map[string]interface{}{"default": eventJson})

	var input = &sns.PublishInput{
		TopicArn:         aws.String(e.eventBusTopic),
		Message:          aws.String(messageJson),
		MessageStructure: aws.String("json"),
		MessageAttributes: map[string]*sns.MessageAttributeValue{
			"event_type": {
				DataType:    aws.String("String"),
				StringValue: aws.String(event.EventType),
			},
		},
	}

	client := e.clientFactory.GetClient(e.eventBusRegion)
	_, err := client.PublishWithContext(ctx, input)

	if err != nil {
		panic(fmt.Errorf("cannot publish event(/%s/%s/%d)\n%v",
			event.SourceType, event.SourceID, event.EventID, err))
	}
}

func NewEventBus(
	eventBusRegion string,
	eventBusTopic string,
	clientFactory ClientFactory) model.EventBus {
	return &eventBusImpl{
		eventBusRegion: eventBusRegion,
		eventBusTopic:  eventBusTopic,
		clientFactory:  clientFactory,
	}
}
