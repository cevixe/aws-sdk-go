package services

type AWSService string

const (
	DynamoDB AWSService = "dynamodb"
	S3                  = "s3"
	SNS                 = "sns"
	SQS                 = "sqs"
)
