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

type eventStoreImpl struct {
	eventStoreTable              string
	eventStoreIndexByTime        string
	eventStoreIndexByDay         string
	eventStoreIndexByType        string
	eventStoreIndexByAuthor      string
	eventStoreIndexByTransaction string
	controlStoreTable            string
	dynamodbClient               dynamodbiface.DynamoDBAPI
}

func NewDynamodbEventStore(
	eventStoreTable string,
	eventStoreIndexByTime string,
	eventStoreIndexByDay string,
	eventStoreIndexByType string,
	eventStoreIndexByAuthor string,
	eventStoreIndexByTransaction string,
	controlStoreTable string,
	dynamodbClient dynamodbiface.DynamoDBAPI) model.AwsEventStore {

	return &eventStoreImpl{
		eventStoreTable:              eventStoreTable,
		eventStoreIndexByTime:        eventStoreIndexByTime,
		eventStoreIndexByDay:         eventStoreIndexByDay,
		eventStoreIndexByType:        eventStoreIndexByType,
		eventStoreIndexByAuthor:      eventStoreIndexByAuthor,
		eventStoreIndexByTransaction: eventStoreIndexByTransaction,
		controlStoreTable:            controlStoreTable,
		dynamodbClient:               dynamodbClient,
	}
}

func NewDefaultDynamodbEventStore(awsFactory factory.AwsFactory) model.AwsEventStore {

	eventStoreTableName := os.Getenv(env.CevixeEventStoreTableName)
	eventStoreIndexByTime := os.Getenv(env.CevixeEventStoreIndexByTime)
	eventStoreIndexByDay := os.Getenv(env.CevixeEventStoreIndexByDay)
	eventStoreIndexByType := os.Getenv(env.CevixeEventStoreIndexByType)
	eventStoreIndexByAuthor := os.Getenv(env.CevixeEventStoreIndexByAuthor)
	eventStoreIndexByTransaction := os.Getenv(env.CevixeEventStoreIndexByTransaction)
	controlStoreTableName := os.Getenv(env.CevixeControlStoreTableName)
	dynamodbClient := awsFactory.DynamodbClient()

	return NewDynamodbEventStore(
		eventStoreTableName,
		eventStoreIndexByTime,
		eventStoreIndexByDay,
		eventStoreIndexByType,
		eventStoreIndexByAuthor,
		eventStoreIndexByTransaction,
		controlStoreTableName,
		dynamodbClient)
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
		panic(errors.Wrap(err, "cannot create control record"))
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
		panic(errors.Wrap(err, "cannot create uncontrolled event record"))
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
		panic(errors.Wrap(err, "cannot create controlled event record"))
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
		panic(errors.Wrap(err, "cannot get control records"))
	}

	records := &[]*model.AwsControlRecord{}
	UnmarshallDynamodbItemList(output.Items, records)
	return *records
}

func (e eventStoreImpl) GetEventRecordByID(ctx context.Context, source string, id string) *model.AwsEventRecord {
	input := &dynamodb.GetItemInput{
		TableName: aws.String(e.eventStoreTable),
		Key: map[string]*dynamodb.AttributeValue{
			"event_source": {S: aws.String(source)},
			"event_id":     {S: aws.String(id)},
		},
	}

	output, err := e.dynamodbClient.GetItemWithContext(ctx, input)
	if err != nil {
		panic(errors.Wrap(err, "cannot get event record by id"))
	}

	if output.Item == nil {
		return nil
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
		panic(errors.Wrapf(err, "cannot get last event record"))
	}

	if len(output.Items) == 0 {
		return nil
	}

	record := &model.AwsEventRecord{}
	UnmarshallDynamodbItem(output.Items[0], record)
	return record
}

func (e eventStoreImpl) GetEventPage(ctx context.Context, index string, pkName string, skName string, pkValue interface{},
	after *time.Time, before *time.Time, nextToken *string, limit *int64) *model.AwsEventRecordPage {

	afterTimeStamp := int64(math.MinInt64)
	if after != nil {
		afterTimeStamp = after.UnixNano() / int64(time.Millisecond)
	}

	beforeTimeStamp := int64(math.MaxInt64)
	if before != nil {
		beforeTimeStamp = before.UnixNano() / int64(time.Millisecond)
	}

	params := &dynamodb.QueryInput{
		TableName:              aws.String(e.eventStoreTable),
		IndexName:              aws.String(index),
		KeyConditionExpression: aws.String("#pk = :pk AND #sk BETWEEN :after AND :before"),
		ExpressionAttributeNames: map[string]*string{
			"#pk": aws.String(pkName),
			"#sk": aws.String(skName),
		},
		ExpressionAttributeValues: map[string]*dynamodb.AttributeValue{
			":pk":     MarshallDynamodbAttribute(pkValue),
			":after":  MarshallDynamodbAttribute(afterTimeStamp),
			":before": MarshallDynamodbAttribute(beforeTimeStamp),
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
			"event_id": MarshallDynamodbAttribute(tokenItems[0]),
			pkName:     MarshallDynamodbAttribute(pkValue),
			skName:     MarshallDynamodbAttribute(timeStamp),
		}
	}

	output, err := e.dynamodbClient.QueryWithContext(ctx, params)
	if err != nil {
		panic(errors.Wrapf(err, "cannot get event page"))
	}

	var newNextToken *string
	if output.LastEvaluatedKey != nil {
		id := *output.LastEvaluatedKey["event_id"].S
		timeStamp := *output.LastEvaluatedKey[skName].N
		token := fmt.Sprintf("%s#%s", id, timeStamp)
		newNextToken = aws.String(base64.StdEncoding.EncodeToString([]byte(token)))
	}

	if len(output.Items) == 0 {
		return &model.AwsEventRecordPage{
			Items:     make([]*model.AwsEventRecord, 0),
			NextToken: newNextToken,
		}
	}

	records := make([]*model.AwsEventRecord, 0)
	UnmarshallDynamodbItemList(output.Items, &records)

	return &model.AwsEventRecordPage{
		Items:     records,
		NextToken: newNextToken,
	}
}

