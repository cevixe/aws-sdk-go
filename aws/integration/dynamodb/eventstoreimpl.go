package dynamodb

import (
	"context"
	"fmt"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbiface"
	"github.com/cevixe/aws-sdk-go/aws/env"
	"github.com/cevixe/aws-sdk-go/aws/integration/session"
	"github.com/cevixe/aws-sdk-go/aws/model"
	"os"
)

type eventStoreImpl struct {
	eventStoreTable   string
	controlStoreTable string
	dynamodbClient    dynamodbiface.DynamoDBAPI
}

func NewDynamodbEventStore(
	eventStoreTable string,
	controlStoreTable string,
	dynamodbClient dynamodbiface.DynamoDBAPI) model.AwsEventStore {

	return &eventStoreImpl{
		eventStoreTable:   eventStoreTable,
		controlStoreTable: controlStoreTable,
		dynamodbClient:    dynamodbClient,
	}
}

func NewDefaultDynamodbEventStore(sessionFactory session.Factory) model.AwsEventStore {

	region := os.Getenv(env.AwsRegion)
	eventStoreTableName := os.Getenv(env.CevixeEventStoreTableName)
	controlStoreTableName := os.Getenv(env.CevixeControlStoreTableName)
	dynamodbClient := dynamodb.New(sessionFactory.GetSession(region))

	return NewDynamodbEventStore(eventStoreTableName, controlStoreTableName, dynamodbClient)
}

func (e eventStoreImpl) CreateControlRecord(ctx context.Context, control *model.AwsControlRecord) {
	input := &dynamodb.PutItemInput{
		TableName:           aws.String(e.controlStoreTable),
		Item:                MarshallDynamodbItem(control),
		ConditionExpression: aws.String("attribute_not_exists(#pk) AND attribute_not_exists(#sk)"),
		ExpressionAttributeNames: map[string]*string{
			"#pk": aws.String("control_group"),
			"#sk": aws.String("control_id"),
		},
	}

	_, err := e.dynamodbClient.PutItemWithContext(ctx, input)
	if err != nil {
		panic(fmt.Errorf("cannot create control record\n%v", err))
	}
}

func (e eventStoreImpl) CreateUncontrolledEventRecord(ctx context.Context, event *model.AwsEventRecord) {
	input := &dynamodb.PutItemInput{
		TableName:           aws.String(e.eventStoreTable),
		Item:                MarshallDynamodbItem(event),
		ConditionExpression: aws.String("attribute_not_exists(#pk) AND attribute_not_exists(#sk)"),
		ExpressionAttributeNames: map[string]*string{
			"#pk": aws.String("event_source"),
			"#sk": aws.String("event_id"),
		},
	}

	_, err := e.dynamodbClient.PutItemWithContext(ctx, input)
	if err != nil {
		panic(fmt.Errorf("cannot create uncontrolled event record\n%v", err))
	}
}

func (e eventStoreImpl) CreateControlledEventRecord(ctx context.Context, event *model.AwsEventRecord, control *model.AwsControlRecord) {
	input := &dynamodb.TransactWriteItemsInput{
		TransactItems: []*dynamodb.TransactWriteItem{
			{
				Put: &dynamodb.Put{
					TableName:           aws.String(e.eventStoreTable),
					Item:                MarshallDynamodbItem(event),
					ConditionExpression: aws.String("attribute_not_exists(#pk) AND attribute_not_exists(#sk)"),
					ExpressionAttributeNames: map[string]*string{
						"#pk": aws.String("event_source"),
						"#sk": aws.String("event_id"),
					},
				},
			},
			{
				Put: &dynamodb.Put{
					TableName:           aws.String(e.controlStoreTable),
					Item:                MarshallDynamodbItem(control),
					ConditionExpression: aws.String("attribute_not_exists(#pk) AND attribute_not_exists(#sk)"),
					ExpressionAttributeNames: map[string]*string{
						"#pk": aws.String("control_group"),
						"#sk": aws.String("control_id"),
					},
				},
			},
		},
	}

	_, err := e.dynamodbClient.TransactWriteItemsWithContext(ctx, input)
	if err != nil {
		panic(fmt.Errorf("cannot create controlled event record\n%v", err))
	}
}

func (e eventStoreImpl) GetControlRecords(ctx context.Context, group string) []*model.AwsControlRecord {
	input := &dynamodb.QueryInput{
		TableName:              aws.String(e.controlStoreTable),
		KeyConditionExpression: aws.String("pk = :pk"),
		ScanIndexForward:       aws.Bool(false),
		ExpressionAttributeValues: map[string]*dynamodb.AttributeValue{
			":pk": {S: aws.String(group)},
		},
		ExpressionAttributeNames: map[string]*string{
			"#pk": aws.String("control_group"),
		},
	}

	output, err := e.dynamodbClient.QueryWithContext(ctx, input)
	if err != nil {
		panic(fmt.Errorf("cannot get control records\n%v", err))
	}

	records := &[]*model.AwsControlRecord{}
	UnmarshallDynamodbItemList(output.Items, records)
	return *records
}

func (e eventStoreImpl) GetEventRecordByID(ctx context.Context, source string, id string) *model.AwsEventRecord {
	input := &dynamodb.GetItemInput{
		TableName: aws.String(e.eventStoreTable),
		Key: map[string]*dynamodb.AttributeValue{
			"#pk": {S: aws.String(source)},
			"#sk": {S: aws.String(id)},
		},
		ExpressionAttributeNames: map[string]*string{
			"#pk": aws.String("event_source"),
			"#sk": aws.String("event_id"),
		},
	}

	output, err := e.dynamodbClient.GetItemWithContext(ctx, input)
	if err != nil {
		panic(fmt.Errorf("cannot get event record by id\n%v", err))
	}

	record := &model.AwsEventRecord{}
	UnmarshallDynamodbItem(output.Item, record)
	return record
}

func (e eventStoreImpl) GetLastEventRecord(ctx context.Context, source string) *model.AwsEventRecord {
	input := &dynamodb.QueryInput{
		TableName:              aws.String(e.eventStoreTable),
		KeyConditionExpression: aws.String("#pk = :pk"),
		ScanIndexForward:       aws.Bool(false),
		Limit:                  aws.Int64(1),
		ExpressionAttributeNames: map[string]*string{
			"#pk": aws.String("event_source"),
		},
		ExpressionAttributeValues: map[string]*dynamodb.AttributeValue{
			":pk": {S: aws.String(source)},
		},
	}

	output, err := e.dynamodbClient.QueryWithContext(ctx, input)
	if err != nil {
		panic(fmt.Errorf("cannot get last event record\n%v", err))
	}

	if len(output.Items) == 0 {
		return nil
	}

	record := &model.AwsEventRecord{}
	UnmarshallDynamodbItem(output.Items[0], record)
	return record
}
