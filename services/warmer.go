package services

import (
	"fmt"
	"net/http"
	"os"
	"sync"
)

type Warmer interface {
	WarmUp(services []AWSService)
}

type warmerImpl struct {
	Warmer
	region     string
	httpClient *http.Client
}

func NewWarmer(region string, httpClient *http.Client) Warmer {
	if region == "" {
		return NewWarmerWithDefaultRegion(httpClient)
	}
	return warmerImpl{
		region:     region,
		httpClient: httpClient,
	}
}

func NewWarmerWithDefaultRegion(httpClient *http.Client) Warmer {
	region := os.Getenv("AWS_REGION")
	return warmerImpl{
		region:     region,
		httpClient: httpClient,
	}
}
func (w warmerImpl) WarmUp(services []AWSService) {
	wg := &sync.WaitGroup{}
	wg.Add(len(services))
	for _, item := range services {
		go w.warmUpSingleService(item, wg)
	}
	wg.Wait()
}

const URLTemplate = "https://%s.%s.amazonaws.com"

func (w warmerImpl) warmUpSingleService(service AWSService, waitGroup *sync.WaitGroup) {
	url := fmt.Sprintf(URLTemplate, string(service), w.region)
	_, _ = w.httpClient.Head(url)
	waitGroup.Done()
}
