package notion

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
)

// See: https://developers.notion.com/reference/errors.
var (
	ErrInvalidJSON        = errors.New("notion: request body could not be decoded as JSON")
	ErrInvalidRequestURL  = errors.New("notion: request URL is not valid")
	ErrInvalidRequest     = errors.New("notion: request is not supported")
	ErrValidation         = errors.New("notion: request body does not match the schema for the expected parameters")
	ErrUnauthorized       = errors.New("notion: bearer token is not valid")
	ErrRestrictedResource = errors.New("notion: client doesn't have permission to perform this operation")
	ErrObjectNotFound     = errors.New("notion: the resource does not exist")
	ErrConflict           = errors.New("notion: the transaction could not be completed, potentially due to a data collision")
	ErrRateLimited        = errors.New("notion: this request exceeds the number of requests allowed")
	ErrInternalServer     = errors.New("notion: an unexpected error occurred")
	ErrServiceUnavailable = errors.New("notion: service is unavailable")
)

type APIError struct {
	Object  string `json:"object"`
	Status  int    `json:"status"`
	Code    string `json:"code"`
	Message string `json:"message"`
}

// Error implements `error`.
func (err *APIError) Error() string {
	return fmt.Sprintf("%v (code: %v, status: %v)", err.Message, err.Code, err.Status)
}

func (err *APIError) Unwrap() error {
	switch err.Code {
	case "invalid_json":
		return ErrInvalidJSON
	case "invalid_request_url":
		return ErrInvalidRequestURL
	case "invalid_request":
		return ErrInvalidRequest
	case "validation_error":
		return ErrValidation
	case "unauthorized":
		return ErrUnauthorized
	case "restricted_resource":
		return ErrRestrictedResource
	case "object_not_found":
		return ErrObjectNotFound
	case "conflict_error":
		return ErrConflict
	case "rate_limited":
		return ErrRateLimited
	case "internal_server_error":
		return ErrInternalServer
	case "service_unavailable":
		return ErrServiceUnavailable
	default:
		return fmt.Errorf("notion: unknown error (%v)", err.Code)
	}
}

func parseErrorResponse(res *http.Response) error {
	var apiErr APIError

	err := json.NewDecoder(res.Body).Decode(&apiErr)
	if err != nil {
		return fmt.Errorf("failed to parse error from HTTP response: %w", err)
	}

	return &apiErr
}
