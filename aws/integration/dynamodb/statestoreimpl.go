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
	"github.com/pkg/errors"
	"os"
	"strconv"
	"time"
)

type stateStoreImpl struct {
	eventStoreTableName   string
	stateStoreTableName   string
	stateStoreIndexByTime string
	dynamodbClient        dynamodbiface.DynamoDBAPI
}

func NewDynamodbStateStore(
	eventStoreTableName string,
	stateStoreTableName string,
	stateStoreIndexByTime string,
	dynamodbClient dynamodbiface.DynamoDBAPI) model.AwsStateStore {

	return &stateStoreImpl{
		eventStoreTableName:   eventStoreTableName,
		stateStoreTableName:   stateStoreTableName,
		stateStoreIndexByTime: stateStoreIndexByTime,
		dynamodbClient:        dynamodbClient,
	}
}

func NewDefaultDynamodbStateStore(awsFactory factory.AwsFactory) model.AwsStateStore {

	eventStoreTableName := os.Getenv(env.CevixeEventStoreTableName)
	stateStoreTableName := os.Getenv(env.CevixeStateStoreTableName)
	stateStoreIndexByTime := os.Getenv(env.CevixeStateStoreIndexByTime)
	dynamodbClient := awsFactory.DynamodbClient()

	return NewDynamodbStateStore(eventStoreTableName, stateStoreTableName, stateStoreIndexByTime, dynamodbClient)
}

func (s stateStoreImpl) UpdateState(ctx context.Context, state *model.AwsStateRecord) {
	if state.State == nil {
		input := &dynamodb.DeleteItemInput{
			TableName: aws.String(s.stateStoreTableName),
			Key: map[string]*dynamodb.AttributeValue{
				"type": {S: aws.String(state.Type)},
				"id":   {S: aws.String(state.ID)},
			},
		}

		_, err := s.dynamodbClient.DeleteItemWithContext(ctx, input)
		if err != nil {
			panic(errors.Wrap(err, "cannot update state record"))
		}
	} else {
		input := &dynamodb.PutItemInput{
			TableName: aws.String(s.stateStoreTableName),
			Item:      MarshallDynamodbItem(state),
		}

		_, err := s.dynamodbClient.PutItemWithContext(ctx, input)
		if err != nil {
			panic(errors.Wrap(err, "cannot update state record"))
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
			s.stateStoreTableName: requests,
		},
	}

	_, err := s.dynamodbClient.BatchWriteItemWithContext(ctx, input)
	if err != nil {
		panic(errors.Wrap(err, "cannot update state records"))
	}
}

func (s stateStoreImpl) GetStates(ctx context.Context, typ string, after *time.Time, nextToken *string, limit *int64) *model.AwsStateRecordPage {

	var afterTime time.Time
	if after != nil {
		afterTime = *after
	}
	unixTime := afterTime.Unix() / int64(time.Millisecond)
	selectStatement := fmt.Sprintf(
		"SELECT * FROM %s.%s WHERE type = ? AND updated_at >= ? ORDER BY updated_at DESC LIMIT %d",
		s.stateStoreTableName, s.stateStoreIndexByTime, *FixPaginationLimit(limit))

	params := &dynamodb.ExecuteStatementInput{
		Statement: aws.String(selectStatement),
		NextToken: nextToken,
		Parameters: []*dynamodb.AttributeValue{
			{S: aws.String(typ)},
			{N: aws.String(strconv.FormatInt(unixTime, 64))},
		},
	}

	output, err := s.dynamodbClient.ExecuteStatementWithContext(ctx, params)
	if err != nil {
		panic(errors.Wrapf(err, "cannot get state records"))
	}

	if len(output.Items) == 0 {
		return &model.AwsStateRecordPage{
			Items:     make([]*model.AwsStateRecord, 0),
			NextToken: output.NextToken,
		}
	}

	records := make([]*model.AwsStateRecord, 0)
	UnmarshallDynamodbItemList(output.Items, &records)

	return &model.AwsStateRecordPage{
		Items:     records,
		NextToken: output.NextToken,
	}
}
