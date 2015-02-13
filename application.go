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
	"errors"
	"fmt"

)

var (
	ErrApplicationExists = errors.New("The application already exists in marathon, you must update")
)

type Applications struct {
	Apps []Application `json:"apps"`
}

type ApplicationWrap struct {
	Application Application	`json:"app"`
}

type Application struct {
	ID            string            `json:"id"`
	Cmd           string            `json:"cmd,omitempty"`
	Constraints   [][]string        `json:"constraints,omitempty"`
	Container     *Container        `json:"container,omitempty"`
	CPUs          float32           `json:"cpus,omitempty"`
	Env           map[string]string `json:"env,omitempty"`
	Executor      string            `json:"executor,omitempty"`
	HealthChecks  []*HealthCheck    `json:"healthChecks,omitempty"`
	Instances     int               `json:"instances,omitemptys"`
	Mem           float32           `json:"mem,omitempty"`
	Tasks         []*Task           `json:"tasks,omitempty"`
	Ports         []int             `json:"ports,omitempty"`
	RequirePorts  bool              `json:"requirePorts,omitempty"`
	BackoffFactor float32           `json:"backoffFactor,omitempty"`
	TasksRunning  int               `json:"tasksRunning,omitempty"`
	TasksStaged   int               `json:"tasksStaged,omitempty"`
	Uris          []string          `json:"uris,omitempty"`
	Version       string            `json:"version,omitempty"`
}

func (client *Client) Applications() (Applications, error) {
	var apps Applications
	if err := client.ApiGet(MARATHON_API_APPS, "", &apps); err != nil {
		return Applications{}, err
	} else {
		return apps, nil
	}
}

func (client *Client) ListApplications() ([]string, error) {
	if applications, err := client.Applications(); err != nil {
		return nil, err
	} else {
		list := make([]string, 0)
		for _, application := range applications.Apps {
			list = append(list, application.ID)
		}
		return list, nil
	}
}

func (client *Client) Application(id string) (Application, error) {
	var application ApplicationWrap
	if err := client.ApiGet(fmt.Sprintf("%s%s", MARATHON_API_APPS, id), "", &application); err != nil {
		return Application{}, err
	} else {
		return application.Application, nil
	}
}

func (client *Client) CreateApplication(application Application) (bool, error) {
	/* step: check of the application already exists */
	if found, err := client.HasApplication(application.ID); err != nil {
		return false, err
	} else if found {
		return false, ErrApplicationExists
	}
	/* step: post the application to marathon */
	if err := client.ApiPost(MARATHON_API_APPS, &application, nil); err != nil {
		return false, err
	}
	return true, nil
}

func (client *Client) HasApplication(name string) (bool, error) {
	if applications, err := client.ListApplications(); err != nil {
		return false, err
	} else {
		for _, id := range applications {
			if name == id {
				return true, nil
			}
		}
	}
	return false, nil
}

func (client *Client) DeleteApplication(app Application) (bool, error) {

	return false, nil
}

func (client *Client) RestartApplication(app Application, force bool) (Deployment, error) {

	return Deployment{}, nil
}
