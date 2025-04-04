package deployment

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/devgymbr/kubeclient/errors"
	"github.com/google/uuid"
)

type Service struct {
	client  *http.Client
	urlBase string
}

func NewService(client *http.Client, url string) Service {
	return Service{
		client:  client,
		urlBase: url,
	}
}

func (s *Service) Create(ctx context.Context, deploy Deployment) (*Deployment, error) {
	j, err := json.Marshal(deploy)
	if err != nil {
		return nil, err
	}

	endpoint := fmt.Sprintf("%s/deployments", s.urlBase)

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, endpoint, bytes.NewBuffer(j))
	if err != nil {
		return nil, err
	}

	req.Header.Set("Content-Type", "application/json")
	resp, err := s.client.Do(req)

	if err != nil {
		return nil, err
	}

	defer resp.Body.Close()

	if resp.StatusCode == http.StatusBadRequest {
		return nil, errors.FromBadRequest(resp)
	}

	if resp.StatusCode != http.StatusCreated {
		return nil, errors.FromHTTPResponse(resp)
	}

	createdDeploy := Deployment{}
	if err := json.NewDecoder(resp.Body).Decode(&createdDeploy); err != nil {
		return nil, err
	}

	return &createdDeploy, nil
}

func (s *Service) Delete(ctx context.Context, ID uuid.UUID) error {
	endpoint := fmt.Sprintf("%s/deployments/%s", s.urlBase, ID.String())

	req, err := http.NewRequestWithContext(ctx, http.MethodDelete, endpoint, nil)
	if err != nil {
		return err
	}

	resp, err := s.client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusNoContent {
		return errors.FromHTTPResponse(resp)
	}

	return nil
}

func (s *Service) Get(ctx context.Context, ID uuid.UUID) (*Deployment, error) {
	endpoint := fmt.Sprintf("%s/deployments/%s", s.urlBase, ID.String())

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, endpoint, nil)
	if err != nil {
		return nil, err
	}

	resp, err := s.client.Do(req)
	if err != nil {
		return nil, err
	}

	defer resp.Body.Close()

	if resp.StatusCode == http.StatusNotFound {
		return nil, errors.NewNotFound(ID, "deployment")
	}

	if resp.StatusCode != http.StatusOK {
		return nil, errors.FromHTTPResponse(resp)
	}

	deploy := Deployment{}
	if err := json.NewDecoder(resp.Body).Decode(&deploy); err != nil {
		return nil, err
	}

	return &deploy, nil
}
