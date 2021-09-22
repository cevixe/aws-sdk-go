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

func (s counterStoreImpl) updateValue(ctx context.Context, category string, name string) (uint64, error) {

	params := &dynamodb.UpdateItemInput{
		TableName: aws.String(s.counterStoreTableName),
		Key: map[string]*dynamodb.AttributeValue{
			"category": MarshallDynamodbAttribute(category),
			"name":     MarshallDynamodbAttribute(name),
		},
		UpdateExpression: aws.String("SET #v = #v + :i"),
		ExpressionAttributeNames: map[string]*string{
			"#v": aws.String("value"),
		},
		ExpressionAttributeValues: map[string]*dynamodb.AttributeValue{
			":i": MarshallDynamodbAttribute(1),
		},
		ReturnValues: aws.String("UPDATED_NEW"),
	}

	output, err := s.dynamodbClient.UpdateItemWithContext(ctx, params)
	if err != nil {
		return 0, err
	}

	newCounterString := *output.Attributes["value"].N
	newCounter, err := strconv.ParseUint(newCounterString, 10, 64)
	if err != nil {
		panic(errors.Wrap(err, "cannot parse counter value"))
	}
	return newCounter, nil
}

func (s counterStoreImpl) initiateValue(ctx context.Context, category string, name string) (uint64, error) {

	params := &dynamodb.PutItemInput{
		TableName: aws.String(s.counterStoreTableName),
		Item: map[string]*dynamodb.AttributeValue{
			"category": MarshallDynamodbAttribute(category),
			"name":     MarshallDynamodbAttribute(name),
			"value":    MarshallDynamodbAttribute(1),
		},
		ConditionExpression: aws.String("attribute_not_exists(#pk) AND attribute_not_exists(#sk)"),
		ExpressionAttributeNames: map[string]*string{
			"#pk": aws.String("category"),
			"#sk": aws.String("name"),
		},
	}

	_, err := s.dynamodbClient.PutItemWithContext(ctx, params)
	if err != nil {
		return 0, err
	}

	return 1, nil
}

func (s counterStoreImpl) NewValue(ctx context.Context, category string, name string) uint64 {

	value, err := s.updateValue(ctx, category, name)
	if err != nil {
		if ae, ok := err.(awserr.RequestFailure); ok && ae.Code() == "ValidationException" {
			value, err = s.initiateValue(ctx, category, name)
			if err != nil {
				if ae, ok := err.(awserr.RequestFailure); ok && ae.Code() == "ConditionalCheckFailedException" {
					value, err = s.updateValue(ctx, category, name)
					if err != nil {
						panic(errors.Wrap(err, "cannot update counter value"))
					}
				} else {
					panic(errors.Wrap(err, "cannot initiate counter value"))
				}
			}
		} else {
			panic(errors.Wrap(err, "cannot update counter value"))
		}
	}
	return value
}

func (s counterStoreImpl) GetValue(ctx context.Context, category string, name string) uint64 {

	params := &dynamodb.GetItemInput{
		TableName: aws.String(s.counterStoreTableName),
		Key: map[string]*dynamodb.AttributeValue{
			"category": MarshallDynamodbAttribute(category),
			"name":     MarshallDynamodbAttribute(name),
		},
	}

	output, err := s.dynamodbClient.GetItemWithContext(ctx, params)
	if err != nil {
		panic(errors.Wrap(err, "cannot get counter value"))
	}

	if output.Item == nil {
		return 0
	}

	counterString := *output.Item["counter"].N
	counter, err := strconv.ParseUint(counterString, 10, 64)
	if err != nil {
		panic(errors.Wrap(err, "cannot parse counter value"))
	}
	return counter
}
