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
	"net/url"
	"strconv"
	"strings"
)

// Tasks ... a collection of marathon tasks
type Tasks struct {
	Tasks []Task `json:"tasks"`
}

// Task ... the definition for a marathon task
type Task struct {
	ID                string               `json:"id"`
	AppID             string               `json:"appId"`
	Host              string               `json:"host"`
	HealthCheckResult []*HealthCheckResult `json:"healthCheckResults"`
	Ports             []int                `json:"ports"`
	ServicePorts      []int                `json:"servicePorts"`
	StagedAt          string               `json:"stagedAt"`
	StartedAt         string               `json:"startedAt"`
	Version           string               `json:"version"`
}

// String returns a string representation of the struct
func (r Task) String() string {
	return fmt.Sprintf("id: %s, application: %s, host: %s, ports: %v, created: %s",
		r.ID, r.AppID, r.Host, r.Ports, r.StartedAt)
}

// HasHealthCheckResults ... Check if the task has any health checks
func (r *Task) HasHealthCheckResults() bool {
	if r.HealthCheckResult == nil || len(r.HealthCheckResult) <= 0 {
		return false
	}
	return true
}

// AllTasks ... Retrieve all the tasks currently running
func (r *marathonClient) AllTasks() (*Tasks, error) {
	tasks := new(Tasks)
	if err := r.apiGet(MARATHON_API_TASKS, nil, tasks); err != nil {
		return nil, err
	}

	return tasks, nil
}

// Tasks ... Retrieve a list of tasks for an application
//		application_id:		the id for the application
func (r *marathonClient) Tasks(id string) (*Tasks, error) {
	tasks := new(Tasks)
	if err := r.apiGet(fmt.Sprintf("%s/%s/tasks", MARATHON_API_APPS, trimRootPath(id)), nil, tasks); err != nil {
		return nil, err
	}

	return tasks, nil
}

// ListTasks ... Retrieve an array of task ids currently running in marathon
func (r *marathonClient) ListTasks() ([]string, error) {
	tasks, err := r.AllTasks()
	if err != nil {
		return nil, err
	}
	list := make([]string, 0)
	for _, task := range tasks.Tasks {
		list = append(list, task.ID)
	}

	return list, nil
}

// KillApplicationTasks ... Kill all tasks relating to an application
//		application_id:		the id for the application
//      host:				kill only those tasks on a specific host (optional)
//		scale:              Scale the app down (i.e. decrement its instances setting by the number of tasks killed) after killing the specified tasks
func (r *marathonClient) KillApplicationTasks(id, hostname string, scale bool) (*Tasks, error) {
	var options struct {
		Host  string `json:"host"`
		Scale bool   `json:"bool"`
	}
	options.Host = hostname
	options.Scale = scale
	tasks := new(Tasks)
	if err := r.apiDelete(fmt.Sprintf("%s/%s/tasks", MARATHON_API_APPS, trimRootPath(id)), &options, tasks); err != nil {
		return nil, err
	}

	return tasks, nil
}

// KillTask ... Kill the task associated with a given ID
// 	task_id:		the id for the task
// 	scale:		Scale the app down
func (r *marathonClient) KillTask(taskId string, scale bool) (*Task, error) {
	var options struct {
		Scale bool `json:"bool"`
	}
	options.Scale = scale
	task := new(Task)
	appName := taskId[0:strings.LastIndex(taskId, ".")]
	if err := r.apiDelete(fmt.Sprintf("%s/%s/tasks/%s", MARATHON_API_APPS, appName, taskId), &options, task); err != nil {
		return nil, err
	}

	return task, nil
}

// KillTasks ... Kill tasks associated with given array of ids
// 	tasks: 	the array of task ids
// 	scale: 	Scale the app down
func (r *marathonClient) KillTasks(tasks []string, scale bool) error {
	v := url.Values{}
	v.Add("scale", strconv.FormatBool(scale))
	var post struct {
		TaskIDs []string `json:"ids"`
	}
	post.TaskIDs = tasks

	return r.apiPost(fmt.Sprintf("%s/delete?%s", MARATHON_API_TASKS, v.Encode()), &post, nil)
}

// TaskEndpoints ... Get the endpoints i.e. HOST_IP:DYNAMIC_PORT for a specific application service
// I.e. a container running apache, might have ports 80/443 (translated to X dynamic ports), but i want
// port 80 only and i only want those whom have passed the health check
//
// Note: I've NO IDEA how to associate the health_check_result to the actual port, I don't think it's
// possible at the moment, however, given marathon will fail and restart an application even if one of x ports of a task is
// down, the per port check is redundant??? .. personally, I like it anyhow, but hey
//

//		name:		the identifier for the application
//		port:		the container port you are interested in
//		health: 	whether to check the health or not
func (r *marathonClient) TaskEndpoints(name string, port int, health_check bool) ([]string, error) {
	// step: get the application details
	application, err := r.Application(name)
	if err != nil {
		return nil, err
	}

	// step: we need to get the port index of the service we are interested in
	port_index, err := application.Container.Docker.ServicePortIndex(port)
	if err != nil {
		return nil, err
	}

	list := make([]string, 0)
	// step: do we have any tasks?
	if application.Tasks == nil || len(application.Tasks) <= 0 {
		return list, nil
	}

	// step: iterate the tasks and extract the dynamic ports
	for _, task := range application.Tasks {
		// step: if we are checking health the 'service' has a health check?
		if health_check && application.HasHealthChecks() {
			/*
				check: does the task have a health check result, if NOT, it's because the
				health of the task hasn't yet been performed, hence we assume it as DOWN
			*/
			if task.HasHealthCheckResults() == false {
				continue
			}

			// step: check the health results then
			skip_endpoint := false
			for _, health := range task.HealthCheckResult {
				if health.Alive == false {
					skip_endpoint = true
				}
			}

			if skip_endpoint == true {
				continue
			}
		}
		// else we can just add it
		list = append(list, fmt.Sprintf("%s:%d", task.Host, task.Ports[port_index]))
	}

	return list, nil
}