func (e eventStoreImpl) GetEventHeaders(ctx context.Context, source string,
	after *string, before *string, nextToken *string, limit *int64) *model.AwsEventHeaderRecordPage {

	afterToken := fmt.Sprintf("%020d", 0)
	if after != nil {
		afterToken = *after
	}

	beforeToken := "99999999999999999999"
	if before != nil {
		beforeToken = *before
	}

	params := &dynamodb.QueryInput{
		TableName:              aws.String(e.eventStoreTable),
		ProjectionExpression:   aws.String("#pk,#sk,#class,#type,#at,#by,#day,#eid,#etype,#edeleted,#trx,#tsrc,#tid"),
		KeyConditionExpression: aws.String("#pk = :pk AND #sk BETWEEN :after AND :before"),
		ExpressionAttributeNames: map[string]*string{
			"#pk":       aws.String("event_source"),
			"#sk":       aws.String("event_id"),
			"#class":    aws.String("event_class"),
			"#type":     aws.String("event_type"),
			"#at":       aws.String("event_time"),
			"#by":       aws.String("event_author"),
			"#day":      aws.String("event_day"),
			"#eid":      aws.String("entity_id"),
			"#etype":    aws.String("entity_type"),
			"#edeleted": aws.String("entity_deleted"),
			"#trx":      aws.String("transaction"),
			"#tsrc":     aws.String("trigger_source"),
			"#tid":      aws.String("trigger_id"),
		},
		ExpressionAttributeValues: map[string]*dynamodb.AttributeValue{
			":pk":     MarshallDynamodbAttribute(source),
			":after":  MarshallDynamodbAttribute(afterToken),
			":before": MarshallDynamodbAttribute(beforeToken),
		},
		ScanIndexForward: aws.Bool(false),
		Limit:            FixPaginationLimit(limit),
	}

	if nextToken != nil {
		eventId, err := base64.StdEncoding.DecodeString(*nextToken)
		if err != nil {
			panic(errors.Wrapf(err, "cannot decode next token"))
		}
		params.ExclusiveStartKey = map[string]*dynamodb.AttributeValue{
			"event_source": MarshallDynamodbAttribute(source),
			"event_id":     MarshallDynamodbAttribute(string(eventId)),
		}
	}

	output, err := e.dynamodbClient.QueryWithContext(ctx, params)
	if err != nil {
		panic(errors.Wrapf(err, "cannot get event header page"))
	}

	var newNextToken *string
	if output.LastEvaluatedKey != nil {
		eventId := *output.LastEvaluatedKey["event_id"].S
		newNextToken = aws.String(base64.StdEncoding.EncodeToString([]byte(eventId)))
	}

	if len(output.Items) == 0 {
		return &model.AwsEventHeaderRecordPage{
			Items:     make([]*model.AwsEventHeaderRecord, 0),
			NextToken: newNextToken,
		}
	}

	records := make([]*model.AwsEventHeaderRecord, 0)
	UnmarshallDynamodbItemList(output.Items, &records)

	return &model.AwsEventHeaderRecordPage{
		Items:     records,
		NextToken: newNextToken,
	}
}

func (e eventStoreImpl) GetSourceEvents(ctx context.Context, source string,
	after *time.Time, before *time.Time, nextToken *string, limit *int64) *model.AwsEventRecordPage {

	return e.GetEventPage(ctx, e.eventStoreIndexByTime, "event_source", "event_time",
		source, after, before, nextToken, limit)
}

func (e eventStoreImpl) GetDayEvents(ctx context.Context, day string,
	after *time.Time, before *time.Time, nextToken *string, limit *int64) *model.AwsEventRecordPage {

	return e.GetEventPage(ctx, e.eventStoreIndexByDay, "event_day", "event_time",
		day, after, before, nextToken, limit)
}

func (e eventStoreImpl) GetTypeEvents(ctx context.Context, typ string,
	after *time.Time, before *time.Time, nextToken *string, limit *int64) *model.AwsEventRecordPage {

	return e.GetEventPage(ctx, e.eventStoreIndexByType, "event_type", "event_time",
		typ, after, before, nextToken, limit)
}

func (e eventStoreImpl) GetAuthorEvents(ctx context.Context, author string,
	after *time.Time, before *time.Time, nextToken *string, limit *int64) *model.AwsEventRecordPage {

	return e.GetEventPage(ctx, e.eventStoreIndexByAuthor, "event_author", "event_time",
		author, after, before, nextToken, limit)
}

func (e eventStoreImpl) GetTransactionEvents(ctx context.Context, transaction string,
	after *time.Time, before *time.Time, nextToken *string, limit *int64) *model.AwsEventRecordPage {

	return e.GetEventPage(ctx, e.eventStoreIndexByTransaction, "transaction", "event_time",
		transaction, after, before, nextToken, limit)
}
