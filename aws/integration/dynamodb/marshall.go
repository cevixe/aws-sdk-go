package dynamodb

import (
	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute"
	"github.com/cevixe/aws-sdk-go/aws/util"
	"github.com/pkg/errors"
)

func MarshallDynamodbAttribute(attribute interface{}) *dynamodb.AttributeValue {
	output, err := dynamodbattribute.Marshal(attribute)
	if err != nil {
		panic(errors.Wrap(err, "cannot marshal dynamodb attribute"))
	}
	return output
}

func MarshallDynamodbItem(record interface{}) map[string]*dynamodb.AttributeValue {
	output, err := dynamodbattribute.MarshalMap(record)
	if err != nil {
		panic(errors.Wrap(err, "cannot marshal dynamodb item"))
	}
	return output
}

func UnmarshallDynamodbItem(item map[string]*dynamodb.AttributeValue, record interface{}) {
	err := dynamodbattribute.UnmarshalMap(item, record)
	if err != nil {
		panic(errors.Wrap(err, "cannot unmarshal dynamodb item"))
	}
}

func UnmarshallDynamodbItemList(items []map[string]*dynamodb.AttributeValue, records interface{}) {
	err := dynamodbattribute.UnmarshalListOfMaps(items, records)
	if err != nil {
		panic(errors.Wrap(err, "cannot unmarshal dynamodb item list"))
	}
}

func UnmarshallDynamodbStreamItem(item map[string]events.DynamoDBAttributeValue, record interface{}) {

	imageJsonBuffer := util.MarshalJsonString(item)

	var dynamoMap map[string]*dynamodb.AttributeValue
	util.UnmarshalJsonString(imageJsonBuffer, &dynamoMap)

	err := dynamodbattribute.UnmarshalMap(dynamoMap, record)
	if err != nil {
		panic(errors.Wrap(err, "cannot unmarshal dynamodb stream item"))
	}
}
