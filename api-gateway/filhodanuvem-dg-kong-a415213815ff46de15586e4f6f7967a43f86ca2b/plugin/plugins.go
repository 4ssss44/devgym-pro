package plugin

import (
	"bytes"
	"errors"
	"fmt"
	"io"
	"log/slog"
	"net/http"
	"strings"

	"github.com/devgymbr/kong"
	"github.com/golang-jwt/jwt/v5"
)

var availablePlugins = map[string]Middleware{}

type Middleware func(p kong.Plugin, f http.HandlerFunc) http.HandlerFunc

func Register(pluginName string, middleware Middleware) {
	availablePlugins[pluginName] = middleware
}

func FindMiddleware(pluginName string) (Middleware, error) {
	if middleware, ok := availablePlugins[pluginName]; ok {
		return middleware, nil
	}

	return nil, kong.ErrPluginNotFound
}

func Log(p kong.Plugin, f http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		slog.Info("request received",
			slog.String("method", r.Method),
			slog.String("url", r.URL.Path),
		)
		f(w, r)
	}
}

func AddHeader(p kong.Plugin, f http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		for header, value := range p.Input {
			r.Header.Add(header, value.(string))
		}

		f(w, r)
	}
}

func RequestSizeLimit(p kong.Plugin, f http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		body, err := io.ReadAll(r.Body)
		if err != nil {
			slog.Error("could not read request body", slog.String("error", err.Error()))
		}

		if p.Input["allowed_payload_size"] == nil || len(body) >= p.Input["allowed_payload_size"].(int) {
			w.WriteHeader(http.StatusRequestEntityTooLarge)
			return
		}

		// nós precisamos fechar o body
		// e garantir que da próxima vez que ele seja lido, que seja do início novamente
		// então fazemos reset dele
		// leia: https://stackoverflow.com/questions/46948050/how-to-read-request-body-twice-in-golang-middleware
		// outro approach seria criar uma nova request do zero e copiar os dados da request original
		r.Body.Close()
		r.Body = io.NopCloser(bytes.NewBuffer(body))

		f(w, r)
	}
}

func getToken(p kong.Plugin, r *http.Request) (string, error) {
	keyName := fmt.Sprintf("%s", p.Input["key_name"])
	if p.Input["key_in_header"] != nil && p.Input["key_in_header"].(bool) {
		// Expected header value in format "Bearer <token>""
		header := r.Header.Get(keyName)
		parts := strings.Split(header, " ")
		if len(parts) < 2 {
			return "", errors.New("invalid header format")
		}

		return parts[1], nil
	}

	if p.Input["key_in_query"] != nil && p.Input["key_in_query"].(bool) {
		return r.URL.Query().Get(keyName), nil
	}

	return "", errors.New("one of the key_in_header or key_in_query must be true")
}

func JWTAuth(p kong.Plugin, f http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		tokenString, err := getToken(p, r)
		if err != nil {
			slog.Error("could not get token", slog.String("error", err.Error()))
			w.WriteHeader(http.StatusUnauthorized)
			return
		}

		_, err = jwt.Parse(tokenString, func(_ *jwt.Token) (interface{}, error) {
			value := fmt.Sprintf("%s", p.Input["secret"])
			return []byte(value), nil
		})

		if err != nil {
			slog.Error("could not parse token", slog.String("error", err.Error()))
			w.WriteHeader(http.StatusUnauthorized)
			return
		}

		f(w, r)
	}
}
