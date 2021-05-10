package model

import "context"

type AwsGraphqlGateway interface {
	ExecuteGraphql(ctx context.Context, request *AwsGraphqlRequest) *AwsGraphqlResponse
}
