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

import "fmt"

// EventType is a wrapper for a marathon event
type EventType struct {
	EventType string `json:"eventType"`
}

const (
	EVENT_API_REQUEST = 1 << iota
	EVENT_STATUS_UPDATE
	EVENT_FRAMEWORK_MESSAGE
	EVENT_SUBSCRIPTION
	EVENT_UNSUBSCRIBED
	EVENT_STREAM_ATTACHED
	EVENT_STREAM_DETACHED
	EVENT_ADD_HEALTH_CHECK
	EVENT_REMOVE_HEALTH_CHECK
	EVENT_FAILED_HEALTH_CHECK
	EVENT_CHANGED_HEALTH_CHECK
	EVENT_GROUP_CHANGE_SUCCESS
	EVENT_GROUP_CHANGE_FAILED
	EVENT_DEPLOYMENT_SUCCESS
	EVENT_DEPLOYMENT_FAILED
	EVENT_DEPLOYMENT_INFO
	EVENT_DEPLOYMENT_STEP_SUCCESS
	EVENT_DEPLOYMENT_STEP_FAILED
	EVENT_APP_TERMINATED
)

const (
	EVENTS_APPLICATIONS  = EVENT_STATUS_UPDATE | EVENT_CHANGED_HEALTH_CHECK | EVENT_FAILED_HEALTH_CHECK | EVENT_APP_TERMINATED
	EVENTS_SUBSCRIPTIONS = EVENT_SUBSCRIPTION | EVENT_UNSUBSCRIBED | EVENT_STREAM_ATTACHED | EVENT_STREAM_DETACHED
)

var (
	Events map[string]int
)

func init() {
	Events = map[string]int{
		"api_post_event":              EVENT_API_REQUEST,
		"status_update_event":         EVENT_STATUS_UPDATE,
		"framework_message_event":     EVENT_FRAMEWORK_MESSAGE,
		"subscribe_event":             EVENT_SUBSCRIPTION,
		"unsubscribe_event":           EVENT_UNSUBSCRIBED,
		"event_stream_attached":       EVENT_STREAM_ATTACHED,
		"event_stream_detached":       EVENT_STREAM_DETACHED,
		"add_health_check_event":      EVENT_ADD_HEALTH_CHECK,
		"remove_health_check_event":   EVENT_REMOVE_HEALTH_CHECK,
		"failed_health_check_event":   EVENT_FAILED_HEALTH_CHECK,
		"health_status_changed_event": EVENT_CHANGED_HEALTH_CHECK,
		"group_change_success":        EVENT_GROUP_CHANGE_SUCCESS,
		"group_change_failed":         EVENT_GROUP_CHANGE_FAILED,
		"deployment_success":          EVENT_DEPLOYMENT_SUCCESS,
		"deployment_failed":           EVENT_DEPLOYMENT_FAILED,
		"deployment_info":             EVENT_DEPLOYMENT_INFO,
		"deployment_step_success":     EVENT_DEPLOYMENT_STEP_SUCCESS,
		"deployment_step_failure":     EVENT_DEPLOYMENT_STEP_FAILED,
		"app_terminated_event":        EVENT_APP_TERMINATED,
	}
}

//
//  Events taken from: https://mesosphere.github.io/marathon/docs/event-bus.html
//

// Event is the definition for a event in marathon
type Event struct {
	ID    int
	Name  string
	Event interface{}
}

func (event *Event) String() string {
	return fmt.Sprintf("type: %s, event: %s", event.Name, event.Event)
}

// EventsChannel is a channel to receive events upon
type EventsChannel chan *Event

/* --- API Request --- */

type EventAPIRequest struct {
	EventType     string       `json:"eventType"`
	ClientIp      string       `json:"clientIp"`
	Timestamp     string       `json:"timestamp"`
	Uri           string       `json:"uri"`
	AppDefinition *Application `json:"appDefinition"`
}

/* --- Status Update --- */

