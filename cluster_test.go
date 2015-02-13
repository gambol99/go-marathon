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
)


var cluster Cluster

func GetFakeCluster() {
	if cluster == nil {
		cluster, _ = NewMarathonCluster(FAKE_MARATHON_URL)
	}
}

func TestUrl(t *testing.T) {
	GetFakeCluster()
	AssertOnString(cluster.Url(), FAKE_MARATHON_URL, t)
}

func TestSize(t *testing.T) {
	GetFakeCluster()
	AssertOnInteger(cluster.Size(), 2, t)
}

func TestActive(t *testing.T) {
	GetFakeCluster()
	AssertOnInteger(len(cluster.Active()), 2, t)
}

func TestNonActive(t *testing.T) {
	GetFakeCluster()
	AssertOnInteger(len(cluster.NonActive()), 0, t)
}

func TestGetMember(t *testing.T) {
	GetFakeCluster()
	member, err := cluster.GetMember()
	AssertOnError(err, t)
	AssertOnString(member, "http://127.0.0.1:3000", t)
}

func TestMarkdown(t *testing.T) {
	GetFakeCluster()
	cluster.MarkDown()
	AssertOnInteger(len(cluster.Active()), 2, t)
	time.Sleep(200 * time.Millisecond)
	AssertOnInteger(len(cluster.Active()), 2, t)
}


