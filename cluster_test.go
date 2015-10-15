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
	"time"

	"github.com/stretchr/testify/assert"
)

func TestUrl(t *testing.T) {
	cluster, _ := newCluster(fakeMarathonURL)
	assert.Equal(t, cluster.URL(), fakeMarathonURL)
}

func TestSize(t *testing.T) {
	cluster, _ := newCluster(fakeMarathonURL)
	assert.Equal(t, cluster.Size(), 3)
}

func TestActive(t *testing.T) {
	cluster, _ := newCluster(fakeMarathonURL)
	assert.Equal(t, len(cluster.Active()), 3)
}

func TestNonActive(t *testing.T) {
	cluster, _ := newCluster(fakeMarathonURL)
	assert.Equal(t, len(cluster.NonActive()), 0)
}

func TestGetMember(t *testing.T) {
	cluster, _ := newCluster(fakeMarathonURL)
	member, err := cluster.GetMember()
	assert.Nil(t, err)
	assert.Equal(t, member, "http://127.0.0.1:3000")
}

func TestMarkdown(t *testing.T) {
	cluster, _ := newCluster(fakeMarathonURL)
	assert.Equal(t, len(cluster.Active()), 3)
	cluster.MarkDown()
	cluster.MarkDown()
	assert.Equal(t, len(cluster.Active()), 1)
	time.Sleep(10 * time.Millisecond)
	assert.Equal(t, len(cluster.Active()), 3)
}
