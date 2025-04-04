package tests

import (
	"io"
	nethttp "net/http"
	"net/http/httptest"
	"regexp"
	"strings"
	"testing"

	"github.com/devgymbr/kong"
	"github.com/devgymbr/kong/http"
)

func TestServerWithoutConfigRules(t *testing.T) {
	s := http.NewServer(&kong.Config{})

	req := httptest.NewRequest(nethttp.MethodGet, "/", nil)
	w := httptest.NewRecorder()

	s.ServeHTTP(w, req)

	res := w.Result()
	defer res.Body.Close()

	data, err := io.ReadAll(res.Body)
	if err != nil {
		t.Errorf("expected error to be nil got %v", err)
	}

	if len(data) != 0 {
		t.Errorf("expected data to be empty got %v", data)
	}

	if res.StatusCode != nethttp.StatusNotFound {
		t.Errorf("expected status code to be 404 got %v", res.StatusCode)
	}
}

func TestServerWithOneEndpointShouldMatch(t *testing.T) {
	expectedMethod := nethttp.MethodPost
	expectedPath := "/posts"

	api := httptest.NewServer(nethttp.HandlerFunc(func(w nethttp.ResponseWriter, r *nethttp.Request) {
		if r.Method != expectedMethod ||
			r.URL.Path != expectedPath {
			w.WriteHeader(nethttp.StatusNotFound)
			return
		}
		w.WriteHeader(nethttp.StatusCreated)
		w.Write([]byte("created payment"))
	}))
	defer api.Close()

	s := http.NewServer(&kong.Config{
		Services: []kong.Service{
			{
				Name: "test",
				URL:  api.URL,
				Routes: []kong.Route{
					{
						Name:       "test",
						Paths:      []string{expectedMethod},
						PathRegexp: regexp.MustCompile(`^/posts$`),
						Methods:    []string{expectedMethod},
					},
				},
			},
		},
	})

	req := httptest.NewRequest(expectedMethod, expectedPath, nil)
	w := httptest.NewRecorder()

	s.ServeHTTP(w, req)

	res := w.Result()
	defer res.Body.Close()

	data, err := io.ReadAll(res.Body)
	if err != nil {
		t.Errorf("expected error to be nil got %v", err)
	}

	if string(data) != "created payment" {
		t.Errorf("expected data to be empty got %v", data)
	}

	if res.StatusCode != nethttp.StatusCreated {
		t.Errorf("expected status code to be 201 got %v", res.StatusCode)
	}
}

func TestServerWithOneEndpointShouldNotMatch(t *testing.T) {
	expectedMethod := nethttp.MethodPost
	expectedPath := "/posts"

	api := httptest.NewServer(nethttp.HandlerFunc(func(w nethttp.ResponseWriter, r *nethttp.Request) {
		if r.Method != expectedMethod ||
			r.URL.Path != expectedPath {
			w.WriteHeader(nethttp.StatusNotFound)
			return
		}
		w.WriteHeader(nethttp.StatusCreated)
		w.Write([]byte("created payment"))
	}))
	defer api.Close()

	s := http.NewServer(&kong.Config{
		Services: []kong.Service{
			{
				Name: "test",
				URL:  api.URL,
				Routes: []kong.Route{
					{
						Name:       "test",
						Paths:      []string{expectedMethod},
						PathRegexp: regexp.MustCompile(`^/posts$`),
						Methods:    []string{expectedMethod},
					},
				},
			},
		},
	})

	requests := map[string]string{
		"GET":    "/posts",
		"PUT":    "/posts",
		"DELETE": "/posts",
		"POST":   "/posts/1",
	}

	for method, path := range requests {
		t.Run(method+" "+path, func(t *testing.T) {
			req := httptest.NewRequest(method, path, nil)
			w := httptest.NewRecorder()

			s.ServeHTTP(w, req)

			res := w.Result()
			defer res.Body.Close()

			data, err := io.ReadAll(res.Body)
			if err != nil {
				t.Errorf("expected error to be nil got %v", err)
			}

			if len(data) != 0 {
				t.Errorf("expected data to be empty got %v", data)
			}

			if res.StatusCode != nethttp.StatusNotFound {
				t.Errorf("expected status code to be 404 got %v", res.StatusCode)
			}
		})
	}

}

func TestServerWithTwoServicesForwardCorrectly(t *testing.T) {
	counterReqAPI1 := 0
	counterReqAPI2 := 0

	api1 := httptest.NewServer(nethttp.HandlerFunc(func(w nethttp.ResponseWriter, r *nethttp.Request) {
		counterReqAPI1++
		if r.Method != nethttp.MethodPost ||
			r.URL.Path != "/posts" {
			w.WriteHeader(nethttp.StatusNotFound)
			return
		}
		w.WriteHeader(nethttp.StatusCreated)
		w.Write([]byte("created payment"))
	}))
	defer api1.Close()

	api2 := httptest.NewServer(nethttp.HandlerFunc(func(w nethttp.ResponseWriter, r *nethttp.Request) {
		counterReqAPI2++
		if r.Method != nethttp.MethodGet ||
			strings.Index(r.URL.Path, "/posts") != 0 {
			w.WriteHeader(nethttp.StatusNotFound)
			return
		}
		w.WriteHeader(nethttp.StatusOK)
		w.Write([]byte("post"))
	}))
	defer api2.Close()

	s := http.NewServer(&kong.Config{
		Services: []kong.Service{
			{
				Name: "API 1",
				URL:  api1.URL,
				Routes: []kong.Route{
					{
						Name:       "Endpoint POST",
						Paths:      []string{"/posts"},
						PathRegexp: regexp.MustCompile(`^/posts$`),
						Methods:    []string{nethttp.MethodPost},
					},
				},
			},
			{
				Name: "API 2",
				URL:  api2.URL,
				Routes: []kong.Route{
					{
						Name:       "Endpoint Get",
						Paths:      []string{"/posts/{id}"},
						PathRegexp: regexp.MustCompile(`^/posts/[0-9a-zA-z]+$`),
						Methods:    []string{nethttp.MethodGet},
					},
				},
			},
		},
	})

	req := httptest.NewRequest(nethttp.MethodGet, "/posts/1", nil)
	w := httptest.NewRecorder()

	s.ServeHTTP(w, req)

	if counterReqAPI2 != 1 {
		t.Errorf("expected GET /posts/1 to reach API 2 once, hit %d times", counterReqAPI2)
	}

	if w.Result().StatusCode != nethttp.StatusOK {
		t.Errorf("expected status code to be 200 got %v", w.Result().StatusCode)
	}

	req = httptest.NewRequest(nethttp.MethodPost, "/posts", nil)
	w = httptest.NewRecorder()

	s.ServeHTTP(w, req)

	if counterReqAPI1 != 1 {
		t.Errorf("expected GET /posts/1 to reach API 2 once, hit %d times", counterReqAPI2)
	}

	if w.Result().StatusCode != nethttp.StatusCreated {
		t.Errorf("expected status code to be 201 got %v", w.Result().StatusCode)
	}
}
