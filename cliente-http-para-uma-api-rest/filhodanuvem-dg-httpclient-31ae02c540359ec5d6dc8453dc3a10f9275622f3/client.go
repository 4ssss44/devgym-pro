package kube

import (
	"net/http"
	"time"

	"github.com/devgymbr/kubeclient/deployment"
)

const productionURL = "https://api.k8s.io"

type Client struct {
	url        string
	timeout    time.Duration
	httpClient *http.Client

	Deployment deployment.Service
}

func NewClient(options ...option) (*Client, error) {
	c := Client{
		url:        productionURL,
		httpClient: &http.Client{},
	}
	for _, option := range options {
		if err := option(&c); err != nil {
			return nil, err
		}
	}

	if c.timeout != 0 {
		c.httpClient.Timeout = c.timeout
	}

	c.Deployment = deployment.NewService(c.httpClient, c.url)

	return &c, nil
}