type EventStatusUpdate struct {
	EventType  string `json:"eventType"`
	Timestamp  string `json:"timestamp,omitempty"`
	SlaveID    string `json:"slaveId,omitempty"`
	TaskID     string `json:"taskId"`
	TaskStatus string `json:"taskStatus"`
	AppID      string `json:"appId"`
	Host       string `json:"host"`
	Ports      []int  `json:"ports,omitempty"`
	Version    string `json:"version,omitempty"`
}

type EventAppTerminated struct {
	EventType string `json:"eventType"`
	Timestamp string `json:"timestamp,omitempty"`
	AppID     string `json:"appId"`
}

/* --- Framework Message --- */

type EventFrameworkMessage struct {
	EventType  string `json:"eventType"`
	ExecutorId string `json:"executorId"`
	Message    string `json:"message"`
	SlaveId    string `json:"slaveId"`
	Timestamp  string `json:"timestamp"`
}

/* --- Event Subscription --- */

// EventSubscription describes 'subscribe_event' Marathon event
type EventSubscription struct {
	CallbackUrl string `json:"callbackUrl"`
	ClientIp    string `json:"clientIp"`
	EventType   string `json:"eventType"`
	Timestamp   string `json:"timestamp"`
}

// EventUnsubscription describes 'unsubscribe_event' Marathon event
type EventUnsubscription struct {
	CallbackUrl string `json:"callbackUrl"`
	ClientIp    string `json:"clientIp"`
	EventType   string `json:"eventType"`
	Timestamp   string `json:"timestamp"`
}

// EventStreamAttached describes 'event_stream_attached' Marathon event
type EventStreamAttached struct {
	RemoteAddress string `json:"remoteAddress"`
	EventType     string `json:"eventType"`
	Timestamp     string `json:"timestamp"`
}

// EventStreamDetached describes 'event_stream_detached' Marathon event
type EventStreamDetached struct {
	RemoteAddress string `json:"remoteAddress"`
	EventType     string `json:"eventType"`
	Timestamp     string `json:"timestamp"`
}

/* --- Health Checks --- */

type EventAddHealthCheck struct {
	AppId       string `json:"appId"`
	EventType   string `json:"eventType"`
	HealthCheck struct {
		GracePeriodSeconds     float64 `json:"gracePeriodSeconds"`
		IntervalSeconds        float64 `json:"intervalSeconds"`
		MaxConsecutiveFailures float64 `json:"maxConsecutiveFailures"`
		Path                   string  `json:"path"`
		PortIndex              float64 `json:"portIndex"`
		Protocol               string  `json:"protocol"`
		TimeoutSeconds         float64 `json:"timeoutSeconds"`
	} `json:"healthCheck"`
	Timestamp string `json:"timestamp"`
}

type EventRemoveHealthCheck struct {
	AppId       string `json:"appId"`
	EventType   string `json:"eventType"`
	HealthCheck struct {
		GracePeriodSeconds     float64 `json:"gracePeriodSeconds"`
		IntervalSeconds        float64 `json:"intervalSeconds"`
		MaxConsecutiveFailures float64 `json:"maxConsecutiveFailures"`
		Path                   string  `json:"path"`
		PortIndex              float64 `json:"portIndex"`
		Protocol               string  `json:"protocol"`
		TimeoutSeconds         float64 `json:"timeoutSeconds"`
	} `json:"healthCheck"`
	Timestamp string `json:"timestamp"`
}

type EventFailedHealthCheck struct {
	AppId       string `json:"appId"`
	EventType   string `json:"eventType"`
	HealthCheck struct {
		GracePeriodSeconds     float64 `json:"gracePeriodSeconds"`
		IntervalSeconds        float64 `json:"intervalSeconds"`
		MaxConsecutiveFailures float64 `json:"maxConsecutiveFailures"`
		Path                   string  `json:"path"`
		PortIndex              float64 `json:"portIndex"`
		Protocol               string  `json:"protocol"`
		TimeoutSeconds         float64 `json:"timeoutSeconds"`
	} `json:"healthCheck"`
	Timestamp string `json:"timestamp"`
}

