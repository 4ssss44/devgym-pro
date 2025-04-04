package integration

import (
	"context"
	"net/http"
	"testing"
	"time"

	"slices"

	kube "github.com/devgymbr/kubeclient"
	"github.com/devgymbr/kubeclient/deployment"
	"github.com/devgymbr/kubeclient/errors"
	"github.com/google/uuid"
)

const url = "http://localhost:3000"

func TestCreateDeployment(t *testing.T) {
	c, err := kube.NewClient(
		kube.WithURL(url),
	)
	if err != nil {
		t.Errorf("should not fail to create client: %s", err)
		return
	}

	deploy := deployment.Deployment{
		ID:       uuid.New(),
		Image:    "nginx",
		Replicas: 1,
		Ports: []deployment.Port{
			{
				Name:   "http",
				Number: 80,
			},
		},
	}

	createdDeploy, err := c.Deployment.Create(context.Background(), deploy)
	if err != nil {
		t.Errorf("should not fail to create deployment: %s", err)
		return
	}
	defer c.Deployment.Delete(context.Background(), deploy.ID)

	assertDeployment(t, &deploy, createdDeploy)
}

func TestCreateDeploymentWithShortTimeout(t *testing.T) {
	c, err := kube.NewClient(
		kube.WithURL(url),
	)
	if err != nil {
		t.Errorf("should not fail to create client: %s", err)
		return
	}

	deploy := deployment.Deployment{
		ID:       uuid.New(),
		Image:    "nginx",
		Replicas: 1,
		Ports: []deployment.Port{
			{
				Name:   "http",
				Number: 80,
			},
		},
	}

	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Nanosecond)
	defer cancel()
	_, err = c.Deployment.Create(ctx, deploy)
	if err == nil {
		t.Errorf("should fail due to a timeout")
		return
	}
	c.Deployment.Delete(context.Background(), deploy.ID)
}

func TestCreateDeploymentWithShortDefaultTimeout(t *testing.T) {
	c, err := kube.NewClient(
		kube.WithURL(url),
		kube.WithTimeout(1*time.Nanosecond),
	)
	if err != nil {
		t.Errorf("should not fail to create client: %s", err)
		return
	}

	deploy := deployment.Deployment{
		ID:       uuid.New(),
		Image:    "nginx",
		Replicas: 1,
		Ports: []deployment.Port{
			{
				Name:   "http",
				Number: 80,
			},
		},
	}

	_, err = c.Deployment.Create(context.Background(), deploy)
	if err == nil {
		t.Errorf("should fail due to a timeout")
	}
	c.Deployment.Delete(context.Background(), deploy.ID)
}

func TestCreateNonValidDeployment(t *testing.T) {
	c, err := kube.NewClient(
		kube.WithURL(url),
	)
	if err != nil {
		t.Errorf("should not fail to create client: %s", err)
		return
	}

	cases := map[string]struct {
		deployment   deployment.Deployment
		failedFields []string
	}{
		"empty deployment": {
			deployment:   deployment.Deployment{},
			failedFields: []string{"id", "replicas", "image", "ports"},
		},
		"invalid image": {
			deployment: deployment.Deployment{
				ID:       uuid.New(),
				Image:    "", // invalid image
				Replicas: 1,
				Ports: []deployment.Port{
					{
						Name:   "http",
						Number: 80,
					},
				},
			},
			failedFields: []string{"image"},
		},
		"without ports": {
			deployment: deployment.Deployment{
				ID:       uuid.New(),
				Image:    "nginx",
				Replicas: 1,
				Ports:    []deployment.Port{}, // without ports
			},
			failedFields: []string{"ports"},
		},
		"no replicas": {
			deployment: deployment.Deployment{
				ID:       uuid.New(),
				Image:    "nginx",
				Replicas: 0, // no replicas
				Ports: []deployment.Port{
					{
						Name:   "http",
						Number: 80,
					},
				},
			},
			failedFields: []string{"replicas"},
		},
	}

	for title, v := range cases {
		t.Run(title, func(t *testing.T) {
			_, err := c.Deployment.Create(context.Background(), v.deployment)
			if err == nil {
				t.Errorf("should fail to create deployment")
				return
			}

			if _, ok := err.(errors.InvalidResource); !ok {
				t.Errorf("should fail with InvalidResource error: %s", err)
				return
			}

			invalid := err.(errors.InvalidResource)
			if !slices.Equal(invalid.FailedFields, v.failedFields) {
				t.Errorf("should fail with %v failed fields: %v", v.failedFields, invalid.FailedFields)
				return
			}
		})
	}
}

func TestCreateDuplicatedDeployment(t *testing.T) {
	c, err := kube.NewClient(
		kube.WithURL(url),
	)
	if err != nil {
		t.Errorf("should not fail to create client: %s", err)
		return
	}

	deploy := deployment.Deployment{
		ID:       uuid.New(),
		Image:    "nginx",
		Replicas: 1,
		Ports: []deployment.Port{
			{
				Name:   "http",
				Number: 80,
			},
		},
	}

	_, err = c.Deployment.Create(context.Background(), deploy)
	if err != nil {
		t.Errorf("should not fail when creating first deployment: %s", err)
		return
	}
	defer c.Deployment.Delete(context.Background(), deploy.ID)

	_, err = c.Deployment.Create(context.Background(), deploy)
	if err == nil {
		t.Error("should fail to create duplicated deployment")
		return
	}

	if _, ok := err.(errors.ResponseError); !ok {
		t.Errorf("should fail with ResponseError error")
		return
	}

	if err.(errors.ResponseError).StatusCode != http.StatusConflict {
		t.Errorf("should fail with 409 error code")
		return
	}
}

func assertDeployment(t *testing.T, expectedDeploy, foundDeploy *deployment.Deployment) {
	if foundDeploy.ID != expectedDeploy.ID {
		t.Errorf("should have same ID: %s", foundDeploy.ID)
		return
	}

	if foundDeploy.Image != expectedDeploy.Image {
		t.Errorf("should have same Image: %s", foundDeploy.Image)
		return
	}

	if foundDeploy.Replicas != expectedDeploy.Replicas {
		t.Errorf("should have same Replicas: %d", foundDeploy.Replicas)
		return
	}

	if !foundDeploy.CreatedAt.IsZero() {
		t.Errorf("should have CreatedAt set")
		return
	}

	if len(foundDeploy.Ports) != len(expectedDeploy.Ports) {
		t.Error("should have same number of ports")
		return
	}

	for i, p := range foundDeploy.Ports {
		if p.Number != expectedDeploy.Ports[i].Number {
			t.Errorf("should have same port number: %d", p.Number)
			return
		}
	}
}
