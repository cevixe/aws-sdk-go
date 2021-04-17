package model

import "context"

type GraphqlRequest struct {
	Query     string                 `json:"query"`
	Variables map[string]interface{} `json:"variables"`
}

type GraphqlResponse struct {
	Data   map[string]interface{}   `json:"data"`
	Errors []map[string]interface{} `json:"errors"`
}

type GraphqlGateway interface {
	ExecuteGraphql(ctx context.Context, request *GraphqlRequest) (*GraphqlResponse, error)
}
