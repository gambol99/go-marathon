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

func TestAllTasks(t *testing.T) {
	client := NewFakeMarathonEndpoint(t)
	tasks, err := client.AllTasks()
	AssertOnError(err, t)
	AssertOnNull(tasks, t)
	AssertOnInteger(len(tasks.Tasks), 2, t)
}

func TestTaskEndpoints(t *testing.T) {
	client := NewFakeMarathonEndpoint(t)
	endpoints, err := client.TaskEndpoints(FAKE_APP_NAME_BROKEN, 80, true)
	AssertOnNoError(err, t)
	endpoints, err = client.TaskEndpoints(FAKE_APP_NAME_BROKEN, 8080, true)
	AssertOnError(err, t)
	AssertOnNull(endpoints, t)
	AssertOnInteger(len(endpoints), 1, t)
	AssertOnString(endpoints[0], "10.141.141.10:31045", t)
	endpoints, err = client.TaskEndpoints(FAKE_APP_NAME_BROKEN, 8080, false)
	AssertOnError(err, t)
	AssertOnNull(endpoints, t)
	AssertOnInteger(len(endpoints), 2, t)
	AssertOnString(endpoints[0], "10.141.141.10:31045", t)
	AssertOnString(endpoints[1], "10.141.141.10:31234", t)
}

func TestKillApplicationTasks(t *testing.T) {
	client := NewFakeMarathonEndpoint(t)
	tasks, err := client.KillApplicationTasks(FAKE_APP_NAME, "", false)
	AssertOnError(err, t)
	AssertOnNull(tasks, t)
}
