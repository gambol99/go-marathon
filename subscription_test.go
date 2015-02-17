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

func TestSubscriptions(t *testing.T) {
	NewFakeMarathonEndpoint()
	sub, err := test_client.Subscriptions()
	AssertOnError(err, t)
	AssertOnNull(sub, t)
	AssertOnNull(sub.CallbackURLs, t)
	AssertOnInteger(len(sub.CallbackURLs), 1, t)
}

func TestWatch(t *testing.T) {
	NewFakeMarathonEndpoint()
	channel := make(chan string)
	AssertOnNull(test_client.WatchList(), t)
	AssertOnInteger(len(test_client.WatchList()), 0, t)
	test_client.Watch(FAKE_APP_NAME, channel)
	AssertOnNull(test_client.WatchList(), t)
	AssertOnInteger(len(test_client.WatchList()), 1, t)
}

func TestRemove(t *testing.T) {
	NewFakeMarathonEndpoint()
	AssertOnNull(test_client.WatchList(), t)
	channel := make(chan string)
	test_client.Watch(FAKE_APP_NAME, channel)
	AssertOnNull(test_client.WatchList(), t)
	AssertOnInteger(len(test_client.WatchList()), 1, t)
	test_client.RemoveWatch(FAKE_APP_NAME)
	AssertOnNull(test_client.WatchList(), t)
	AssertOnInteger(len(test_client.WatchList()), 0, t)
}
