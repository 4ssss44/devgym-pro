package integration

import (
	"context"
	"errors"
	"testing"
	"time"

	kube "github.com/devgymbr/kubeclient"
	"github.com/devgymbr/kubeclient/deployment"
	"github.com/google/uuid"
)

func TestDeleteDeployment(t *testing.T) {
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

	if err := c.Deployment.Delete(context.Background(), deploy.ID); err != nil {
		t.Errorf("should not fail to delete deployment: %s", err)
		return
	}
}

func TestDeleteNotFoundDeployment(t *testing.T) {
	c, err := kube.NewClient(
		kube.WithURL(url),
	)
	if err != nil {
		t.Errorf("should not fail to create client: %s", err)
		return
	}

	if err := c.Deployment.Delete(context.Background(), uuid.New()); err == nil {
		t.Errorf("should fail to delete non existent deployment")
		return
	}
}

func TestCancelDeleteDeploymentDueTimeout(t *testing.T) {
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

	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Nanosecond)
	defer cancel()
	err = c.Deployment.Delete(ctx, deploy.ID)
	if err == nil {
		t.Errorf("should fail to delete deployment due timeout")
		return
	}

	if errors.Is(err, context.Canceled) {
		t.Errorf("should fail with context.Canceled")
		return
	}
}
