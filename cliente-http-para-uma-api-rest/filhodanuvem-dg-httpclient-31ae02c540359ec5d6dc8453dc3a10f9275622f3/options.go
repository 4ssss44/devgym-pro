package kube

import (
	"net/http"
	"net/url"
	"time"
)

type option func(*Client) error

func WithURL(u string) option {
	return func(c *Client) error {
		if _, err := url.ParseRequestURI(u); err != nil {
			return err
		}
		c.url = u

		return nil
	}
}

func WithHTTPClient(httpClient *http.Client) option {
	return func(c *Client) error {
		c.httpClient = httpClient

		return nil
	}
}

func WithTimeout(timeout time.Duration) option {
	return func(c *Client) error {
		c.timeout = timeout

		return nil
	}
}
