package impl

import (
	"github.com/cevixe/aws-sdk-go/aws/factory"
	"github.com/cevixe/aws-sdk-go/aws/model"
)

const AwsContext = "AwsContext"

type Context struct {
	AwsHandlerID      string
	AwsHandlerVersion string
	AwsHandlerTimeout uint64
	AwsFactory        factory.AwsFactory
	AwsObjectStore    model.AwsObjectStore
	AwsEventStore     model.AwsEventStore
	AwsEventBus       model.AwsEventBus
	AwsStateStore     model.AwsStateStore
	AwsCounterStore   model.AwsCounterStore
	AwsGraphqlGateway model.AwsGraphqlGateway
}
