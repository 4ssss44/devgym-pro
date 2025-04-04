package http

import (
	"io"
	"net/http"
	"net/url"

	"github.com/devgymbr/kong"
)

type ForwardClient struct {
	Client *http.Client
	Config *kong.Config
}

func (c *ForwardClient) ForwardRequest(urlBase string, w http.ResponseWriter, r *http.Request) error {
	forwardURL := urlBase + r.URL.Path

	var err error
	r.URL, err = url.Parse(forwardURL)
	if err != nil {
		return err
	}

	// vamos trocar o destino da request
	// e pra isso precisamos limpar o requestURI.
	// Leia https://stackoverflow.com/questions/19595860/http-request-requesturi-field-when-making-request-in-go
	r.RequestURI = ""
	resp, err := c.Client.Do(r)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	for k := range resp.Header {
		w.Header().Set(k, resp.Header.Get(k))
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return err
	}

	w.WriteHeader(resp.StatusCode)
	w.Write(body)

	return nil
}
