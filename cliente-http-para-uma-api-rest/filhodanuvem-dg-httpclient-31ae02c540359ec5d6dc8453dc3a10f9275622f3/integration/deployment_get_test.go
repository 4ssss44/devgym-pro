package integration

import (
	"context"
	"testing"

	kube "github.com/devgymbr/kubeclient"
	"github.com/devgymbr/kubeclient/deployment"
	ierrors "github.com/devgymbr/kubeclient/errors"
	"github.com/google/uuid"
)

func TestGetDeployment(t *testing.T) {
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
		t.Errorf("should not fail to create deployment: %s", err)
		return
	}
	defer c.Deployment.Delete(context.Background(), deploy.ID)

	foundDeploy, err := c.Deployment.Get(context.Background(), deploy.ID)
	if err != nil {
		t.Errorf("should not fail to get deployment: %s", err)
		return
	}

	assertDeployment(t, &deploy, foundDeploy)
}

func TestGetNotFoundDeployment(t *testing.T) {
	c, err := kube.NewClient(
		kube.WithURL(url),
	)
	if err != nil {
		t.Errorf("should not fail to create client: %s", err)
		return
	}

	foundDeploy, err := c.Deployment.Get(context.Background(), uuid.New())
	if err == nil || foundDeploy != nil {
		t.Errorf("should fail to get non existent deployment")
		return
	}

	if _, ok := err.(ierrors.NotFoundResource); !ok {
		t.Errorf("should fail with ErrNotFound")
		return
	}
}
