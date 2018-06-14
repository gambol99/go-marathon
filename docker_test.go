/*
Copyright 2015 The go-marathon Authors All rights reserved.

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

func createPortMapping(containerPort int, protocol string) *PortMapping {
	return &PortMapping{
		ContainerPort: containerPort,
		HostPort:      0,
		ServicePort:   0,
		Protocol:      protocol,
	}
}

func TestDockerAddParameter(t *testing.T) {
	docker := NewDockerApplication().Container.Docker
	docker.AddParameter("k1", "v1").AddParameter("k2", "v2")

	assert.Equal(t, 2, len(*docker.Parameters))
	assert.Equal(t, (*docker.Parameters)[0].Key, "k1")
	assert.Equal(t, (*docker.Parameters)[0].Value, "v1")
	assert.Equal(t, (*docker.Parameters)[1].Key, "k2")
	assert.Equal(t, (*docker.Parameters)[1].Value, "v2")

	docker.EmptyParameters()
	assert.NotNil(t, docker.Parameters)
	assert.Equal(t, 0, len(*docker.Parameters))
}
func TestDockerExpose(t *testing.T) {
	apps := []*Application{
		NewDockerApplication(),
		NewDockerApplication(),
	}

	// Marathon < 1.5
	apps[0].Container.Docker.Expose(8080).Expose(80, 443)

	// Marathon >= 1.5
	apps[1].Container.Expose(8080).Expose(80, 443)

	portMappings := []*[]PortMapping{
		apps[0].Container.Docker.PortMappings,
		apps[1].Container.PortMappings,
	}

	for _, portMapping := range portMappings {
		assert.Equal(t, 3, len(*portMapping))

		assert.Equal(t, *createPortMapping(8080, "tcp"), (*portMapping)[0])
		assert.Equal(t, *createPortMapping(80, "tcp"), (*portMapping)[1])
		assert.Equal(t, *createPortMapping(443, "tcp"), (*portMapping)[2])
	}
}

func TestDockerExposeUDP(t *testing.T) {
	apps := []*Application{
		NewDockerApplication(),
		NewDockerApplication(),
	}

	// Marathon < 1.5
	apps[0].Container.Docker.ExposeUDP(53).ExposeUDP(5060, 6881)

	// Marathon >= 1.5
	apps[1].Container.ExposeUDP(53).ExposeUDP(5060, 6881)

	portMappings := []*[]PortMapping{
		apps[0].Container.Docker.PortMappings,
		apps[1].Container.PortMappings,
	}

	for _, portMapping := range portMappings {
		assert.Equal(t, 3, len(*portMapping))
		assert.Equal(t, *createPortMapping(53, "udp"), (*portMapping)[0])
		assert.Equal(t, *createPortMapping(5060, "udp"), (*portMapping)[1])
		assert.Equal(t, *createPortMapping(6881, "udp"), (*portMapping)[2])
	}
}

func TestPortMappingLabels(t *testing.T) {
	pm := createPortMapping(80, "tcp")

	pm.AddLabel("hello", "world").AddLabel("foo", "bar")

	assert.Equal(t, 2, len(*pm.Labels))
	assert.Equal(t, "world", (*pm.Labels)["hello"])
	assert.Equal(t, "bar", (*pm.Labels)["foo"])

	pm.EmptyLabels()

	assert.NotNil(t, pm.Labels)
	assert.Equal(t, 0, len(*pm.Labels))
}

func TestPortMappingNetworkNames(t *testing.T) {
	pm := createPortMapping(80, "tcp")

	pm.AddNetwork("test")

	assert.Equal(t, 1, len(*pm.NetworkNames))
	assert.Equal(t, "test", (*pm.NetworkNames)[0])

	pm.EmptyNetworkNames()

	assert.NotNil(t, pm.NetworkNames)
	assert.Equal(t, 0, len(*pm.NetworkNames))
}

func TestVolume(t *testing.T) {
	container := NewDockerApplication().Container

	container.Volume("hp1", "cp1", "RW")
	container.Volume("hp2", "cp2", "R")

	assert.Equal(t, 2, len(*container.Volumes))
	assert.Equal(t, (*container.Volumes)[0].HostPath, "hp1")
	assert.Equal(t, (*container.Volumes)[0].ContainerPath, "cp1")
	assert.Equal(t, (*container.Volumes)[0].Mode, "RW")
	assert.Equal(t, (*container.Volumes)[1].HostPath, "hp2")
	assert.Equal(t, (*container.Volumes)[1].ContainerPath, "cp2")
	assert.Equal(t, (*container.Volumes)[1].Mode, "R")
}

func TestExternalVolume(t *testing.T) {
	container := NewDockerApplication().Container

	container.Volume("", "cp", "RW")
	ev := (*container.Volumes)[0].SetExternalVolume("myVolume", "dvdi")

	ev.AddOption("prop", "pval")
	ev.AddOption("dvdi", "rexray")

	ev1 := (*container.Volumes)[0].External
	assert.Equal(t, ev1.Name, "myVolume")
	assert.Equal(t, ev1.Provider, "dvdi")
	if assert.Equal(t, len(*ev1.Options), 2) {
		assert.Equal(t, (*ev1.Options)["dvdi"], "rexray")
		assert.Equal(t, (*ev1.Options)["prop"], "pval")
	}

	// empty the external volume again
	(*container.Volumes)[0].EmptyExternalVolume()
	ev2 := (*container.Volumes)[0].External
	assert.Equal(t, ev2.Name, "")
	assert.Equal(t, ev2.Provider, "")
}

func TestDockerPersistentVolume(t *testing.T) {
	docker := NewDockerApplication()
	container := docker.Container.Volume("/host", "/container", "RW")
	require.Equal(t, 1, len(*docker.Container.Volumes))

	pVol := (*container.Volumes)[0].SetPersistentVolume()
	pVol.SetType(PersistentVolumeTypeMount)
	pVol.SetSize(256)
	pVol.SetMaxSize(128)
	pVol.AddConstraint("cons1", "EQUAL", "tag1")
	pVol.AddConstraint("cons2", "UNIQUE")

	assert.Equal(t, 256, pVol.Size)
	assert.Equal(t, PersistentVolumeTypeMount, pVol.Type)
	assert.Equal(t, 128, pVol.MaxSize)

	if assert.NotNil(t, pVol.Constraints) {
		constraints := *pVol.Constraints
		require.Equal(t, 2, len(constraints))
		assert.Equal(t, []string{"cons1", "EQUAL", "tag1"}, constraints[0])
		assert.Equal(t, []string{"cons2", "UNIQUE"}, constraints[1])
	}

	pVol.EmptyConstraints()
	if assert.NotNil(t, pVol.Constraints) {
		assert.Empty(t, len(*pVol.Constraints))
	}
}
