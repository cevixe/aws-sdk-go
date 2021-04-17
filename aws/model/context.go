package model

const AwsContext = "AwsContext"

type Context struct {
	AwsObjectStore  ObjectStore
	AwsEventStore   EventStore
	AwsEventBus     EventBus
	AwsStateStore   StateStore
	AwsGraphGateway GraphqlGateway
}
