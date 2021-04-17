package dynamo

import (
	"context"
	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute"
	"github.com/cevixe/aws-sdk-go/aws/impl"
	"github.com/cevixe/aws-sdk-go/aws/model"
	"github.com/cevixe/aws-sdk-go/util"
	"github.com/cevixe/core-sdk-go/core"
)

func MapDynamoEventRecordToCevixeEvent(ctx context.Context, record events.DynamoDBEventRecord) core.Event {

	rawRecord := record.Change.NewImage
	rawJsonBuffer := util.MarshalJsonString(rawRecord)

	var dynamoMap map[string]*dynamodb.AttributeValue
	util.UnmarshalJsonString(rawJsonBuffer, &dynamoMap)

	eventValue := &model.EventObject{}
	err := dynamodbattribute.UnmarshalMap(dynamoMap, eventValue)
	if err != nil {
		panic(err)
	}

	return impl.NewEvent(ctx, eventValue)
}
