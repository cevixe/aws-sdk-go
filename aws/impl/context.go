package impl

import "github.com/cevixe/aws-sdk-go/aws/model"

const AwsContext = "AwsContext"

type Context struct {
	AwsHandlerID      string
	AwsHandlerVersion string
	AwsHandlerTimeout uint64
	AwsObjectStore    model.AwsObjectStore
	AwsEventStore     model.AwsEventStore
	AwsEventBus       model.AwsEventBus
	AwsStateStore     model.AwsStateStore
	AwsGraphqlGateway model.AwsGraphqlGateway
}
