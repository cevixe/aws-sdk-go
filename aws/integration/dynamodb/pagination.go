package dynamodb

import "github.com/aws/aws-sdk-go/aws"

const DefaultPaginationLimit = 20
const MaxPaginationLimit = 100
const MinPaginationLimit = 1

func FixPaginationLimit(limit *int64) *int64 {
	if limit == nil {
		return aws.Int64(DefaultPaginationLimit)
	}
	if *limit < MinPaginationLimit {
		return aws.Int64(MinPaginationLimit)
	}
	if *limit > MaxPaginationLimit {
		return aws.Int64(MaxPaginationLimit)
	}
	return limit
}
