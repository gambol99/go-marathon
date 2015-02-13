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
	"net/http"
	"strings"
	"sync"
	"time"
	"fmt"
)

type EventSubscription struct {
	CallbackURL  string   `json:"CallbackUrl"`
	ClientIP     string   `json:"ClientIp"`
	EventType    string   `json:"eventType"`
	CallbackURLs []string `json:"CallbackUrls"`
}

type Subscriptions struct {
	CallbackURLs []string `json:"callbackUrls"`
}

var (
	subscriptionLock sync.Once
)

func (client *Client) Subscriptions() (Subscriptions, error) {
	var subscriptions Subscriptions
	if err := client.ApiGet(MARATHON_API_SUBSCRIPTION, "", &subscriptions); err != nil {
		return Subscriptions{}, err
	} else {
		return subscriptions, nil
	}
}

func (client *Client) RegisterSubscription() error {
	/* step: create the call back handler */
	subscriptionLock.Do(func() {
		/* step: register the handler */
		http.HandleFunc(DEFAULT_EVENTS_URL, client.HandleMarathonEvent)
		/* step: register and listen */
		go http.ListenAndServe(client.subscription_iface, nil)
		/* step: we register with the subscription service */
	})
	/* step: attempt to register with the marathon callback */
	attempts := 1
	max_attempts := 3
	for {
		uri := fmt.Sprintf("%s", MARATHON_API_SUBSCRIPTION)
		if err := client.ApiPost(uri, "", nil); err != nil {
			return err
		}
		/* check: have we reached the max attempts? */
		if attempts >= max_attempts {
			return ErrInvalidResponse
		}
		/* choice: lets go to sleep for x seconds */
		time.Sleep(3 * time.Second)
		attempts += 1
	}
}

func (client *Client) DeregisterSubscription() error {
	/* step: check if we are already subscripted */
	found, err := client.HasSubscription()
	if err != nil {
		return err
	} else if found {
		/* step: remove from the list of subscriptions */

	}
	return nil
}

func (client *Client) HasSubscription() (bool, error) {
	if subscriptions, err := client.Subscriptions(); err != nil {
		return false, err
	} else {
		for _, subscription := range subscriptions.CallbackURLs {
			if client.subscription_url == subscription {
				return true, nil
			}

		}
	}
	return false, nil
}

func (client *Client) HandleMarathonEvent(writer http.ResponseWriter, request *http.Request) {
	var event Event
	decoder := json.NewDecoder(request.Body)
	if err := decoder.Decode(&event); err != nil {

		//glog.Errorf("Failed to decode the Marathon event: %s, error: %s", request.Body, err )
	} else {
		switch event.EventType {
		case "health_status_changed_event":
			//glog.V(4).Infof("Marathon application: %s health status has been altered, resyncing", event.AppID)
		case "status_update_event":
			//glog.V(4).Infof("Marathon application: %s status update, resyncing endpoints", event.AppID)
		default:
			//glog.V(10).Infof("Skipping the Marathon event, as it's not a status update, type: %s", event.EventType)
			return
		}
		/* step: we notify the receiver */
		for service, listener := range client.services {
			if strings.HasPrefix(service, event.AppID) {

				go func() {
					listener <- true
				}()
			}
		}
	}
}

func (client *Client) WatchList() []string {
	client.RLock()
	defer client.RUnlock()
	list := make([]string,0)
	for name, _ := range client.services {
		list = append(list, name)
	}
	return list
}

func (client *Client) Watch(name string, channel chan bool) {
	client.Lock()
	defer client.Unlock()
	client.services[name] = channel
}

func (client *Client) RemoveWatch(name string, channel chan bool) {
	client.Lock()
	defer client.Unlock()
	delete(client.services, name)
}
