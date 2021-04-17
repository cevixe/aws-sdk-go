package dynamo

import (
	"context"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute"
	"github.com/cevixe/aws-sdk-go/aws/model"
	"strconv"
)

type stateStoreImpl struct {
	model.StateStore
	stateStoreRegion string
	stateStoreTable  string
	clientFactory    ClientFactory
}

func (s stateStoreImpl) UpdateState(ctx context.Context, events []*model.EventObject) {

	requests := make([]*dynamodb.WriteRequest, 0, len(events))

	for _, item := range events {
		var payloadMap *map[string]interface{}
		if item.SourceState != nil {
			payloadMap = item.SourceState
		} else {
			payloadMap = item.EventPayload
		}
		payload, err := dynamodbattribute.MarshalMap(payloadMap)
		if err != nil {
			panic(err)
		}

		rq := &dynamodb.WriteRequest{
			PutRequest: &dynamodb.PutRequest{
				Item: map[string]*dynamodb.AttributeValue{
					"type":       {S: aws.String(item.SourceType)},
					"id":         {S: aws.String(item.SourceID)},
					"version":    {N: aws.String(strconv.FormatUint(item.EventID, 10))},
					"payload":    {M: payload},
					"created_at": {N: aws.String(strconv.FormatInt(item.SourceTime, 10))},
					"created_by": {S: aws.String(item.SourceOwner)},
					"updated_at": {N: aws.String(strconv.FormatInt(item.EventTime, 10))},
					"updated_by": {S: aws.String(item.EventAuthor)},
				},
			},
		}
		requests = append(requests, rq)
	}

	input := &dynamodb.BatchWriteItemInput{
		RequestItems: map[string][]*dynamodb.WriteRequest{
			s.stateStoreTable: requests,
		},
	}

	client := s.clientFactory.GetClient(s.stateStoreRegion)
	_, err := client.BatchWriteItemWithContext(ctx, input)
	if err != nil {
		panic(err)
	}
}

func NewStateStore(
	stateStoreRegion string,
	stateStoreTable string,
	clientFactory ClientFactory) model.StateStore {
	return &stateStoreImpl{
		stateStoreRegion: stateStoreRegion,
		stateStoreTable:  stateStoreTable,
		clientFactory:    clientFactory,
	}
}
