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
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestGetPodStatus(t *testing.T) {
	endpoint := newFakeMarathonEndpoint(t, nil)
	defer endpoint.Close()

	podStatus, err := endpoint.Client.PodStatus(fakePodName)
	require.NoError(t, err)

	if assert.NotNil(t, podStatus) {
		assert.Equal(t, podStatus.Spec.ID, fakePodName)
	}
}

func TestGetAllPodStatus(t *testing.T) {
	endpoint := newFakeMarathonEndpoint(t, nil)
	defer endpoint.Close()

	podStatuses, err := endpoint.Client.PodStatuses()
	require.NoError(t, err)
	assert.Equal(t, podStatuses[0].Spec.ID, fakePodName)
}

func TestWaitOnPod(t *testing.T) {
	endpoint := newFakeMarathonEndpoint(t, nil)
	defer endpoint.Close()

	err := endpoint.Client.WaitOnPod(fakePodName, 1*time.Microsecond)
	require.NoError(t, err)
}

func TestPodIsRunning(t *testing.T) {
	endpoint := newFakeMarathonEndpoint(t, nil)
	defer endpoint.Close()

	exists := endpoint.Client.PodIsRunning(fakePodName)
	assert.True(t, exists)

	exists = endpoint.Client.PodIsRunning("not_existing")
	assert.False(t, exists)

	exists = endpoint.Client.PodIsRunning(secondFakePodName)
	assert.False(t, exists)
}
