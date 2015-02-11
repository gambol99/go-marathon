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
	"time"
	"encoding/json"
	"net/http"
	"fmt"
	"strings"
)

type EventSubscription struct {
	CallbackURL  string   `json:"CallbackUrl"`
	ClientIP     string   `json:"ClientIp"`
	EventType    string   `json:"eventType"`
	CallbackURLs []string `json:"CallbackUrls"`
}

func (client *MarathonClient) RegisterEvents() error {
	registration := fmt.Sprintf("%s/%s?callbackUrl=%s", r.marathon_url, MARATHON_API_SUBSCRIPTION, r.callback_url)

	/* step: register the http handler and start listening */
	http.HandleFunc(DEFAULT_EVENTS_URL, client.HandleMarathonEvent)
	go func() {
		//glog.Infof("Starting to listen to http events from Marathon on %s", service.marathon_url)
		//http.ListenAndServe(config.service_interface, nil)
	}()

	/* step: attempt to register with the marathon callback */
	attempts := 1
	max_attempts := 3
	for {
		if response, err := http.Post(registration, "application/json", nil); err != nil {
			//glog.Errorf("Failed to post Marathon registration for callback service, error: %s", err )
		} else {
			if response.StatusCode < 200 || response.StatusCode >= 300 {
				//glog.Errorf("Failed to register with the Marathon event callback service, error: %s", response.Body)
			} else {
				//glog.Infof("Successfully registered with Marathon to receive events")
				return nil
			}
		}
		/* check: have we reached the max attempts? */
		if attempts >= max_attempts {
			/* choice: if after x attempts we can't register with Marathon, there's not much point */
			//glog.Fatalf("Failed to register with Marathon's callback service %d time, no point in continuing", attempts)
		}

		/* choice: lets go to sleep for x seconds */
		time.Sleep(3 * time.Second)
		attempts += 1
	}
	return nil
}

func (client *MarathonClient) DeregisterEvents(callback string, marathon string) error {
	/** @@TODO := needs to be implemented, not to leave loose callbacks around */
	return nil
}

func (client *MarathonClient) HandleMarathonEvent(writer http.ResponseWriter, request *http.Request) {
	var event MarathonEvent
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
			//glog.Infof("FOUND SERVICE, key: %s, channel: %v", service, listener)
			if strings.HasPrefix(service, event.AppID) {
				//glog.Infof("SENDING EVENT, key: %s, channel: %v", service, listener)
				go func() {
					listener <- true
				}()
			}
		}
	}
}

func (client MarathonClient) Watch(service_name string, service_port int, channel chan bool) {
	client.Lock()
	defer client.Unlock()
	service_key := client.GetServiceKey(service_name, service_port)
	client.services[service_key] = channel
}

func (client MarathonClient) RemoveWatch(service_name string, service_port int, channel chan bool) {
	client.Lock()
	defer client.Unlock()
	service_key := client.GetServiceKey(service_name, service_port)
	delete(client.services,service_key)
}


