package dynamodb

import (
	"context"
	"fmt"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbiface"
	"github.com/cevixe/aws-sdk-go/aws/env"
	"github.com/cevixe/aws-sdk-go/aws/factory"
	"github.com/cevixe/aws-sdk-go/aws/model"
	"os"
)

type stateStoreImpl struct {
	stateStoreTable string
	dynamodbClient  dynamodbiface.DynamoDBAPI
}

func NewDynamodbStateStore(
	stateStoreTable string,
	dynamodbClient dynamodbiface.DynamoDBAPI) model.AwsStateStore {

	return &stateStoreImpl{
		stateStoreTable: stateStoreTable,
		dynamodbClient:  dynamodbClient,
	}
}

func NewDefaultDynamodbStateStore(awsFactory factory.AwsFactory) model.AwsStateStore {

	stateStoreTableName := os.Getenv(env.CevixeStateStoreTableName)
	dynamodbClient := awsFactory.DynamodbClient()

	return NewDynamodbStateStore(stateStoreTableName, dynamodbClient)
}

func (s stateStoreImpl) UpdateState(ctx context.Context, state *model.AwsStateRecord) {
	if state.State == nil {
		input := &dynamodb.DeleteItemInput{
			TableName: aws.String(s.stateStoreTable),
			Key: map[string]*dynamodb.AttributeValue{
				"type": {S: aws.String(state.Type)},
				"id":   {S: aws.String(state.ID)},
			},
		}

		_, err := s.dynamodbClient.DeleteItemWithContext(ctx, input)
		if err != nil {
			panic(fmt.Errorf("cannot update state record\n%v", err))
		}
	} else {
		input := &dynamodb.PutItemInput{
			TableName: aws.String(s.stateStoreTable),
			Item:      MarshallDynamodbItem(state),
		}

		_, err := s.dynamodbClient.PutItemWithContext(ctx, input)
		if err != nil {
			panic(fmt.Errorf("cannot update state record\n%v", err))
		}
	}
}

func (s stateStoreImpl) UpdateStates(ctx context.Context, states []*model.AwsStateRecord) {
	requests := make([]*dynamodb.WriteRequest, 0, len(states))
	for _, item := range states {
		var elem *dynamodb.WriteRequest
		if item.State != nil {
			elem = &dynamodb.WriteRequest{
				PutRequest: &dynamodb.PutRequest{
					Item: MarshallDynamodbItem(item),
				},
			}
		} else {
			elem = &dynamodb.WriteRequest{
				DeleteRequest: &dynamodb.DeleteRequest{
					Key: map[string]*dynamodb.AttributeValue{
						"type": {S: aws.String(item.Type)},
						"id":   {S: aws.String(item.ID)},
					},
				},
			}
		}

		requests = append(requests, elem)
	}
	input := &dynamodb.BatchWriteItemInput{
		RequestItems: map[string][]*dynamodb.WriteRequest{
			s.stateStoreTable: requests,
		},
	}

	_, err := s.dynamodbClient.BatchWriteItemWithContext(ctx, input)
	if err != nil {
		panic(fmt.Errorf("cannot update state records\n%v", err))
	}
}
