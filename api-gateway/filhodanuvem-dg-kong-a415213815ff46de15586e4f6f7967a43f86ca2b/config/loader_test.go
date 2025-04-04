package config

import (
	"testing"
	"time"

	"github.com/devgymbr/kong"
)

func TestEmptyYaml(t *testing.T) {
	c := kong.Config{}
	yaml := []byte(``)

	err := c.Refresh(yaml, time.Now())

	if err != nil {
		t.Errorf("expected error to be nil got %v", err)
	}

	if len(c.Services) != 0 {
		t.Errorf("expected services to be empty got %v", c.Services)
	}
}

func TestSimpleYaml(t *testing.T) {
	c := kong.Config{}
	yaml := []byte(
		`
services:
- name: payments
  url: http://localhost:8081
  plugins:
  - name: http_log
  - name: add_header
  routes:
  - paths:
    - /payments
    methods:
    - GET	
`)

	err := c.Refresh(yaml, time.Now())

	if err != nil {
		t.Errorf("expected error to be nil got %v", err)
	}

	if len(c.Services) != 1 {
		t.Errorf("expected services to have 1 element got %v", len(c.Services))
	}

	if c.Services[0].Name != "payments" || c.Services[0].URL != "http://localhost:8081" ||
		len(c.Services[0].Plugins) != 2 || len(c.Services[0].Routes) != 1 {
		t.Errorf("expected service name to be payments got %v", c.Services[0].Name)
	}
}
