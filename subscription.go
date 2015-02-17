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
	"net/http"
	"strings"
)

type Subscriptions struct {
	CallbackURLs []string `json:"callbackUrls"`
}

func (client *Client) Subscriptions() (*Subscriptions, error) {
	subscriptions := new(Subscriptions)
	if err := client.ApiGet(MARATHON_API_SUBSCRIPTION, "", subscriptions); err != nil {
		return nil, err
	} else {
		return subscriptions, nil
	}
}

func (client *Client) RegisterSubscription() error {
	/* step: lets lock the client */
	client.Lock()
	defer client.Unlock()
	if !client.events_running {
		go func() {
			/* step: generate the call url */
			if _, err := client.SubscriptionURL(); err != nil {
				return
			}
			/* step: register the handler */
			http.HandleFunc(DEFAULT_EVENTS_URL, client.HandleMarathonEvent)
			/* step: register and listen */
			http.ListenAndServe(fmt.Sprintf("%s:%d", client.events_ipaddress, client.config.EventsPort), nil)
			/* step; unset the boolean */
			client.events_running = false
		}()
	}
	/* step: check if we are already subscribed */
	if found, err := client.HasSubscription(); err != nil {
		return err
	} else if found {
		return nil
	}
	/* step: register the event callback */
	uri := fmt.Sprintf("%s?callbackUrl=%s", MARATHON_API_SUBSCRIPTION, client.events_callback_url)
	if err := client.ApiPost(uri, "", nil); err != nil {
		return err
	}
	return nil
}

func (client *Client) DeregisterSubscription() error {
	/* step: check if we are already subscribed */
	found, err := client.HasSubscription()
	if err != nil {
		return err
	} else if found {
		/* step: remove from the list of subscriptions */
		uri := fmt.Sprintf("%s?callbackUrl=%s", MARATHON_API_SUBSCRIPTION, client.events_callback_url)
		if err := client.ApiDelete(uri, "", nil); err != nil {
			return err
		}
	}
	return nil
}

func (client *Client) HasSubscription() (bool, error) {
	/* step: generate our events callback */
	if callback, err := client.SubscriptionURL(); err != nil {
		return false, err
	} else {
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
}

func (client *Client) HandleMarathonEvent(writer http.ResponseWriter, request *http.Request) {
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
		/* step: check the type is handled */
		if event_type_value, found := Events[event_type.EventType]; found {
			event := client.GetEvent(event_type.EventType)
			event.EventType = event_type_value
			/* step: lets decode */
			decoder = json.NewDecoder(strings.NewReader(content))
			if err := decoder.Decode(event.Event); err != nil {
				client.Debug("Failed to decode the event type: %s, error: %s", event_type, err)
			}

		} else {
			client.Debug("The event type: %s was not found", event_type.EventType)
		}
	} else {
		client.Debug("Failed to decode the event type, content: %s, error: %s")
	}
}

func (client *Client) GetEvent(event_type string) *Event {
	event := new(Event)
	event.EventTypeName = event_type
	switch event_type {
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
	return event
}

func (client *Client) SubscriptionURL() (string, error) {
	if ip_address, err := GetInterfaceAddress(client.config.EventsInterface); err != nil {
		return "", err
	} else {
		/* step: construct the url */
		client.events_ipaddress = ip_address
		client.events_callback_url = fmt.Sprintf("http://%s:%d%s",
			client.events_ipaddress, client.config.EventsPort, DEFAULT_EVENTS_URL)

		client.Debug("Subscription callback url: %s", client.events_callback_url)
		return client.events_callback_url, nil
	}
}

func (client *Client) WatchList() []string {
	client.RLock()
	defer client.RUnlock()
	list := make([]string, 0)
	for name, _ := range client.services {
		list = append(list, name)
	}
	return list
}

func (client *Client) Watch(name string, channel chan string) {
	client.Lock()
	defer client.Unlock()
	client.services[name] = channel
}

func (client *Client) RemoveWatch(name string) {
	client.Lock()
	defer client.Unlock()
	delete(client.services, name)
}
