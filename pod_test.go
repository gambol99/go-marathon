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

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

const key = "testKey"
const val = "testValue"

const fakePodName = "/fake-pod"
const secondFakePodName = "/fake-pod2"

func TestPodLabels(t *testing.T) {
	pod := NewPod()
	pod.AddLabel(key, val)
	if assert.Equal(t, len(pod.Labels), 1) {
		assert.Equal(t, pod.Labels[key], val)
	}

	pod.EmptyLabels()
	assert.Equal(t, len(pod.Labels), 0)
}

func TestPodEnvironmentVars(t *testing.T) {
	pod := NewPod()
	pod.AddEnv(key, val)

	newVal, ok := pod.Env[key]
	assert.Equal(t, newVal, val)
	assert.Equal(t, ok, true)

	badVal, ok := pod.Env["fakeKey"]
	assert.Equal(t, badVal, "")
	assert.Equal(t, ok, false)

	pod.EmptyEnvs()
	assert.Equal(t, len(pod.Env), 0)
}

func TestSecrets(t *testing.T) {
	pod := NewPod()
	pod.AddSecret("randomVar", key, val)

	newVal, err := pod.GetSecretSource(key)
	assert.Equal(t, newVal, val)
	assert.Equal(t, err, nil)

	badVal, err := pod.GetSecretSource("fakeKey")
	assert.Equal(t, badVal, "")
	assert.NotNil(t, err)

	pod.EmptySecrets()
	assert.Equal(t, len(pod.Env), 0)
}

func TestSupportsPod(t *testing.T) {
	endpoint := newFakeMarathonEndpoint(t, nil)

	supports, err := endpoint.Client.SupportsPods()
	if assert.Nil(t, err) {

		assert.Equal(t, supports, true)
	}

	// Manually closing to test lack of support
	endpoint.Close()

	supports, err = endpoint.Client.SupportsPods()
	assert.NotNil(t, err)
	assert.Equal(t, supports, false)
}
func TestGetPod(t *testing.T) {
	endpoint := newFakeMarathonEndpoint(t, nil)
	defer endpoint.Close()

	pod, err := endpoint.Client.Pod(fakePodName)
	require.NoError(t, err)
	if assert.NotNil(t, pod) {
		assert.Equal(t, pod.ID, fakePodName)
	}
}

func TestGetAllPods(t *testing.T) {
	endpoint := newFakeMarathonEndpoint(t, nil)
	defer endpoint.Close()

	pods, err := endpoint.Client.Pods()
	require.NoError(t, err)
	if assert.Equal(t, len(pods), 2) {
		assert.Equal(t, pods[0].ID, fakePodName)
		assert.Equal(t, pods[1].ID, secondFakePodName)
	}
}

func TestCreatePod(t *testing.T) {
	endpoint := newFakeMarathonEndpoint(t, nil)
	defer endpoint.Close()

	pod := NewPod().Name(fakePodName)
	pod, err := endpoint.Client.CreatePod(pod)
	require.NoError(t, err)
	if assert.NotNil(t, pod) {
		assert.Equal(t, pod.ID, fakePodName)
	}
}

func TestUpdatePod(t *testing.T) {
	endpoint := newFakeMarathonEndpoint(t, nil)
	defer endpoint.Close()

	pod := NewPod().Name(fakePodName)
	pod, err := endpoint.Client.CreatePod(pod)
	require.NoError(t, err)

	pod, err = endpoint.Client.UpdatePod(pod, true)
	require.NoError(t, err)

	if assert.NotNil(t, pod) {
		assert.Equal(t, pod.ID, fakePodName)
		assert.Equal(t, pod.Scaling.Instances, 2)
	}
}

func TestDeletePod(t *testing.T) {
	endpoint := newFakeMarathonEndpoint(t, nil)
	defer endpoint.Close()

	id, err := endpoint.Client.DeletePod(fakePodName, true)
	require.NoError(t, err)

	if assert.NotNil(t, id) {
		assert.Equal(t, id.DeploymentID, "c0e7434c-df47-4d23-99f1-78bd78662231")
	}
}

func TestVersions(t *testing.T) {
	endpoint := newFakeMarathonEndpoint(t, nil)
	defer endpoint.Close()

	versions, err := endpoint.Client.PodVersions(fakePodName)
	require.NoError(t, err)
	assert.Equal(t, versions[0], "2014-08-18T22:36:41.451Z")
}

func TestGetPodByVersion(t *testing.T) {
	endpoint := newFakeMarathonEndpoint(t, nil)
	defer endpoint.Close()

	pod, err := endpoint.Client.PodByVersion(fakePodName, "2014-08-18T22:36:41.451Z")
	require.NoError(t, err)
	assert.Equal(t, pod.ID, fakePodName)
}
