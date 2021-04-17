package appsync

import (
	"bytes"
	"context"
	v4 "github.com/aws/aws-sdk-go/aws/signer/v4"
	"github.com/cevixe/aws-sdk-go/aws/integration/session"
	"github.com/cevixe/aws-sdk-go/aws/model"
	"github.com/cevixe/aws-sdk-go/util"
	"net/http"
	"time"
)

type graphqlGatewayImpl struct {
	model.GraphqlGateway
	appsyncRegion string
	appsyncUrl    string
	signer        *v4.Signer
	client        *http.Client
}

func (g graphqlGatewayImpl) ExecuteGraphql(_ context.Context, request *model.GraphqlRequest) (*model.GraphqlResponse, error) {

	requestBuffer := util.MarshalJson(request)
	req, err := http.NewRequest("POST", g.appsyncUrl, bytes.NewReader(requestBuffer))
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", "application/json")

	_, err = g.signer.Sign(req, bytes.NewReader(requestBuffer), "appsync", g.appsyncRegion, time.Now())
	if err != nil {
		return nil, err
	}

	response, err := g.client.Do(req)
	if err != nil {
		return nil, err
	}

	buf := new(bytes.Buffer)
	_, err = buf.ReadFrom(response.Body)
	if err != nil {
		return nil, err
	}

	graphqlResponse := &model.GraphqlResponse{}
	util.UnmarshalJson(buf.Bytes(), graphqlResponse)

	return graphqlResponse, nil
}

func NewGraphqlGateway(
	appsyncRegion string,
	appsyncUrl string,
	factory session.Factory) model.GraphqlGateway {

	signer := v4.NewSigner(factory.GetCredentials(appsyncRegion))
	return &graphqlGatewayImpl{
		appsyncRegion: appsyncRegion,
		appsyncUrl:    appsyncUrl,
		signer:        signer,
		client:        factory.GetClient(),
	}
}
