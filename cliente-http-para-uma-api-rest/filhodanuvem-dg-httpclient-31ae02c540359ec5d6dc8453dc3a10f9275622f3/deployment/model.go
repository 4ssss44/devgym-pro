package deployment

import (
	"time"

	"github.com/google/uuid"
)

type Deployment struct {
	ID        uuid.UUID         `json:"id"`
	Labels    map[string]string `json:"labels"`
	Replicas  int               `json:"replicas"`
	Image     string            `json:"image"`
	Name      string            `json:"name"`
	Ports     []Port            `json:"ports"`
	CreatedAt time.Time         `json:"created_at"`
}

type Port struct {
	Name   string `json:"name"`
	Number uint   `json:"port"`
}

type Error struct {
	Message string `json:"message"`
	Code    int    `json:"code"`
}
