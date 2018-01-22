/*
Copyright 2017 The go-marathon Authors All rights reserved.

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
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestPodEnvironmentVariableUnmarshal(t *testing.T) {
	endpoint := newFakeMarathonEndpoint(t, nil)
	defer endpoint.Close()

	pod, err := endpoint.Client.Pod(fakePodName)
	require.NoError(t, err)

	env := pod.Env
	secrets := pod.Secrets

	require.NotNil(t, env)
	assert.Equal(t, "value", env["key1"])
	assert.Equal(t, "key2", secrets["secret0"].EnvVar)
	assert.Equal(t, "source0", secrets["secret0"].Source)

	assert.Equal(t, "value3", pod.Containers[0].Env["key3"])
	assert.Equal(t, "key4", pod.Containers[0].Secrets["secret1"].EnvVar)
	assert.Equal(t, "source1", secrets["secret1"].Source)
}

func TestPodMalformedPayloadUnmarshal(t *testing.T) {
	var tests = []struct {
		expected    string
		given       []byte
		description string
	}{
		{
			expected:    "unexpected secret field",
			given:       []byte(`{"environment": {"FOO": "bar", "SECRET": {"not_secret": "secret1"}}, "secrets": {"secret1": {"source": "/path/to/secret"}}}`),
			description: "Field in environment secret not equal to secret.",
		},
		{
			expected:    "unexpected secret field",
			given:       []byte(`{"environment": {"FOO": "bar", "SECRET": {"secret": 1}}, "secrets": {"secret1": {"source": "/path/to/secret"}}}`),
			description: "Invalid value in environment secret.",
		},
		{
			expected:    "unexpected environment variable type",
			given:       []byte(`{"environment": {"FOO": 1, "SECRET": {"secret": "secret1"}}, "secrets": {"secret1": {"source": "/path/to/secret"}}}`),
			description: "Invalid environment variable type.",
		},
		{
			expected:    "malformed pod definition",
			given:       []byte(`{"environment": "value"}`),
			description: "Bad pod definition.",
		},
	}

	for _, test := range tests {
		tmpPod := new(Pod)

		err := json.Unmarshal(test.given, &tmpPod)
		if assert.Error(t, err, test.description) {
			assert.True(t, strings.HasPrefix(err.Error(), test.expected), test.description)
		}
	}
}

func TestPodEnvironmentVariableMarshal(t *testing.T) {
	testPod := new(Pod)
	targetString := []byte(`{"containers":[{"lifecycle":{},"environment":{"FOO2":"bar2","TOP2":"secret1"}}],"environment":{"FOO":"bar","TOP":{"secret":"secret1"}},"secrets":{"secret1":{"source":"/path/to/secret"}}}`)

	testPod.AddEnv("FOO", "bar")
	testPod.AddSecret("TOP", "secret1", "/path/to/secret")

	testContainer := new(PodContainer)
	testContainer.AddSecret("TOP2", "secret1")
	testContainer.AddEnv("FOO2", "bar2")
	testPod.AddContainer(testContainer)

	pod, err := json.Marshal(testPod)
	if assert.NoError(t, err) {
		assert.Equal(t, targetString, pod)
	}
}
