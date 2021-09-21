package sqs

import (
	"fmt"
	"github.com/aws/aws-lambda-go/events"
	"github.com/cevixe/aws-sdk-go/aws/serdes/json"
	"github.com/cevixe/aws-sdk-go/aws/util"
	"github.com/pkg/errors"
	"reflect"
)

func UnmarshallSQSEvent(sqsEvent events.SQSEvent, record interface{}) {

	recordType := reflect.TypeOf(record)
	if recordType.Kind() != reflect.Ptr {
		panic(errors.New("record parameter must be a pointer"))
	}

	messages := make([]*map[string]interface{}, 0, len(sqsEvent.Records))
	for _, sqsMessage := range sqsEvent.Records {
		generic := make(map[string]interface{})
		snsMessage := events.SNSEntity{}
		util.UnmarshalJsonString(sqsMessage.Body, &snsMessage)
		fmt.Printf("message: %s\n", json.Marshall(snsMessage))
		util.UnmarshalJsonString(snsMessage.Message, &generic)
		messages = append(messages, &generic)
	}

	if len(messages) == 1 &&
		recordType.Elem().Kind() != reflect.Slice &&
		recordType.Elem().Kind() != reflect.Array {
		json := util.MarshalJson(messages[0])
		util.UnmarshalJson(json, record)
	} else {
		json := util.MarshalJson(messages)
		util.UnmarshalJson(json, record)
	}
}
