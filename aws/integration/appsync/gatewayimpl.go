package appsync

import (
	"bytes"
	"context"
	"github.com/aws/aws-sdk-go/aws/signer/v4"
	"github.com/cevixe/aws-sdk-go/aws/env"
	"github.com/cevixe/aws-sdk-go/aws/integration/session"
	"github.com/cevixe/aws-sdk-go/aws/model"
	util2 "github.com/cevixe/aws-sdk-go/aws/util"
	"github.com/pkg/errors"
	"net/http"
	"os"
	"time"
)

type gatewayImpl struct {
	region string
	url    string
	signer *v4.Signer
	client *http.Client
}

func NewAppsyncGateway(
	region string,
	url string,
	signer *v4.Signer,
	client *http.Client) model.AwsGraphqlGateway {

	return &gatewayImpl{
		region: region,
		url:    url,
		signer: signer,
		client: client,
	}
}

func NewDefaultAppsyncGateway(sessionFactory session.Factory) model.AwsGraphqlGateway {

	region := os.Getenv(env.AwsRegion)
	url := os.Getenv(env.CevixeGraphqlGatewayUrl)
	sess := sessionFactory.GetSession(region)
	client := sess.Config.HTTPClient
	signer := v4.NewSigner(sess.Config.Credentials)

	return NewAppsyncGateway(region, url, signer, client)
}

func (g gatewayImpl) buildHttpRequest(request *model.AwsGraphqlRequest) *http.Request {

	body := util2.MarshalJson(request)
	req, err := http.NewRequest("POST", g.url, bytes.NewReader(body))
	if err != nil {
		panic(errors.Wrap(err, "cannot generate http request"))
	}
	req.Header.Set("Content-Type", "application/json")

	_, err = g.signer.Sign(req, bytes.NewReader(body), "appsync", g.region, time.Now())
	if err != nil {
		panic(errors.Wrap(err, "cannot sign http request"))
	}

	return req
}

func (g gatewayImpl) readHttpResponse(response *http.Response) *model.AwsGraphqlResponse {

	buf := new(bytes.Buffer)
	_, err := buf.ReadFrom(response.Body)
	if err != nil {

		panic(errors.Wrap(err, "cannot read http response"))
	}

	res := &model.AwsGraphqlResponse{}
	util2.UnmarshalJson(buf.Bytes(), res)

	return res
}

func (g gatewayImpl) ExecuteGraphql(_ context.Context, request *model.AwsGraphqlRequest) *model.AwsGraphqlResponse {

	httpRequest := g.buildHttpRequest(request)

	response, err := g.client.Do(httpRequest)
	if err != nil {
		panic(errors.Wrap(err, "unexpected error in http call to appsync"))
	}

	return g.readHttpResponse(response)
}
