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
	"io"
	"io/ioutil"
	"net/http"
	"strings"
	"sync"
	"time"
)

const (
	HTTP_GET    = "GET"
	HTTP_PUT    = "PUT"
	HTTP_DELETE = "DELETE"
	HTTP_POST   = "POST"
)

type Marathon interface {
	/* -- APPLICATIONS --- */

	/* check it see if a application exists */
	HasApplication(name string) (bool, error)
	/* get a listing of the application ids */
	ListApplications() ([]string, error)
	/* a list of application versions */
	ApplicationVersions(name string) (*ApplicationVersions, error)
	/* check a application version exists */
	HasApplicationVersion(name, version string) (bool, error)
	/* change an application to a different version */
	SetApplicationVersion(name string, version *ApplicationVersion) (*DeploymentID, error)
	/* check if an application is ok */
	ApplicationOK(name string) (bool, error)
	/* create an application in marathon */
	CreateApplication(application *Application) error
	/* delete an application */
	DeleteApplication(name string) error
	/* scale a application */
	ScaleApplicationInstances(name string, instances int) error
	/* restart an application */
	RestartApplication(name string, force bool) (*DeploymentID, error)
	/* get a list of applications from marathon */
	Applications() (*Applications, error)
	/* get a specific application */
	Application(name string) (*Application, error)

	/* -- TASKS --- */

	/* get a list of tasks for a specific application */
	Tasks(application string) (*Tasks, error)
	/* get a list of all tasks */
	AllTasks() (*Tasks, error)

	/* --- GROUPS --- */

	/* list all the groups in the system */
	Groups() (*Groups, error)
	/* retrieve a specific group from marathon */
	Group(name string) (*Group, error)
	/* create a group deployment */
	CreateGroup(group *Group) (*ApplicationVersion, error)
	/* delete a group */
	DeleteGroup(name string) (*ApplicationVersion, error)
	/* check if a group exists */
	HasGroup(name string) (bool, error)

	/* --- DEPLOYMENTS --- */

	/* get a list of the deployments */
	Deployments() ([]Deployment, error)
	/* delete a deployment */
	DeleteDeployment(deployment Deployment, force bool) (Deployment, error)

	/* --- SUBSCRIPTIONS --- */

	/* a list of current subscriptions */
	Subscriptions() (*Subscriptions, error)
	/* add a events listener */
	AddEventsListener(channel EventsChannel, filter int) error
	/* remove a events listener */
	RemoveEventsListener(channel EventsChannel)
	/* notify me of changes to application */

	/* --- MISC --- */

	/* get the marathon url */
	GetMarathonURL() string
	/* ping the marathon */
	Ping() (bool, error)
	/* grab the marathon server info */
	Info() (*Info, error)
}

var (
	/* the url specified was invalid */
	ErrInvalidEndpoint = errors.New("Invalid Marathon endpoint specified")
	/* invalid or error response from marathon */
	ErrInvalidResponse = errors.New("Invalid response from Marathon")
	/* some resource does not exists */
	ErrDoesNotExist = errors.New("The resource does not exist")
	/* all the marathon endpoints are down */
	ErrMarathonDown = errors.New("All the Marathon hosts are presently down")
	/* unable to decode the response */
	ErrInvalidResult = errors.New("Unable to decode the response from Marathon")
	/* invalid argument */
	ErrInvalidArgument = errors.New("The argument passed is invalid")
	/* error return by marathon */
	ErrMarathonError = errors.New("Marathon error")
)

type Client struct {
	sync.RWMutex
	/* the configuration for the client */
	config Config
	/* the ip addess of the client */
	ipaddress string
	/* the http server */
	events_http *http.Server
	/* the http client */
	http *http.Client
	/* the marathon cluster */
	cluster Cluster
	/* a map of service you wish to listen to */
	listeners map[EventsChannel]int
}

type Message struct {
	Message string `json:"message"`
}

func NewClient(config Config) (Marathon, error) {
	/* step: we parse the url and build a cluster */
	if cluster, err := NewMarathonCluster(config.URL); err != nil {
		return nil, err
	} else {
		/* step: create the service marathon client */
		service := new(Client)
		service.config = config
		service.listeners = make(map[EventsChannel]int, 0)
		service.cluster = cluster
		service.http = &http.Client{
			Timeout: (time.Duration(config.RequestTimeout) * time.Second),
		}
		return service, nil
	}
}

