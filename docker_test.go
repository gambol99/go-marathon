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
	"testing"

	"github.com/stretchr/testify/assert"
)

func createPortMapping(containerPort int, protocol string) *PortMapping {
	return &PortMapping{
		ContainerPort: containerPort,
		HostPort:      0,
		ServicePort:   0,
		Protocol:      protocol,
	}
}

func TestDockerExpose(t *testing.T) {
	app := NewDockerApplication()
	app.Container.Docker.Expose(8080).Expose(80, 443)

	portMappings := app.Container.Docker.PortMappings
	assert.Equal(t, 3, len(portMappings))
	assert.Equal(t, createPortMapping(8080, "tcp"), portMappings[0])
	assert.Equal(t, createPortMapping(80, "tcp"), portMappings[1])
	assert.Equal(t, createPortMapping(443, "tcp"), portMappings[2])
}

func TestDockerExposeUDP(t *testing.T) {
	app := NewDockerApplication()
	app.Container.Docker.ExposeUDP(53).ExposeUDP(5060, 6881)

	portMappings := app.Container.Docker.PortMappings
	assert.Equal(t, 3, len(portMappings))
	assert.Equal(t, createPortMapping(53, "udp"), portMappings[0])
	assert.Equal(t, createPortMapping(5060, "udp"), portMappings[1])
	assert.Equal(t, createPortMapping(6881, "udp"), portMappings[2])
}
