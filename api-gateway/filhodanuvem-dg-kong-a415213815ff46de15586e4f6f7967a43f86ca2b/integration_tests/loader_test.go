package tests

import (
	"io"
	"os"
	"testing"
	"time"

	"github.com/devgymbr/kong"
	"github.com/devgymbr/kong/config"
)

func TestInitialLoader(t *testing.T) {
	pwd, err := os.Getwd()
	if err != nil {
		t.Fatalf("unable to get current directory: %s", err)
	}
	path := pwd + "/files/initial_loader.yaml"

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

	fp, err := os.Create(path)
	if err != nil {
		t.Fatalf("unable to create temporary file %s: %s", path, err)
	}
	defer fp.Close()

	_, err = fp.Write(yaml)
	if err != nil {
		t.Fatalf("unable to write on temporary file %s: %s", path, err)
	}

	tooLong := 10 * time.Second
	c := kong.Config{}
	err = config.Loader(&c, path, tooLong)
	if err != nil {
		t.Fatalf("unable to load config: %s", err)
	}

	if len(c.Services) != 1 || len(c.Services[0].Routes) != 1 {
		t.Fatalf("expected 1 service and 1 route (no refresh before 10sec), got %d services and %d routes", len(c.Services), len(c.Services[0].Routes))
	}

	if c.Services[0].Name != "payments" {
		t.Fatalf("expected service name to be payments, got %s", c.Services[0].Name)
	}

	if c.Services[0].URL != "http://localhost:8081" {
		t.Fatalf("expected service url to be http://localhost:8081, got %s", c.Services[0].URL)
	}

	if c.Services[0].Routes[0].Paths[0] != "/payments" {
		t.Fatalf("expected route path to be /payments, got %s", c.Services[0].Routes[0].Paths[0])
	}

	if c.Services[0].Routes[0].Methods[0] != "GET" {
		t.Fatalf("expected route method to be GET, got %s", c.Services[0].Routes[0].Methods[0])
	}
}

func TestLongRefreshTickerDoesNotRefresh(t *testing.T) {
	pwd, err := os.Getwd()
	if err != nil {
		t.Fatalf("unable to get current directory: %s", err)
	}
	path := pwd + "/files/one_route_get_payments.yaml"

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

	fp, err := os.Create(path)
	if err != nil {
		t.Fatalf("unable to create temporary file %s: %s", path, err)
	}
	defer fp.Close()

	_, err = fp.Write(yaml)
	if err != nil {
		t.Fatalf("unable to write on temporary file %s: %s", path, err)
	}

	tooLong := 10 * time.Second
	c := kong.Config{}
	err = config.Loader(&c, path, tooLong)
	if err != nil {
		t.Fatalf("unable to load config: %s", err)
	}

	// update yaml
	yaml = []byte(
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
		- POST # !!!!! added POST method 
`)

	if _, err = fp.Seek(int64(0), 0); err != nil {
		t.Fatalf("unable to truncate temporary file %s: %s", path, err)
	}

	_, err = fp.Write(yaml)
	if err != nil {
		t.Fatalf("unable to override on temporary file %s: %s", path, err)
	}

	// wait for the a few seconds but before the interval to refresh
	time.Sleep(2 * time.Second)

	if len(c.Services) != 1 || len(c.Services[0].Routes) != 1 || len(c.Services[0].Routes[0].Methods) != 1 {
		t.Fatalf("expected 1 service with a refreshed route with %d methods", len(c.Services[0].Routes[0].Methods))
	}
}

func TestShortRefreshTickerDoesRefresh(t *testing.T) {
	pwd, err := os.Getwd()
	if err != nil {
		t.Fatalf("unable to get current directory: %s", err)
	}
	path := pwd + "/files/two_route_get_payments.yaml"

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

	fp, err := os.Create(path)
	if err != nil {
		t.Fatalf("unable to create temporary file %s: %s", path, err)
	}
	defer fp.Close()

	_, err = fp.Write(yaml)
	if err != nil {
		t.Fatalf("unable to write on temporary file %s: %s", path, err)
	}

	tooShort := 5 * time.Millisecond
	c := kong.Config{}
	err = config.Loader(&c, path, tooShort)
	if err != nil {
		t.Fatalf("unable to load config: %s", err)
	}

	// update yaml
	yaml = []byte(
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
    - POST # !!!!! added POST method 
`)

	if _, err = fp.Seek(int64(0), io.SeekStart); err != nil {
		t.Fatalf("unable to truncate temporary file %s: %s", path, err)
	}

	_, err = fp.Write(yaml)
	if err != nil {
		t.Fatalf("unable to override on temporary file %s: %s", path, err)
	}
	fp.Sync()

	// wait for the a few seconds but before the interval to refresh
	time.Sleep(30 * time.Millisecond)

	if len(c.Services) != 1 || len(c.Services[0].Routes) != 1 || len(c.Services[0].Routes[0].Methods) != 2 {
		t.Fatalf("expected 1 service with a refreshed route with %d methods", len(c.Services[0].Routes[0].Methods))
	}
}
