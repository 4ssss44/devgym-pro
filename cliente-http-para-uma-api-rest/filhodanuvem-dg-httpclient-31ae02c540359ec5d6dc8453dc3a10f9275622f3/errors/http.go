package errors

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"github.com/google/uuid"
)

type ResponseError struct {
	StatusCode int
	Message    string `json:"message"`
	Code       int    `json:"code"`
}

type InvalidResource struct {
	ResponseError

	FailedFields []string `json:"failed_fields"`
}

type NotFoundResource struct {
	ID       uuid.UUID
	Resource string
}

type response struct {
	Message string              `json:"message"`
	Code    int                 `json:"code"`
	Extras  map[string][]string `json:"extras"`
}

func (err ResponseError) Error() string {
	return fmt.Sprintf("StatusCode=%d, Message=%s", err.StatusCode, err.Message)
}

func FromHTTPResponse(resp *http.Response) ResponseError {
	var r response
	body, _ := io.ReadAll(resp.Body)
	if err := json.Unmarshal(body, &r); err != nil {
		return ResponseError{}
	}

	return ResponseError{
		StatusCode: resp.StatusCode,
		Message:    r.Message,
		Code:       r.Code,
	}
}

func FromBadRequest(resp *http.Response) InvalidResource {
	var r response
	body, _ := io.ReadAll(resp.Body)
	if err := json.Unmarshal(body, &r); err != nil {
		return InvalidResource{}
	}

	extras := []string{}
	if failedFields, ok := r.Extras["failed_fields"]; ok {
		for _, f := range failedFields {
			extras = append(extras, f)
		}
	}

	return InvalidResource{
		ResponseError: ResponseError{
			StatusCode: resp.StatusCode,
			Message:    r.Message,
			Code:       r.Code,
		},
		FailedFields: extras,
	}
}

func NewNotFound(id uuid.UUID, resource string) error {
	return NotFoundResource{
		ID:       id,
		Resource: resource,
	}
}

func (e InvalidResource) Error() string {
	return fmt.Sprintf("StatusCode=%d, Message=%s, FailedFields=%v", e.StatusCode, e.Message, e.FailedFields)
}

func (e NotFoundResource) Error() string {
	return fmt.Sprintf("%s with ID=%s not found", e.Resource, e.ID)
}
