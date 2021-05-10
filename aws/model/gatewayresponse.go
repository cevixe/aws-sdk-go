package model

type AwsGraphqlResponse struct {
	Data   map[string]interface{}   `json:"data"`
	Errors []map[string]interface{} `json:"errors"`
}
