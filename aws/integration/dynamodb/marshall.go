package dynamodb

import (
	"fmt"
	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute"
	"github.com/cevixe/aws-sdk-go/aws/util"
)

func MarshallDynamodbItem(record interface{}) map[string]*dynamodb.AttributeValue {
	output, err := dynamodbattribute.MarshalMap(record)
	if err != nil {
		panic(fmt.Errorf("cannot marshal dynamodb item\n%v", err))
	}
	return output
}

func UnmarshallDynamodbItem(item map[string]*dynamodb.AttributeValue, record interface{}) {
	err := dynamodbattribute.UnmarshalMap(item, record)
	if err != nil {
		panic(fmt.Errorf("cannot unmarshal dynamodb item\n%v", err))
	}
}

func UnmarshallDynamodbItemList(items []map[string]*dynamodb.AttributeValue, records interface{}) {
	err := dynamodbattribute.UnmarshalListOfMaps(items, records)
	if err != nil {
		panic(fmt.Errorf("cannot unmarshal dynamodb item list\n%v", err))
	}
}

func UnmarshallDynamodbStreamItem(item map[string]events.DynamoDBAttributeValue, record interface{}) {

	imageJsonBuffer := util.MarshalJsonString(item)

	var dynamoMap map[string]*dynamodb.AttributeValue
	util.UnmarshalJsonString(imageJsonBuffer, &dynamoMap)

	err := dynamodbattribute.UnmarshalMap(dynamoMap, record)
	if err != nil {
		panic(fmt.Errorf("cannot unmarshal dynamodb stream item\n%v", err))
	}
}