func (client *Client) GetMarathonURL() string {
	return client.cluster.Url()
}

func (client *Client) Ping() (bool, error) {
	if err := client.ApiGet(MARATHON_API_PING, "", nil); err != nil {
		return false, err
	} else {
		return true, nil
	}
}

func (client *Client) MarshallJSON(data interface{}) (string, error) {
	if response, err := json.Marshal(data); err != nil {
		return "", err
	} else {
		return string(response), err
	}
}

func (client *Client) UnMarshallDataToJson(stream io.Reader, result interface{}) error {
	decoder := json.NewDecoder(stream)
	if err := decoder.Decode(result); err != nil {
		return err
	}
	return nil
}

func (client *Client) ApiGet(uri, body string, result interface{}) error {
	client.Debug("ApiGet() uri: %s, body: %s", uri, body)
	_, _, error := client.ApiCall(HTTP_GET, uri, body, result)
	return error
}

func (client *Client) ApiPut(uri string, post interface{}, result interface{}) error {
	var content string
	var err error
	if post == nil {
		content = ""
	} else {
		content, err = client.MarshallJSON(post)
		if err != nil {
			return err
		}
	}
	_, _, error := client.ApiCall(HTTP_PUT, uri, content, result)
	return error
}

func (client *Client) ApiPost(uri string, post interface{}, result interface{}) error {
	/* step: we need to marshall the post data into json */
	var content string
	var err error
	if post == nil {
		content = ""
	} else {
		content, err = client.MarshallJSON(post)
		if err != nil {
			return err
		}
	}
	_, _, error := client.ApiCall(HTTP_POST, uri, content, result)
	return error
}

func (client *Client) ApiDelete(uri, body string, result interface{}) error {
	_, _, error := client.ApiCall(HTTP_DELETE, uri, body, result)
	return error
}

func (client *Client) ApiCall(method, uri, body string, result interface{}) (int, string, error) {
	client.Debug("ApiCall() method: %s, uri: %s, body: %s", method, uri, body)
	if status, content, _, err := client.HttpCall(method, uri, body); err != nil {
		return 0, "", err
	} else {
		client.Debug("ApiCall() status: %s, content: %s\n", status, content)
		if status >= 200 && status <= 299 {
			if result != nil {
				if err := client.UnMarshallDataToJson(strings.NewReader(content), result); err != nil {
					return status, content, err
				}
			}
			return status, content, nil
		}
		switch status {
		case 500:
			return status, "", ErrInvalidResponse
		case 404:
			return status, "", ErrDoesNotExist
		}

		/* step: lets decode into a error message */
		var message Message
		if err := client.UnMarshallDataToJson(strings.NewReader(content), &message); err != nil {
			return status, content, ErrInvalidResponse
		} else {
			return status, message.Message, ErrMarathonError
		}
	}
}

func (client *Client) HttpCall(method, uri, body string) (int, string, *http.Response, error) {
	/* step: get a member from the cluster */
	if marathon, err := client.cluster.GetMember(); err != nil {
		return 0, "", nil, err
	} else {
		url := fmt.Sprintf("%s/%s", marathon, uri)
		client.Debug("HTTP method: %s, uri: %s, url: %s", method, uri, url)

		if request, err := http.NewRequest(method, url, strings.NewReader(body)); err != nil {
			return 0, "", nil, err
		} else {
			request.Header.Add("Content-Type", "application/json")
			var content string
			/* step: perform the request */
			if response, err := client.http.Do(request); err != nil {
				/* step: mark the endpoint as down */
				client.cluster.MarkDown()
				/* step: retry the request with another endpoint */
				return client.HttpCall(method, uri, body)
			} else {
				/* step: lets read in any content */
				client.Debug("HTTP method: %s, uri: %s, url: %s\n", method, uri, url)
				if response.ContentLength != 0 {
					/* step: read in the content from the request */
					response_content, err := ioutil.ReadAll(response.Body)
					if err != nil {
						return response.StatusCode, "", response, err
					}
					content = string(response_content)
				}
				/* step: return the request */
				return response.StatusCode, content, response, nil
			}
		}
	}
	return 0, "", nil, errors.New("Unable to make call to marathon")
}
