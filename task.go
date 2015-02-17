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
	"fmt"
)

type Tasks struct {
	Tasks []Task `json:"tasks"`
}

type Task struct {
	AppID             string               `json:"appId"`
	Host              string               `json:"host"`
	ID                string               `json:"id"`
	HealthCheckResult []*HealthCheckResult `json:"healthCheckResults"`
	Ports             []int                `json:"ports"`
	ServicePorts      []int                `json:"servicePorts"`
	StagedAt          string               `json:"stagedAt"`
	StartedAt         string               `json:"startedAt"`
	Version           string               `json:"version"`
}

func (task Task) String() string {
	return fmt.Sprintf("id: %s, application: %s, host: %s, ports: %s, created: %s",
		task.ID, task.AppID, task.Host, task.Ports, task.StartedAt)
}

func (client *Client) AllTasks() (*Tasks, error) {
	tasks := new(Tasks)
	if err := client.ApiGet(MARATHON_API_TASKS, "", tasks); err != nil {
		return nil, err
	} else {
		return tasks, nil
	}
}

func (client *Client) Tasks(application_id string) (*Tasks, error) {
	tasks := new(Tasks)
	if err := client.ApiGet(fmt.Sprintf("%s%s/tasks", MARATHON_API_APPS, application_id), "", tasks); err != nil {
		return nil, err
	} else {
		return tasks, nil
	}
}

// Get the endpoints i.e. HOST_IP:DYNAMIC_PORT for a specific application service
// I.e. a container running apache, might have ports 80/443 (translated to X dynamic ports), but i want
// port 80 only and i only want those whom have passed the health check
// Params:
//		name:		the identifier for the application
//		port:		the container port you are interested in
//		health: 	whether to check the health or not
func (client *Client) TaskEndpoints(name string, port int) ([]string, error) {
	/* step: get the application details */
	if application, err := client.Application(name); err != nil {
		return nil, err
	} else {
		/* step: we need to get the port index of the service we are interested in */
		if port_index, err := application.Container.Docker.ServicePortIndex(port); err != nil {
			return nil, err
		} else {
			var _ = port_index
		}
	}
	return nil, nil
}
