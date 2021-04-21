package dynamo

import (
	"context"
	"fmt"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute"
	"github.com/cevixe/aws-sdk-go/aws/model"
	"strconv"
)

type eventStoreImpl struct {
	model.EventStore
	eventStoreRegion string
	eventStoreTable  string
	clientFactory    ClientFactory
}

func (r eventStoreImpl) GetLatestEvent(ctx context.Context, source string) *model.EventObject {

	input := &dynamodb.QueryInput{
		TableName:              aws.String(r.eventStoreTable),
		Limit:                  aws.Int64(1),
		ScanIndexForward:       aws.Bool(false),
		KeyConditionExpression: aws.String("#source = :source"),
		ExpressionAttributeNames: map[string]*string{
			"#source": aws.String("source"),
		},
		ExpressionAttributeValues: map[string]*dynamodb.AttributeValue{
			":source": {S: aws.String(source)},
		},
	}

	client := r.clientFactory.GetClient(r.eventStoreRegion)
	output, err := client.QueryWithContext(ctx, input)

	if err != nil {
		panic(fmt.Errorf("cannot get latest event from dynamodb(Source:%s)\n%v", source, err))
	}

	if len(output.Items) == 0 {
		return nil
	}

	event := &model.EventObject{}
	fromDynamoItem(output.Items[0], event)
	return event
}

func (r eventStoreImpl) GetEvent(ctx context.Context, source string, id uint64) *model.EventObject {

	input := &dynamodb.GetItemInput{
		TableName: aws.String(r.eventStoreTable),
		Key: map[string]*dynamodb.AttributeValue{
			"source": {S: aws.String(source)},
			"id":     {N: aws.String(strconv.FormatUint(id, 10))},
		},
	}

	client := r.clientFactory.GetClient(r.eventStoreRegion)
	output, err := client.GetItemWithContext(ctx, input)

	if err != nil {
		panic(fmt.Errorf("cannot get event from dynamodb(Source:%s, ID:%d)\n%v", source, id, err))
	}

	if output.Item == nil {
		return nil
	}

	event := &model.EventObject{}
	fromDynamoItem(output.Item, event)
	return event
}

func (r eventStoreImpl) SaveEvent(ctx context.Context, event *model.EventObject) {

	input := &dynamodb.PutItemInput{
		TableName: aws.String(r.eventStoreTable),
		Item:      toDynamoItem(event),
	}

	client := r.clientFactory.GetClient(r.eventStoreRegion)
	_, err := client.PutItemWithContext(ctx, input)

	if err != nil {
		panic(fmt.Errorf("cannot put event to dynamodb(/%s/%s/%d)\n%v", event.SourceType, event.SourceID, event.EventID, err))
	}
}

func fromDynamoItem(item map[string]*dynamodb.AttributeValue, event *model.EventObject) {
	err := dynamodbattribute.UnmarshalMap(item, event)

	if err != nil {
		panic(fmt.Errorf("cannot unmarshal dynamodb item\n%v", err))
	}
}

func toDynamoItem(event *model.EventObject) map[string]*dynamodb.AttributeValue {
	output, err := dynamodbattribute.MarshalMap(event)
	output["source_key"] = &dynamodb.AttributeValue{
		S: aws.String("/" + event.SourceType + "/" + event.SourceID),
	}

	if err != nil {
		panic(fmt.Errorf("cannot marshal dynamodb item\n%v", err))
	}

	return output
}

func NewEventStore(eventStoreRegion string, eventStoreTable string, clientFactory ClientFactory) model.EventStore {
	return &eventStoreImpl{
		eventStoreRegion: eventStoreRegion,
		eventStoreTable:  eventStoreTable,
		clientFactory:    clientFactory,
	}
}
