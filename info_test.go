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
)

func TestInfo(t *testing.T) {
	client := NewFakeMarathonEndpoint(t)
	info, err := client.Info()
	assertOnError(err, t)
	assertOnString(info.FrameworkId, "20140730-222531-1863654316-5050-10422-0000", t)
	assertOnString(info.Leader, "127.0.0.1:8080", t)
	assertOnString(info.Version, "0.7.0-SNAPSHOT", t)
}

func TestLeader(t *testing.T) {
	client := NewFakeMarathonEndpoint(t)
	leader, err := client.Leader()
	assertOnError(err, t)
	assertOnString(leader, "127.0.0.1:8080", t)
}

func TestAbdicateLeader(t *testing.T) {
	client := NewFakeMarathonEndpoint(t)
	message, err := client.AbdicateLeader()
	assertOnError(err, t)
	assertOnString(message, "Leadership abdicted", t)
}
