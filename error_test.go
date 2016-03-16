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
	"fmt"
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test400Error(t *testing.T) {
	content := []byte(`{
	"message": "Invalid JSON",
	"details": [
		{
			"path": "/id",
		 	"errors": ["error.expected.jsstring", "error.something.else"]
		},
		{
			"path": "/name",
		 	"errors": ["error.not.inventive"]
		}
	]
}`)

	e := validatedAPIError(t, http.StatusBadRequest, content, false)

	assert.Equal(t, ErrCodeBadRequest, e.ErrCode)
	assert.Contains(t, e.Error(), "Invalid JSON (path: '/id' errors: error.expected.jsstring, error.something.else; path: '/name' errors: error.not.inventive)")
}

func Test401Error(t *testing.T) {
	content := []byte(`{"message": "invalid username or password."}`)

	e := validatedAPIError(t, http.StatusUnauthorized, content, false)

	assert.Equal(t, ErrCodeUnauthorized, e.ErrCode)
	assert.Contains(t, e.Error(), "invalid username or password")
}

func Test403Error(t *testing.T) {
	content := []byte(`{"message": "Not Authorized to perform this action!"}`)

	e := validatedAPIError(t, http.StatusForbidden, content, false)

	assert.Equal(t, ErrCodeForbidden, e.ErrCode)
	assert.Contains(t, e.Error(), "Not Authorized to perform this action!")
}

func Test404Error(t *testing.T) {
	content := []byte(`{"message": "App '/not_existent' does not exist"}`)

	e := validatedAPIError(t, http.StatusNotFound, content, false)

	assert.Equal(t, ErrCodeNotFound, e.ErrCode)
	assert.Contains(t, e.Error(), "App '/not_existent' does not exist")
}

func Test409POSTError(t *testing.T) {
	content := []byte(`{"message": "An app with id [/existing_app] already exists."}`)

	e := validatedAPIError(t, http.StatusConflict, content, false)

	assert.Equal(t, ErrCodeDuplicateID, e.ErrCode)
	assert.Contains(t, e.Error(), "An app with id [/existing_app] already exists.")
}

func Test409PUTError(t *testing.T) {
	content := []byte(`{"message":"App is locked", "deployments": [ { "id": "97c136bf-5a28-4821-9d94-480d9fbb01c8" } ] }`)

	e := validatedAPIError(t, http.StatusConflict, content, false)

	assert.Equal(t, ErrCodeAppLocked, e.ErrCode)
	assert.Contains(t, e.Error(), "App is locked (locking deployment IDs: 97c136bf-5a28-4821-9d94-480d9fbb01c8)")
}

func Test422Error(t *testing.T) {
	for _, detailsPropKey := range []string{"details", "errors"} {
		content := []byte(fmt.Sprintf(`{
	"message": "Something is not valid",
	"%s": [
		{
			"attribute": "upgradeStrategy.minimumHealthCapacity",
			"error": "is greater than 1"
		},
		{
			"attribute": "foobar",
			"error": "foo does not have enough bar"
		}
	]
}
`, detailsPropKey))

		e := validatedAPIError(t, 422, content, false)

		assert.Equal(t, ErrCodeInvalidBean, e.ErrCode)
		assert.Contains(t, e.Error(), "Something is not valid (attribute 'upgradeStrategy.minimumHealthCapacity': is greater than 1; attribute 'foobar': foo does not have enough bar)")
	}
}

func TestServerError(t *testing.T) {
	content := []byte(`{"message": "internal server error"}`)

	for _, code := range []int{500, 501} {
		e := validatedAPIError(t, code, content, false)

		assert.Equal(t, ErrCodeServer, e.ErrCode, "code: %d", code)
		assert.Contains(t, e.Error(), "internal server error")
	}
}

func TestUnknownError(t *testing.T) {
	content := []byte("unknown error")

	e := validatedAPIError(t, 499, content, false)

	assert.Equal(t, ErrCodeUnknown, e.ErrCode)
	assert.Contains(t, e.Error(), "unknown error")
}

func TestInvalidJSON(t *testing.T) {
	content := []byte{}
	for _, code := range []int{400, 401, 403, 404, 409, 422, 500, 501} {
		validatedAPIError(t, code, content, true)
	}
}

func validatedAPIError(t *testing.T, code int, content []byte, parseErr bool) *APIError {
	e, err := NewAPIError(code, content)
	if parseErr {
		assert.Error(t, err, "code: %d", code)
	} else {
		assert.NoError(t, err, "code: %d", code)
	}

	return e
}
