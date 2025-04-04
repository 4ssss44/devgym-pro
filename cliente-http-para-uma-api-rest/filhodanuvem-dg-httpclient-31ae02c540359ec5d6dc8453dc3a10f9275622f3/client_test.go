package kube

import (
	"testing"
	"time"
)

func TestDefaultClient(t *testing.T) {
	c, err := NewClient()
	if err != nil {
		t.Errorf("should not fail to create client: %s", err)
	}

	if c.httpClient == nil {
		t.Error("expected httpClient to be set")
	}
}

func TestWithURL(t *testing.T) {
	c, err := NewClient(WithURL("http://localhost:8080"))
	if err != nil {
		t.Errorf("should not fail to create client: %s", err)
	}
	if c.url != "http://localhost:8080" {
		t.Errorf("expected url to be http://localhost:8080, got %s", c.url)
	}
}

func TestWithInvalidURL(t *testing.T) {
	_, err := NewClient(WithURL("noturl"))
	if err == nil {
		t.Error("should fail to create client")
	}
}

func TestWithHTTPClient(t *testing.T) {
	c, err := NewClient(WithHTTPClient(nil))
	if err != nil {
		t.Errorf("should not fail to create client: %s", err)
	}
	if c.httpClient != nil {
		t.Error("expected httpClient to be nil")
	}
}

func TestWithTimeout(t *testing.T) {
	c, err := NewClient(WithTimeout(30 * time.Second))
	if err != nil {
		t.Errorf("should not fail to create client: %s", err)
	}
	if c.timeout != 30*time.Second {
		t.Errorf("expected timeout to be 30s, got %s", c.timeout)
	}
}
