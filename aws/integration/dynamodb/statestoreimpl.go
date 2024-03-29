package dynamodb

import (
	"context"
	"encoding/base64"
	"fmt"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbiface"
	"github.com/cevixe/aws-sdk-go/aws/env"
	"github.com/cevixe/aws-sdk-go/aws/factory"
	"github.com/cevixe/aws-sdk-go/aws/model"
	"github.com/pkg/errors"
	"math"
	"os"
	"strconv"
	"strings"
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
	if state.Deleted {
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
		if item.Deleted {
			elem = &dynamodb.WriteRequest{
				DeleteRequest: &dynamodb.DeleteRequest{
					Key: map[string]*dynamodb.AttributeValue{
						"type": {S: aws.String(item.Type)},
						"id":   {S: aws.String(item.ID)},
					},
				},
			}
		} else {
			elem = &dynamodb.WriteRequest{
				PutRequest: &dynamodb.PutRequest{
					Item: MarshallDynamodbItem(item),
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

	afterTimeStamp := int64(math.MinInt64)
	if after != nil {
		afterTimeStamp = after.UnixNano() / int64(time.Millisecond)
	}

	params := &dynamodb.QueryInput{
		TableName:              aws.String(s.stateStoreTableName),
		IndexName:              aws.String(s.stateStoreIndexByTime),
		KeyConditionExpression: aws.String("#pk = :pk AND #sk >= :after"),
		ExpressionAttributeNames: map[string]*string{
			"#pk": aws.String("type"),
			"#sk": aws.String("updated_at"),
		},
		ExpressionAttributeValues: map[string]*dynamodb.AttributeValue{
			":pk":    MarshallDynamodbAttribute(typ),
			":after": MarshallDynamodbAttribute(afterTimeStamp),
		},
		ScanIndexForward: aws.Bool(false),
		Limit:            FixPaginationLimit(limit),
	}

	if nextToken != nil {
		token, err := base64.StdEncoding.DecodeString(*nextToken)
		if err != nil {
			panic(errors.Wrapf(err, "cannot decode next token"))
		}
		tokenItems := strings.Split(string(token), "#")
		timeStamp, err := strconv.ParseInt(tokenItems[1], 10, 64)
		if err != nil {
			panic(errors.Wrapf(err, "invalid next token value"))
		}
		params.ExclusiveStartKey = map[string]*dynamodb.AttributeValue{
			"id":         MarshallDynamodbAttribute(tokenItems[0]),
			"type":       MarshallDynamodbAttribute(typ),
			"updated_at": MarshallDynamodbAttribute(timeStamp),
		}
	}

	output, err := s.dynamodbClient.QueryWithContext(ctx, params)
	if err != nil {
		panic(errors.Wrapf(err, "cannot get state records"))
	}

	var newNextToken *string
	if output.LastEvaluatedKey != nil {
		id := *output.LastEvaluatedKey["id"].S
		timeStamp := *output.LastEvaluatedKey["updated_at"].N
		token := fmt.Sprintf("%s#%s", id, timeStamp)
		newNextToken = aws.String(base64.StdEncoding.EncodeToString([]byte(token)))
	}

	if len(output.Items) == 0 {
		return &model.AwsStateRecordPage{
			Items:     make([]*model.AwsStateRecord, 0),
			NextToken: newNextToken,
		}
	}

	records := make([]*model.AwsStateRecord, 0)
	UnmarshallDynamodbItemList(output.Items, &records)

	return &model.AwsStateRecordPage{
		Items:     records,
		NextToken: newNextToken,
	}
}
