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

	"github.com/stretchr/testify/assert"
)

type testCase struct {
	source      string
	expectation interface{}
}

var testCases = map[string]*testCase{
	"status_update_event": &testCase{
		`{
			"eventType": "status_update_event",
			"timestamp": "2014-03-01T23:29:30.158Z",
			"slaveId": "20140909-054127-177048842-5050-1494-0",
			"taskId": "my-app_0-1396592784349",
			"taskStatus": "TASK_RUNNING",
			"appId": "/my-app",
			"host": "slave-1234.acme.org",
			"ports": [31372],
			"version": "2014-04-04T06:26:23.051Z"
		}`,
		&EventStatusUpdate{
			EventType:  "status_update_event",
			Timestamp:  "2014-03-01T23:29:30.158Z",
			SlaveID:    "20140909-054127-177048842-5050-1494-0",
			TaskID:     "my-app_0-1396592784349",
			TaskStatus: "TASK_RUNNING",
			AppID:      "/my-app",
			Host:       "slave-1234.acme.org",
			Ports:      []int{31372},
			Version:    "2014-04-04T06:26:23.051Z",
		},
	},
	"health_status_changed_event": &testCase{
		`{
			"eventType": "health_status_changed_event",
			"timestamp": "2014-03-01T23:29:30.158Z",
			"appId": "/my-app",
			"taskId": "my-app_0-1396592784349",
			"version": "2014-04-04T06:26:23.051Z",
			"alive": true
		}`,
		&EventHealthCheckChanged{
			EventType: "health_status_changed_event",
			Timestamp: "2014-03-01T23:29:30.158Z",
			AppID:     "/my-app",
			TaskID:    "my-app_0-1396592784349",
			Version:   "2014-04-04T06:26:23.051Z",
			Alive:     true,
		},
	},
	"failed_health_check_event": &testCase{
		`{
			"eventType": "failed_health_check_event",
			"timestamp": "2014-03-01T23:29:30.158Z",
			"appId": "/my-app",
			"taskId": "my-app_0-1396592784349",
			"healthCheck": {
				"protocol": "HTTP",
				"path": "/health",
				"portIndex": 0,
				"gracePeriodSeconds": 5,
				"intervalSeconds": 10,
				"timeoutSeconds": 10,
				"maxConsecutiveFailures": 3
			}
		}`,
		&EventFailedHealthCheck{
			EventType: "failed_health_check_event",
			Timestamp: "2014-03-01T23:29:30.158Z",
			AppId:     "/my-app",
			HealthCheck: struct {
				GracePeriodSeconds     float64 `json:"gracePeriodSeconds"`
				IntervalSeconds        float64 `json:"intervalSeconds"`
				MaxConsecutiveFailures float64 `json:"maxConsecutiveFailures"`
				Path                   string  `json:"path"`
				PortIndex              float64 `json:"portIndex"`
				Protocol               string  `json:"protocol"`
				TimeoutSeconds         float64 `json:"timeoutSeconds"`
			}{
				GracePeriodSeconds:     5,
				IntervalSeconds:        10,
				MaxConsecutiveFailures: 3,
				Path:           "/health",
				PortIndex:      0,
				Protocol:       "HTTP",
				TimeoutSeconds: 10,
			},
		},
	},
}

func TestSubscriptions(t *testing.T) {
	endpoint := newFakeMarathonEndpoint(t, nil)
	defer endpoint.Close()

	sub, err := endpoint.Client.Subscriptions()
	assert.NoError(t, err)
	assert.NotNil(t, sub)
	assert.NotNil(t, sub.CallbackURLs)
	assert.Equal(t, len(sub.CallbackURLs), 1)
}

func TestEventStreamConnectionErrorsForwarded(t *testing.T) {
	config := NewDefaultConfig()
	config.EventsTransport = EventsTransportSSE
	config.URL = "http://non-existing-marathon-host.local:5555"
	endpoint := newFakeMarathonEndpoint(t, &config)
	defer endpoint.Close()

	events := make(EventsChannel)
	err := endpoint.Client.AddEventsListener(events, EVENTS_APPLICATIONS)
	assert.Error(t, err)
}

func TestEventStreamEventsReceived(t *testing.T) {
	config := NewDefaultConfig()
	config.EventsTransport = EventsTransportSSE
	endpoint := newFakeMarathonEndpoint(t, &config)
	defer endpoint.Close()

	events := make(EventsChannel)
	err := endpoint.Client.AddEventsListener(events, EVENTS_APPLICATIONS)
	assert.NoError(t, err)

	// Publish test events
	go func() {
		for _, testCase := range testCases {
			endpoint.Server.PublishEvent(testCase.source)
		}
	}()

	// Receive test events
	for i := 0; i < len(testCases); i++ {
		event := <-events
		assert.Equal(t, testCases[event.Name].expectation, event.Event)
	}
}
