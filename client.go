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
	"strings"
	"fmt"
	"errors"
	"net/http"
	"net/url"
	"io/ioutil"
	"encoding/json"
)

type Marathon interface {
	/* watch for changes on a application */
	Watch(application_id string, channel chan bool)
	/* remove me from watching this service */
	RemoveWatch(application_id string, channel chan bool)
	/* get a list of applications from marathon */
	Applications() (Applications,error)
	/* get a specific application */
	Application(id string) (Application, error)
	/* get a list of tasks for a specific application */
	Tasks(id string) (Tasks, error)
	/* get a list of all tasks */
	AllTasks() (Tasks, error)
	/* get the marathon url */
	GetMarathonURL() string
	/* get the call back url */
	GetCallbackURL() string
}

var (
	/* the url specified was invalid */
	ErrInvalidEndpoint = errors.New("Invalid Marathon endpoint specified")
	/* invalid or error response from marathon */
	ErrInvalidResponse = errors.New("Invalid response from Marathon")
	/* some resource does not exists */
	ErrDoesNotExist = errors.New("The resource does not exist")
)

type MarathonClient struct {
	/* the configuration for the client */
	config Config
	/* the marathon url */
	hosts []string
	/* protocol */
	protocol string
	/* the http clinet */
	http *http.Client
}

func NewClient(config Config) (Marathon, error) {
	/* step: we need to get the ip address of the interface */
	var ip_address string
	var err errors

	/* step: get the ip address, will be required for call backs */
	if config.event_ipaddress != "" {

	} else {
		ip_address, err = GetInterfaceAddress(config.Options.Proxy_interface)
		if err != nil {
			return nil, err
		}
	}

	/* step: create the service */
	service := new(MarathonClient)
	service.services = make(map[string]chan bool,0)
	service.marathon_url = fmt.Sprintf("http://%s", strings.TrimPrefix(config.marathon_url, "marathon://") )
	service.http = http.Transport{Dial: 5}
	/* step: register with marathon service as a callback for events */
	service.service_interface = fmt.Sprintf("%s:%d", ip_address, config.events_port)
	service.callback_url 	  = fmt.Sprintf("http://%s%s", service.service_interface, DEFAULT_EVENTS_URL)
	return service, nil
}

func (client *MarathonClient) GetMarathonURL() string {
	return client.marathon_url
}

func (client *MarathonClient) GetCallbackURL() string {
	return client.callback_url
}

func (client *MarathonClient) GetServiceKey(service_name string, service_port int) string {
	return fmt.Sprintf("%s:%d", service_name, service_port)
}

func (client *MarathonClient) ParseMarathonURL(uri string) error {
	if marathon, err := url.Parse(uri); err != nil {
		return ErrInvalidEndpoint
	} else {
		/* check the protocol */
		if marathon.Scheme != "http" && marathon.Scheme != "https" {
			return errors.New("Invalid protocol type for marathon url, must be http/https")
		}
		client.hosts := strings.SplitN(marathon.Host, ",", -1)
		client.protocol = marathon.Scheme
	}
	return nil
}

func (client *MarathonClient) ApiGet(uri string, response *interface {}) error {
	if result, _, err := client.HttpGet(uri); err != nil {
		return err
	} else {
		decoder := json.NewDecoder(result)
		if err := decoder.Decode(response); err != nil {
			return err
		} else {
			return nil
		}
	}
}

func (client *MarathonClient) ApiPost(uri string, post interface {}, result interface {}) error {


	return nil
}

func (client *MarathonClient) ApiDelete(uri string, result interface {}) error {

	return nil
}

func (client *MarathonClient) HttpGet(uri string) (string, int, error) {
	/* step: we can try any of the endpoints */
	for _, marathon := range client.hosts {
		/* @@todo will move this over to a cluster formation later */
		request_url := fmt.Sprintf("%s://%s%s", client.protocol, marathon, uri)
		if response, err := client.http.Get(request_url); err != nil {
			/* step: lets try another host perhaps? */
			continue
		} else {
			/* step: lets read in the http body */
			if body, err := ioutil.ReadAll(response.Body); err != nil {
				return "", 0, err
			} else {
				status_code := response.StatusCode
				if status_code >= 200 || status_code <= 299 {
					return body, status_code, nil
				} else {
					switch status_code {
					case 404:
						return "", status_code, ErrDoesNotExist
					default:
						return body, status_code, ErrInvalidResponse
					}
				}
			}
		}
	}
	return "", 0, errors.New("Unable to make call to marathon")
}
