package tests

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/devgymbr/kong"
	"github.com/devgymbr/kong/plugin"
	"github.com/golang-jwt/jwt/v5"
)

func TestPluginAddHeader(t *testing.T) {
	// request vai ter vários cabeçalhos por padrão
	headers := map[string]string{
		"X-custom1": "value-1",
		"X-custom2": "value-2",
	}

	p := kong.Plugin{
		Name: "add_header",
		Input: map[string]any{
			"X-youtuber": "filhodanuvem",
		},
	}

	// a função f vai ser o handler do http que testa se a request chegou com os cabeçalhos esperados
	f := func(w http.ResponseWriter, r *http.Request) {
		if r.Header.Get("X-youtuber") != "filhodanuvem" {
			t.Errorf("expected header X-youtuber to be filhodanuvem got %v", r.Header.Get("X-youtuber"))
		}

		for k, v := range headers {
			if r.Header.Get(k) != v {
				t.Errorf("expected header %v to be %v got %v", k, v, r.Header.Get(k))
			}
		}
	}

	f = plugin.AddHeader(p, f)

	// adiciona todos os cabeçalhos na request
	r := httptest.NewRequest(http.MethodGet, "/", nil)
	for k, v := range headers {
		r.Header.Set(k, v)
	}

	w := httptest.NewRecorder()

	f(w, r)
}

func TestPluginRequestSizeLimitFail(t *testing.T) {
	// request com limite máximo de 10 bytes
	p := kong.Plugin{
		Name: "request_size_limit",
		Input: map[string]any{
			"allowed_payload_size": 10,
		},
	}

	// a função f vai ser o handler do http e sempre falha porque request deveria ser bloqueada
	f := func(_ http.ResponseWriter, _ *http.Request) {
		t.Error("expected to return 413 before reaching this point")
	}

	f = plugin.RequestSizeLimit(p, f)

	// request greater than size limit
	r := httptest.NewRequest(http.MethodPost, "/", strings.NewReader("body size greater than 10 bytes"))
	w := httptest.NewRecorder()

	f(w, r)

	if w.Code != http.StatusRequestEntityTooLarge {
		t.Errorf("expected status code 413 got %v", w.Code)
	}
}

func TestPluginRequestSizeLimitSuccess(t *testing.T) {
	// request com limite máximo de 10 bytes
	p := kong.Plugin{
		Name: "request_size_limit",
		Input: map[string]any{
			"allowed_payload_size": 10,
		},
	}

	// a função f vai ser o handler do http e sempre retorna 200
	f := func(w http.ResponseWriter, _ *http.Request) {
		w.WriteHeader(http.StatusOK)
	}

	f = plugin.RequestSizeLimit(p, f)

	// request allowed based on size limit
	r := httptest.NewRequest(http.MethodPost, "/", strings.NewReader("smallbody"))
	w := httptest.NewRecorder()

	f(w, r)

	if w.Code != http.StatusOK {
		t.Errorf("expected status code 200 got %v", w.Code)
	}
}

func TestPluginJWTAuthSuccessWithTokenOnHeader(t *testing.T) {
	// request com jwt secret
	secret := "th1s1ss3cr3t"
	p := kong.Plugin{
		Name: "request_size_limit",
		Input: map[string]any{
			"secret":        secret,
			"key_in_header": true,
			"key_name":      "Authorization",
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{}, nil)
	tokenString, err := token.SignedString([]byte(secret))
	if err != nil {
		t.Error("could not sign jwt token: " + err.Error())
		t.FailNow()
	}

	// a função f vai ser o handler do http e sempre retorna 200
	f := func(w http.ResponseWriter, _ *http.Request) {
		w.WriteHeader(http.StatusOK)
	}

	f = plugin.JWTAuth(p, f)

	// request with Bearer token
	r := httptest.NewRequest(http.MethodPost, "/", nil)
	r.Header.Set("Authorization", "Bearer "+tokenString)
	w := httptest.NewRecorder()

	f(w, r)

	if w.Code != http.StatusOK {
		t.Errorf("expected status code 200 got %v", w.Code)
	}
}

func TestPluginJWTAuthSuccessWithTokenOnQueryString(t *testing.T) {
	// request com jwt secret
	secret := "th1s1ss3cr3t"
	p := kong.Plugin{
		Name: "request_size_limit",
		Input: map[string]any{
			"secret":       secret,
			"key_in_query": true,
			"key_name":     "jwt",
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{}, nil)
	tokenString, err := token.SignedString([]byte(secret))
	if err != nil {
		t.Error("could not sign jwt token: " + err.Error())
		t.FailNow()
	}

	// a função f vai ser o handler do http e sempre retorna 200
	f := func(w http.ResponseWriter, _ *http.Request) {
		w.WriteHeader(http.StatusOK)
	}

	f = plugin.JWTAuth(p, f)

	// request with token as query string
	r := httptest.NewRequest(http.MethodPost, "/?jwt="+tokenString, nil)
	w := httptest.NewRecorder()

	f(w, r)

	if w.Code != http.StatusOK {
		t.Errorf("expected status code 200 got %v", w.Code)
	}
}

func TestPluginJWTAuthFailWithWrongSecret(t *testing.T) {
	// request com jwt secret
	secret := "th1s1ss3cr3t"
	p := kong.Plugin{
		Name: "request_size_limit",
		Input: map[string]any{
			"secret":        secret,
			"key_in_header": true,
			"key_name":      "Authorization",
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{}, nil)
	tokenString, err := token.SignedString([]byte("wrongsecret"))
	if err != nil {
		t.Error("could not sign jwt token: " + err.Error())
		t.FailNow()
	}

	// a função f vai ser o handler do http e falha pois o middleware deveria bloquear a request
	f := func(_ http.ResponseWriter, _ *http.Request) {
		t.Error("should not reach this point because jwt token is invalid")
	}

	f = plugin.JWTAuth(p, f)

	// request with Bearer token
	r := httptest.NewRequest(http.MethodPost, "/", nil)
	r.Header.Set("Authorization", "Bearer "+tokenString)
	w := httptest.NewRecorder()

	f(w, r)

	if w.Code != http.StatusUnauthorized {
		t.Errorf("expected status code 200 got %v", w.Code)
	}
}

func TestPluginJWTAuthFailWithMissingToken(t *testing.T) {
	// request com jwt secret
	secret := "th1s1ss3cr3t"
	p := kong.Plugin{
		Name: "request_size_limit",
		Input: map[string]any{
			"secret":        secret,
			"key_in_header": true,
			"key_name":      "Authorization",
		},
	}

	// a função f vai ser o handler do http e falha pois o middleware deveria bloquear a request
	f := func(_ http.ResponseWriter, _ *http.Request) {
		t.Error("should not reach this point because jwt token is invalid")
	}

	f = plugin.JWTAuth(p, f)

	// request without token
	r := httptest.NewRequest(http.MethodPost, "/", nil)
	w := httptest.NewRecorder()

	f(w, r)

	if w.Code != http.StatusUnauthorized {
		t.Errorf("expected status code 200 got %v", w.Code)
	}
}