type EventHealthCheckChanged struct {
	EventType string `json:"eventType"`
	Timestamp string `json:"timestamp,omitempty"`
	AppID     string `json:"appId"`
	TaskID    string `json:"taskId"`
	Version   string `json:"version,omitempty"`
	Alive     bool   `json:"alive"`
}

/* --- Deployments --- */

type EventGroupChangeSuccess struct {
	EventType string `json:"eventType"`
	GroupId   string `json:"groupId"`
	Timestamp string `json:"timestamp"`
	Version   string `json:"version"`
}

type EventGroupChangeFailed struct {
	EventType string `json:"eventType"`
	GroupId   string `json:"groupId"`
	Timestamp string `json:"timestamp"`
	Version   string `json:"version"`
	Reason    string `json:"reason"`
}

type EventDeploymentSuccess struct {
	ID        string `json:"id"`
	EventType string `json:"eventType"`
	Timestamp string `json:"timestamp"`
}

type EventDeploymentFailed struct {
	ID        string `json:"id"`
	EventType string `json:"eventType"`
	Timestamp string `json:"timestamp"`
}

type EventDeploymentInfo struct {
	EventType   string          `json:"eventType"`
	CurrentStep *DeploymentStep `json:"currentStep"`
	Timestamp   string          `json:"timestamp"`
	Plan        *DeploymentPlan `json:"plan"`
}

type EventDeploymentStepSuccess struct {
	EventType   string          `json:"eventType"`
	CurrentStep *DeploymentStep `json:"currentStep"`
	Timestamp   string          `json:"timestamp"`
	Plan        *DeploymentPlan `json:"plan"`
}

type EventDeploymentStepFailure struct {
	EventType   string          `json:"eventType"`
	CurrentStep *DeploymentStep `json:"currentStep"`
	Timestamp   string          `json:"timestamp"`
	Plan        *DeploymentPlan `json:"plan"`
}

// GetEvent returns allocated empty event object which corresponds to provided event type
//		eventType:			the type of Marathon event
func GetEvent(eventType string) (*Event, error) {
	// step: check it's supported
	id, found := Events[eventType]
	if found {
		event := new(Event)
		event.ID = id
		event.Name = eventType
		switch eventType {
		case "api_post_event":
			event.Event = new(EventAPIRequest)
		case "status_update_event":
			event.Event = new(EventStatusUpdate)
		case "framework_message_event":
			event.Event = new(EventFrameworkMessage)
		case "subscribe_event":
			event.Event = new(EventSubscription)
		case "unsubscribe_event":
			event.Event = new(EventUnsubscription)
		case "event_stream_attached":
			event.Event = new(EventStreamAttached)
		case "event_stream_detached":
			event.Event = new(EventStreamDetached)
		case "add_health_check_event":
			event.Event = new(EventAddHealthCheck)
		case "remove_health_check_event":
			event.Event = new(EventRemoveHealthCheck)
		case "failed_health_check_event":
			event.Event = new(EventFailedHealthCheck)
		case "health_status_changed_event":
			event.Event = new(EventHealthCheckChanged)
		case "group_change_success":
			event.Event = new(EventGroupChangeSuccess)
		case "group_change_failed":
			event.Event = new(EventGroupChangeFailed)
		case "deployment_success":
			event.Event = new(EventDeploymentSuccess)
		case "deployment_failed":
			event.Event = new(EventDeploymentFailed)
		case "deployment_info":
			event.Event = new(EventDeploymentInfo)
		case "deployment_step_success":
			event.Event = new(EventDeploymentStepSuccess)
		case "deployment_step_failure":
			event.Event = new(EventDeploymentStepFailure)
		case "app_terminated_event":
			event.Event = new(EventAppTerminated)
		}
		return event, nil
	}

	return nil, fmt.Errorf("the event type: %s was not found or supported", eventType)
}
