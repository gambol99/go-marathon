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
	"errors"
	"fmt"
	"io/ioutil"
	"net"
	"net/http"
	"strings"
	"time"
)

type Subscriptions struct {
	CallbackURLs []string `json:"callbackUrls"`
}

// Retrieve a list of registered subscriptions
func (client *Client) Subscriptions() (*Subscriptions, error) {
	subscriptions := new(Subscriptions)
	if err := client.ApiGet(MARATHON_API_SUBSCRIPTION, "", subscriptions); err != nil {
		return nil, err
	} else {
		return subscriptions, nil
	}
}

func (client *Client) AddEventsListener(channel EventsChannel, filter int) error {
	client.Lock()
	defer client.Unlock()
	/* step: someone has asked to start listening to event, we need to register for events
	if we haven't done so already
	*/
	if err := client.RegisterSubscription(); err != nil {
		return err
	}

	if _, found := client.listeners[channel]; !found {
		client.listeners[channel] = filter
	}
	return nil
}

func (client *Client) RemoveEventsListener(channel EventsChannel) {
	client.Lock()
	defer client.Unlock()
	if _, found := client.listeners[channel]; found {
		delete(client.listeners, channel)
		/* step: if there is no one listening anymore, lets remove our self
		from the events callback */
		if len(client.listeners) <= 0 {
			client.UnSubscribe()
		}
	}
}

func (client *Client) SubscriptionURL() string {
	return fmt.Sprintf("http://%s:%d%s", client.ipaddress, client.config.EventsPort, DEFAULT_EVENTS_URL)
}

func (client *Client) RegisterSubscription() error {
	if client.events_http == nil {
		if ip_address, err := GetInterfaceAddress(client.config.EventsInterface); err != nil {
			return errors.New(fmt.Sprintf("Unable to get the ip address from the interface: %s, error: %s",
				client.config.EventsInterface, err))
		} else {
			/* step: set the ip address */
			client.ipaddress = ip_address
			binding := fmt.Sprintf("%s:%d", ip_address, client.config.EventsPort)
			/* step: register the handler */
			http.HandleFunc(DEFAULT_EVENTS_URL, client.HandleMarathonEvent)
			/* step: create the http server */
			client.Debug("Creating a event http server")
			client.events_http = &http.Server{
				Addr:           binding,
				Handler:        nil,
				ReadTimeout:    10 * time.Second,
				WriteTimeout:   10 * time.Second,
				MaxHeaderBytes: 1 << 20,
			}
			/* step: try and listen on the port */
			if listener, err := net.Listen("tcp", binding); err != nil {
				return nil
			} else {
				go func() {
					for {
						client.Debug("Starting to listen on http events service")
						// step: not sure how to handle an error from this, panic???
						client.events_http.Serve(listener)
						client.Debug("Exitted the http events service")
					}
				}()
			}
		}
	}

	/* step: get the callback url */
	callback := client.SubscriptionURL()
	/* step: check if the callback is registered */
	if found, err := client.HasSubscription(callback); err != nil {
		return err
	} else if !found {
		client.Debug("Registering a subscription with Marathon: callback: %s", callback)
		/* step: we need to register our self */
		uri := fmt.Sprintf("%s?callbackUrl=%s", MARATHON_API_SUBSCRIPTION, callback)
		if err := client.ApiPost(uri, "", nil); err != nil {
			return err
		}
	} else {
		client.Debug("A subscription already exists for this callback: %s", callback)
	}
	return nil
}

func (client *Client) UnSubscribe() error {
	/* step: remove from the list of subscriptions */
	return client.ApiDelete(fmt.Sprintf("%s?callbackUrl=%s", MARATHON_API_SUBSCRIPTION, client.SubscriptionURL()), "", nil)
}

func (client *Client) HasSubscription(callback string) (bool, error) {
	client.Debug("Checking for subscription: %s", callback)
	/* step: generate our events callback */
	if subscriptions, err := client.Subscriptions(); err != nil {
		return false, err
	} else {
		for _, subscription := range subscriptions.CallbackURLs {
			if callback == subscription {
				return true, nil
			}
		}
	}
	return false, nil
}

func (client *Client) HandleMarathonEvent(writer http.ResponseWriter, request *http.Request) {
	client.Debug("Recieved a possible marathon event")
	/* step: lets read in the post body */
	if body, err := ioutil.ReadAll(request.Body); err == nil {
		content := string(body[:])
		/* step: phase one, get the event type */
		decoder := json.NewDecoder(strings.NewReader(content))
		/* step: decode the event type */
		event_type := new(EventType)
		if err := decoder.Decode(event_type); err != nil {
			client.Debug("Failed to decode the event type, content: %s, error: %s", content, err)
			return
		}

		client.Debug("Recieved marathon event, %s", event_type.EventType)
		/* step: check the type is handled */
		if event, err := client.GetEvent(event_type.EventType); err != nil {
			client.Debug("Unable to retrieve the event, type: %s", event_type.EventType)
		} else {
			/* step: lets decode */
			decoder = json.NewDecoder(strings.NewReader(content))
			if err := decoder.Decode(event.Event); err != nil {
				client.Debug("Failed to decode the event type: %d, name: %s error: %s", event.Type, err)
			}
			client.Debug("Marathon event, %s", event)
			client.RLock()
			defer client.RUnlock()
			/* step: check if anyone is listen for this event */
			for channel, filter := range client.listeners {
				/* step: check if this person wants this event type */
				if event.Type&filter != 0 {
					client.Debug("Event type: %d being listen to, sending to lister: %s", event.Type, channel)
					go func() {
						channel <- event
					}()
				}
			}
		}
	} else {
		client.Debug("Failed to decode the event type, content: %s, error: %s")
	}
}

func (client *Client) GetEvent(name string) (*Event, error) {
	/* step: check it's supported */
	if event_id, found := Events[name]; found {
		event := new(Event)
		event.Type = event_id
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
		}
		return event, nil
	} else {
		return nil, errors.New(fmt.Sprintf("The event type: %d was not found or supported", name))
	}
}
