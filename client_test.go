/*
Copyright 2014 Rohith All rights reserved.

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
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestPing(t *testing.T) {
	endpoint := newFakeMarathonEndpoint(t, nil)
	defer endpoint.Close()

	pong, err := endpoint.Client.Ping()
	assert.NoError(t, err)
	assert.True(t, pong)
}

func TestGetMarathonURL(t *testing.T) {
	endpoint := newFakeMarathonEndpoint(t, nil)
	defer endpoint.Close()

	assert.Equal(t, endpoint.Client.GetMarathonURL(), endpoint.URL)
}

func TestOneLogLine(t *testing.T) {
	in := `
	a
	b    c
	d\n
	  efgh
	i\r\n
	j\t
	{"json":  "works",
		"f o o": "ba    r"
	}
	`
	assert.Equal(t, `a\n b    c\n d\n\n efgh\n i\r\n\n j\t\n {"json":  "works",\n "f o o": "ba    r"\n }\n `, string(oneLogLine([]byte(in))))
}
