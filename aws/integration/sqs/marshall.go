package sqs

import (
	"encoding/base64"
	"fmt"
	"github.com/aws/aws-lambda-go/events"
	"github.com/cevixe/aws-sdk-go/aws/serdes/gzip"
	"github.com/cevixe/aws-sdk-go/aws/util"
	"reflect"
)

func UnmarshallSQSEvent(sqsEvent events.SQSEvent, record interface{}) {

	recordType := reflect.TypeOf(record)
	if recordType.Kind() != reflect.Ptr {
		panic(fmt.Errorf("record parameter must be a pointer"))
	}

	messages := make([]*map[string]interface{}, 0, len(sqsEvent.Records))
	for _, sqsMessage := range sqsEvent.Records {
		generic := &map[string]interface{}{}
		snsMessage := &events.SNSEntity{}
		util.UnmarshalJsonString(sqsMessage.Body, snsMessage)

		buff, err := base64.StdEncoding.DecodeString(snsMessage.Message)
		if err != nil {
			util.UnmarshalJsonString(snsMessage.Message, generic)
		} else {
			jsonBuff := gzip.Decompress(buff)
			util.UnmarshalJsonString(string(jsonBuff), generic)
		}
		messages = append(messages, generic)
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
