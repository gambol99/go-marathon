/*
Copyright 2015 Rohith All rights reserved.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package marathon

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
)

const (
	// ErrCodeBadRequest specifies a 400 Bad Request error.
	ErrCodeBadRequest = iota
	// ErrCodeUnauthorized specifies a 401 Unauthorized error.
	ErrCodeUnauthorized
	// ErrCodeForbidden specifies a 403 Forbidden error.
	ErrCodeForbidden
	// ErrCodeNotFound specifies a 404 Not Found error.
	ErrCodeNotFound
	// ErrCodeDuplicateID specifies a PUT 409 Conflict error.
	ErrCodeDuplicateID
	// ErrCodeAppLocked specifies a POST 409 Conflict error.
	ErrCodeAppLocked
	// ErrCodeInvalidBean specifies a 422 UnprocessableEntity error.
	ErrCodeInvalidBean
	// ErrCodeServer specifies a 500+ Server error.
	ErrCodeServer
	// ErrCodeUnknown specifies an unknown error.
	ErrCodeUnknown
)

// APIError represents a generic API error.
type APIError struct {
	// ErrCode specifies the nature of the error.
	ErrCode int
	message string
}

func (e *APIError) Error() string {
	return fmt.Sprintf("Marathon API error: %s", e.message)
}

// NewAPIError creates a new APIError instance from the given response code and content.
func NewAPIError(code int, content []byte) (*APIError, error) {
	switch {
	case code == http.StatusBadRequest:
		return parseContent(&badRequestDef{}, content)
	case code == http.StatusUnauthorized:
		return parseContent(&simpleErrDef{code: ErrCodeUnauthorized}, content)
	case code == http.StatusForbidden:
		return parseContent(&simpleErrDef{code: ErrCodeForbidden}, content)
	case code == http.StatusNotFound:
		return parseContent(&simpleErrDef{code: ErrCodeNotFound}, content)
	case code == http.StatusConflict:
		return parseContent(&conflictDef{}, content)
	case code == 422:
		return parseContent(&unprocessableEntityDef{}, content)
	case code >= http.StatusInternalServerError:
		return parseContent(&simpleErrDef{code: ErrCodeServer}, content)
	default:
		return &APIError{ErrCodeUnknown, "unknown error"}, nil
	}
}

type errorDefinition interface {
	message() string
	errCode() int
}

func parseContent(errDef errorDefinition, content []byte) (*APIError, error) {
	if err := json.Unmarshal(content, errDef); err != nil {
		return nil, err
	}

	return &APIError{message: errDef.message(), ErrCode: errDef.errCode()}, nil
}

type simpleErrDef struct {
	Message string `json:"message"`
	code    int
}

func (def *simpleErrDef) message() string {
	return def.Message
}

func (def *simpleErrDef) errCode() int {
	return def.code
}

type badRequestDef struct {
	Message string `json:"message"`
	Details []struct {
		Path   string   `json:"path"`
		Errors []string `json:"errors"`
	} `json:"details"`
}

func (def *badRequestDef) message() string {
	var details []string
	for _, detail := range def.Details {
		errDesc := fmt.Sprintf("path: '%s' errors: %s", detail.Path,
			strings.Join(detail.Errors, ", "))
		details = append(details, errDesc)
	}

	return fmt.Sprintf("%s (%s)", def.Message, strings.Join(details, "; "))
}

func (def *badRequestDef) errCode() int {
	return ErrCodeBadRequest
}

type conflictDef struct {
	Message     string `json:"message"`
	Deployments []struct {
		ID string `json:"id"`
	} `json:"deployments"`
}

func (def *conflictDef) message() string {
	if len(def.Deployments) == 0 {
		// 409 Conflict response to "POST /v2/apps".
		return def.Message
	}

	// 409 Conflict response to "PUT /v2/apps/{appId}".
	var ids []string
	for _, deployment := range def.Deployments {
		ids = append(ids, deployment.ID)
	}
	return fmt.Sprintf("%s (locking deployment IDs: %s)", def.Message, strings.Join(ids, ", "))
}

func (def *conflictDef) errCode() int {
	if len(def.Deployments) == 0 {
		return ErrCodeDuplicateID
	}

	return ErrCodeAppLocked
}

type unprocessableEntityDetails []struct {
	Attribute string `json:"attribute"`
	Error     string `json:"error"`
}

type unprocessableEntityDef struct {
	Message string `json:"message"`
	// Name used in Marathon 0.15.0+.
	Details unprocessableEntityDetails `json:"details"`
	// Name used in Marathon < 0.15.0.
	Errors unprocessableEntityDetails `json:"errors"`
}

func (def *unprocessableEntityDef) message() string {
	joinDetails := func(details unprocessableEntityDetails) []string {
		var res []string
		for _, detail := range details {
			res = append(res,
				fmt.Sprintf("attribute '%s': %s", detail.Attribute, detail.Error))
		}
		return res
	}

	details := joinDetails(def.Details)
	details = append(details, joinDetails(def.Errors)...)

	return fmt.Sprintf("%s (%s)", def.Message, strings.Join(details, "; "))
}

func (def *unprocessableEntityDef) errCode() int {
	return ErrCodeInvalidBean
}
