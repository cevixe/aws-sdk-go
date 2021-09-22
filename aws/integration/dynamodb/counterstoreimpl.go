package dynamodb

import (
	"context"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/awserr"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbiface"
	"github.com/cevixe/aws-sdk-go/aws/env"
	"github.com/cevixe/aws-sdk-go/aws/factory"
	"github.com/cevixe/aws-sdk-go/aws/model"
	"github.com/pkg/errors"
	"os"
	"strconv"
)

type counterStoreImpl struct {
	counterStoreTableName string
	dynamodbClient        dynamodbiface.DynamoDBAPI
}

func NewDynamodbCounterStore(
	counterStoreTableName string,
	dynamodbClient dynamodbiface.DynamoDBAPI) model.AwsCounterStore {

	return &counterStoreImpl{
		counterStoreTableName: counterStoreTableName,
		dynamodbClient:        dynamodbClient,
	}
}

func NewDefaultDynamodbCounterStore(awsFactory factory.AwsFactory) model.AwsCounterStore {

	counterStoreTableName := os.Getenv(env.CevixeCounterStoreTableName)
	dynamodbClient := awsFactory.DynamodbClient()

	return NewDynamodbCounterStore(counterStoreTableName, dynamodbClient)
}

func (s counterStoreImpl) NewValue(ctx context.Context, category string) uint64 {

	params := &dynamodb.UpdateItemInput{
		TableName: aws.String(s.counterStoreTableName),
		Key: map[string]*dynamodb.AttributeValue{
			"category": MarshallDynamodbAttribute(category),
		},
		UpdateExpression: aws.String("SET #counter = #counter + :incr"),
		ExpressionAttributeNames: map[string]*string{
			"#counter": aws.String("counter"),
		},
		ExpressionAttributeValues: map[string]*dynamodb.AttributeValue{
			":incr": MarshallDynamodbAttribute(1),
		},
		ReturnValues: aws.String("UPDATED_NEW"),
	}

	output, err := s.dynamodbClient.UpdateItemWithContext(ctx, params)
	if err != nil {
		panic(errors.Wrap(err, "cannot update category counter"))
	}

	if len(output.Attributes) == 0 {
		putParams := &dynamodb.PutItemInput{
			TableName: aws.String(s.counterStoreTableName),
			Item: map[string]*dynamodb.AttributeValue{
				"category": MarshallDynamodbAttribute(category),
				"counter":  MarshallDynamodbAttribute(1),
			},
			ConditionExpression: aws.String("attribute_not_exists(#pk) AND attribute_not_exists(#sk)"),
			ExpressionAttributeNames: map[string]*string{
				"#pk": aws.String("category"),
				"#sk": aws.String("counter"),
			},
		}
		if _, err := s.dynamodbClient.PutItemWithContext(ctx, putParams); err != nil {
			if ae, ok := err.(awserr.RequestFailure); ok && ae.Code() == "ConditionalCheckFailedException" {
				output, err = s.dynamodbClient.UpdateItemWithContext(ctx, params)
				if err != nil {
					panic(errors.Wrap(err, "cannot update category counter"))
				}
			} else {
				panic(errors.Wrap(err, "cannot initiate category counter"))
			}
		}
		return 1
	}

	newCounterString := *output.Attributes["counter"].N
	newCounter, err := strconv.ParseUint(newCounterString, 10, 64)
	if err != nil {
		panic(errors.Wrap(err, "cannot parse category counter"))
	}
	return newCounter
}

func (s counterStoreImpl) GetValue(ctx context.Context, category string) uint64 {

	params := &dynamodb.QueryInput{
		TableName:              aws.String(s.counterStoreTableName),
		KeyConditionExpression: aws.String("#category = :category"),
		ExpressionAttributeNames: map[string]*string{
			"#category": aws.String("category"),
		},
		ExpressionAttributeValues: map[string]*dynamodb.AttributeValue{
			":category": MarshallDynamodbAttribute(category),
		},
	}

	output, err := s.dynamodbClient.QueryWithContext(ctx, params)
	if err != nil {
		panic(errors.Wrap(err, "cannot get category counter"))
	}

	if len(output.Items) == 0 {
		return 0
	}

	counterString := *output.Items[0]["counter"].N
	counter, err := strconv.ParseUint(counterString, 10, 64)
	if err != nil {
		panic(errors.Wrap(err, "cannot parse category counter"))
	}
	return counter
}
