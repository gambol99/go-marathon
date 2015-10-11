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
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net"
	"net/http"
	"strings"
	"time"

	"github.com/golang/glog"
)

// Subscriptions is a collection to urls that marathon is implmenting a callback on
type Subscriptions struct {
	CallbackURLs []string `json:"callbackUrls"`
}

// Subscriptions retrieves a list of registered subscriptions
func (r *marathonClient) Subscriptions() (*Subscriptions, error) {
	subscriptions := new(Subscriptions)
	if err := r.apiGet(MARATHON_API_SUBSCRIPTION, nil, subscriptions); err != nil {
		return nil, err
	}

	return subscriptions, nil
}

// AddEventsListener adds your self as a listener to events from Marathon
//		channel:	a EventsChannel used to receive event on
func (r *marathonClient) AddEventsListener(channel EventsChannel, filter int) error {
	r.Lock()
	defer r.Unlock()
	// step: someone has asked to start listening to event, we need to register for events
	// if we haven't done so already
	if err := r.RegisterSubscription(); err != nil {
		return err
	}

	if _, found := r.listeners[channel]; !found {
		glog.V(DEBUG_LEVEL).Infof("adding a watch for events: %d, channel: %v", filter, channel)
		r.listeners[channel] = filter
	}
	return nil
}

// RemoveEventsListener removes the channel from the events listeners
//		channel:			the channel you are removing
func (r *marathonClient) RemoveEventsListener(channel EventsChannel) {
	r.Lock()
	defer r.Unlock()
	if _, found := r.listeners[channel]; found {
		delete(r.listeners, channel)
		/* step: if there is no one listening anymore, lets remove our self
		from the events callback */
		if len(r.listeners) <= 0 {
			r.UnSubscribe()
		}
	}
}

// SubscriptionURL retrieves the subscription call back URL used when registering
func (r *marathonClient) SubscriptionURL() string {
	return fmt.Sprintf("http://%s:%d%s", r.ipAddress, r.config.EventsPort, DEFAULT_EVENTS_URL)
}

// RegisterSubscription registers ourselves with Marathon to receive events from it's callback facility
func (r *marathonClient) RegisterSubscription() error {
	if r.eventsHTTP == nil {
		ipAddress, err := getInterfaceAddress(r.config.EventsInterface)
		if err != nil {
			return fmt.Errorf("Unable to get the ip address from the interface: %s, error: %s",
				r.config.EventsInterface, err)
		}

		// step: set the ip address
		r.ipAddress = ipAddress
		binding := fmt.Sprintf("%s:%d", ipAddress, r.config.EventsPort)
		// step: register the handler
		http.HandleFunc(DEFAULT_EVENTS_URL, r.handleMarathonEvent)
		// step: create the http server
		r.eventsHTTP = &http.Server{
			Addr:           binding,
			Handler:        nil,
			ReadTimeout:    10 * time.Second,
			WriteTimeout:   10 * time.Second,
			MaxHeaderBytes: 1 << 20,
		}

		// @todo need to add a timeout value here
		listener, err := net.Listen("tcp", binding)
		if err != nil {
			return nil
		}

		go func() {
			for {
				r.eventsHTTP.Serve(listener)
			}
		}()
	}

	// step: get the callback url
	callback := r.SubscriptionURL()

	// step: check if the callback is registered
	found, err := r.HasSubscription(callback)
	if err != nil {
		return err
	}
	if !found {
		// step: we need to register our self
		uri := fmt.Sprintf("%s?callbackUrl=%s", MARATHON_API_SUBSCRIPTION, callback)
		if err := r.apiPost(uri, "", nil); err != nil {
			return err
		}
	}

	return nil
}

// UnSubscribe removes ourselves from Marathon's callback facility
func (r *marathonClient) UnSubscribe() error {
	// step: remove from the list of subscriptions
	return r.apiDelete(fmt.Sprintf("%s?callbackUrl=%s", MARATHON_API_SUBSCRIPTION, r.SubscriptionURL()), nil, nil)
}

// HasSubscription checks to see a subscription already exists with Marathon
//		callback:			the url of the callback
func (r *marathonClient) HasSubscription(callback string) (bool, error) {
	// step: generate our events callback
	subscriptions, err := r.Subscriptions()
	if err != nil {
		return false, err
	}

	for _, subscription := range subscriptions.CallbackURLs {
		if callback == subscription {
			return true, nil
		}
	}

	return false, nil
}

func (r *marathonClient) handleMarathonEvent(writer http.ResponseWriter, request *http.Request) {
	body, err := ioutil.ReadAll(request.Body)
	if err != nil {
		return
	}

	// step: process the event and decode the event
	content := string(body[:])
	eventType := new(EventType)
	err = json.NewDecoder(strings.NewReader(content)).Decode(eventType)
	if err != nil {
		glog.V(DEBUG_LEVEL).Infof("failed to decode the event type, content: %s, error: %s", content, err)
		return
	}

	// step: check the type is handled
	event, err := r.GetEvent(eventType.EventType)
	if err != nil {
		glog.V(DEBUG_LEVEL).Infof("unable to retrieve the event, type: %s", eventType.EventType)
		return
	}

	// step: lets decode message
	err = json.NewDecoder(strings.NewReader(content)).Decode(event.Event)
	if err != nil {
		glog.V(DEBUG_LEVEL).Infof("failed to decode the event type: %d, name: %s error: %s", event.ID, err)
		return
	}

	r.RLock()
	defer r.RUnlock()

	// step: check if anyone is listen for this event
	for channel, filter := range r.listeners {
		// step: check if this listener wants this event type
		if event.ID&filter != 0 {
			go func(ch EventsChannel, e *Event) {
				ch <- e
			}(channel, event)
		}
	}
}

func (r *marathonClient) GetEvent(name string) (*Event, error) {
	// step: check it's supported
	id, found := Events[name]
	if found {
		event := new(Event)
		event.ID = id
		event.Name = name
		switch name {
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

	return nil, fmt.Errorf("the event type: %s was not found or supported", name)
}
