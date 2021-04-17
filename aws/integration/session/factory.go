package session

import (
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-xray-sdk-go/xray"
	http2 "github.com/cevixe/aws-sdk-go/http"
	"github.com/cevixe/aws-sdk-go/util"
	"net/http"
)

const AwsSessionFactory string = "AwsSessionFactory"

type Factory interface {
	GetClient() *http.Client
	GetCredentials(region string) *credentials.Credentials
	GetSession(region string) *session.Session
}

type factoryImpl struct {
	client   *http.Client
	sessions map[string]*session.Session
}

func (f factoryImpl) GetClient() *http.Client {
	return f.client
}

func (f factoryImpl) GetCredentials(region string) *credentials.Credentials {
	return f.GetSession(region).Config.Credentials
}

func (f factoryImpl) GetSession(region string) *session.Session {
	if f.sessions[region] != nil {
		return f.sessions[region]
	}

	newSession := NewSessionWithRegion(f.client, region)
	f.sessions[region] = newSession

	return newSession
}

func NewSessionWithRegion(client *http.Client, region string) *session.Session {
	sess := session.Must(
		session.NewSessionWithOptions(
			session.Options{
				Config: aws.Config{
					Region:                  aws.String(region),
					S3ForcePathStyle:        aws.Bool(true),
					DisableParamValidation:  aws.Bool(true),
					DisableComputeChecksums: aws.Bool(true),
					HTTPClient:              client,
				},
				SharedConfigState: session.SharedConfigEnable,
			}))
	xray.AWSSession(sess)
	warmer := util.NewWarmer(region, client)
	warmer.WarmUp([]string{"s3", "dynamodb", "sns"})
	return sess
}

func NewSessionFactory() Factory {
	factory := &factoryImpl{
		client:   http2.NewDefaultHTTPClient(),
		sessions: make(map[string]*session.Session),
	}
	return factory
}
