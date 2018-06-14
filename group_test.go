/*
Copyright 2014 The go-marathon Authors All rights reserved.

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

func TestGroups(t *testing.T) {
	endpoint := newFakeMarathonEndpoint(t, nil)
	defer endpoint.Close()

	groups, err := endpoint.Client.Groups()
	assert.NoError(t, err)
	assert.NotNil(t, groups)
	assert.Equal(t, 1, len(groups.Groups))
	group := groups.Groups[0]
	assert.Equal(t, fakeGroupName, group.ID)
}

func TestNewGroup(t *testing.T) {
	endpoint := newFakeMarathonEndpoint(t, nil)
	defer endpoint.Close()

	group, err := endpoint.Client.Group(fakeGroupName)
	assert.NoError(t, err)
	assert.NotNil(t, group)
	assert.Equal(t, 1, len(group.Apps))
	assert.Equal(t, fakeGroupName, group.ID)

	group, err = endpoint.Client.Group(fakeGroupName1)

	assert.NoError(t, err)
	assert.NotNil(t, group)
	assert.Equal(t, fakeGroupName1, group.ID)
	assert.NotNil(t, group.Groups)
	assert.Equal(t, 1, len(group.Groups))

	frontend := group.Groups[0]
	assert.Equal(t, "frontend", frontend.ID)
	assert.Equal(t, 3, len(frontend.Apps))
	for _, app := range frontend.Apps {
		assert.NotNil(t, app.Container)
		assert.NotNil(t, app.Container.Docker)
		for _, network := range *app.Networks {
			assert.Equal(t, BridgeNetworkMode, network.Mode)
		}
		if len(*app.Container.PortMappings) == 0 {
			t.Fail()
		}
	}
}

// TODO @kamsz: How to work with old and new endpoints from methods.yml?
// func TestGroup(t *testing.T) {
// 	endpoint := newFakeMarathonEndpoint(t, nil)
// 	defer endpoint.Close()

// 	group, err := endpoint.Client.Group(fakeGroupName)
// 	assert.NoError(t, err)
// 	assert.NotNil(t, group)
// 	assert.Equal(t, 1, len(group.Apps))
// 	assert.Equal(t, fakeGroupName, group.ID)

// 	group, err = endpoint.Client.Group(fakeGroupName1)

// 	assert.NoError(t, err)
// 	assert.NotNil(t, group)
// 	assert.Equal(t, fakeGroupName1, group.ID)
// 	assert.NotNil(t, group.Groups)
// 	assert.Equal(t, 1, len(group.Groups))

// 	frontend := group.Groups[0]
// 	assert.Equal(t, "frontend", frontend.ID)
// 	assert.Equal(t, 3, len(frontend.Apps))
// 	for _, app := range frontend.Apps {
// 		assert.NotNil(t, app.Container)
// 		assert.NotNil(t, app.Container.Docker)
// 		assert.Equal(t, "BRIDGE", app.Container.Docker.Network)
// 		if len(*app.Container.Docker.PortMappings) == 0 {
// 			t.Fail()
// 		}
// 	}
// }
