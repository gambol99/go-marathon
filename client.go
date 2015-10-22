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
	"log"
	"net/http"
	"net/url"
	"strings"
	"sync"
	"time"

	"github.com/golang/glog"
)

// Marathon is the interface to the marathon API
type Marathon interface {
	// -- APPLICATIONS ---

	// check it see if a application exists
	HasApplication(name string) (bool, error)
	// get a listing of the application ids
	ListApplications(url.Values) ([]string, error)
	// a list of application versions
	ApplicationVersions(name string) (*ApplicationVersions, error)
	// check a application version exists
	HasApplicationVersion(name, version string) (bool, error)
	// change an application to a different version
	SetApplicationVersion(name string, version *ApplicationVersion) (*DeploymentID, error)
	// check if an application is ok
	ApplicationOK(name string) (bool, error)
	// create an application in marathon
	CreateApplication(application *Application) (*Application, error)
	// delete an application
	DeleteApplication(name string) (*DeploymentID, error)
	// update an application in marathon
	UpdateApplication(application *Application) (*DeploymentID, error)
	// a list of deployments on a application
	ApplicationDeployments(name string) ([]*DeploymentID, error)
	// scale a application
	ScaleApplicationInstances(name string, instances int, force bool) (*DeploymentID, error)
	// restart an application
	RestartApplication(name string, force bool) (*DeploymentID, error)
	// get a list of applications from marathon
	Applications(url.Values) (*Applications, error)
	// get a specific application
	Application(name string) (*Application, error)
	// wait of application
	WaitOnApplication(name string, timeout time.Duration) error

	// -- TASKS ---

	// get a list of tasks for a specific application
	Tasks(application string) (*Tasks, error)
	// get a list of all tasks
	AllTasks() (*Tasks, error)
	// get a listing of the task ids
	ListTasks() ([]string, error)
	// get the endpoints for a service on a application
	TaskEndpoints(name string, port int, healthCheck bool) ([]string, error)
	// kill all the tasks for any application
	KillApplicationTasks(applicationID, hostname string, scale bool) (*Tasks, error)
	// kill a single task
	KillTask(taskID string, scale bool) (*Task, error)
	// kill the given array of tasks
	KillTasks(taskIDs []string, scale bool) error

	// --- GROUPS ---

	// list all the groups in the system
	Groups() (*Groups, error)
	// retrieve a specific group from marathon
	Group(name string) (*Group, error)
	// create a group deployment
	CreateGroup(group *Group) error
	// delete a group
	DeleteGroup(name string) (*DeploymentID, error)
	// update a groups
	UpdateGroup(id string, group *Group) (*DeploymentID, error)
	// check if a group exists
	HasGroup(name string) (bool, error)
	// wait for an group to be deployed
	WaitOnGroup(name string, timeout time.Duration) error

	// --- DEPLOYMENTS ---

	// get a list of the deployments
	Deployments() ([]*Deployment, error)
	// delete a deployment
	DeleteDeployment(id string, force bool) (*DeploymentID, error)
	// check to see if a deployment exists
	HasDeployment(id string) (bool, error)
	// wait of a deployment to finish
	WaitOnDeployment(id string, timeout time.Duration) error

	// --- SUBSCRIPTIONS ---

	// a list of current subscriptions
	Subscriptions() (*Subscriptions, error)
	// add a events listener
	AddEventsListener(channel EventsChannel, filter int) error
	// remove a events listener
	RemoveEventsListener(channel EventsChannel)
	// remove our self from subscriptions
	UnSubscribe() error

	// --- MISC ---

	// get the marathon url
	GetMarathonURL() string
	// ping the marathon
	Ping() (bool, error)
	// grab the marathon server info
	Info() (*Info, error)
	// retrieve the leader info
	Leader() (string, error)
	// cause the current leader to abdicate
	AbdicateLeader() (string, error)
}

var (
	// ErrInvalidEndpoint is thrown when the marathon url specified was invalid
	ErrInvalidEndpoint = errors.New("invalid Marathon endpoint specified")
	// ErrInvalidResponse is thrown when marathon responds with invalid or error response
	ErrInvalidResponse = errors.New("invalid response from Marathon")
	// ErrDoesNotExist is thrown when the resource does not exists
	ErrDoesNotExist = errors.New("the resource does not exist")
	// ErrMarathonDown is thrown when all the marathon endpoints are down
	ErrMarathonDown = errors.New("all the Marathon hosts are presently down")
	// ErrInvalidArgument is thrown when invalid argument
	ErrInvalidArgument = errors.New("the argument passed is invalid")
	// ErrTimeoutError is thrown when the operation has timed out
	ErrTimeoutError = errors.New("the operation has timed out")
)

type marathonClient struct {
	sync.RWMutex
	// the configuration for the client
	config Config
	// the ip address of the client
	ipAddress string
	// the http server */
	eventsHTTP *http.Server
	// the http client use for making requests
	httpClient *http.Client
	// the output for the logger
	logger *log.Logger
	// the marathon cluster
	cluster Cluster
	// a map of service you wish to listen to
	listeners map[EventsChannel]int
}

// NewClient creates a new marathon client
//		config:			the configuration to use
func NewClient(config Config) (Marathon, error) {
	// step: we parse the url and build a cluster
	cluster, err := newCluster(config.URL)
	if err != nil {
		return nil, err
	}

	service := new(marathonClient)
	service.config = config
	service.listeners = make(map[EventsChannel]int, 0)
	service.cluster = cluster
	service.httpClient = &http.Client{
		Timeout: (time.Duration(config.RequestTimeout) * time.Second),
	}

	return service, nil
}

// GetMarathonURL retrieves the marathon url
func (r *marathonClient) GetMarathonURL() string {
	return r.cluster.URL()
}

// Ping pings the current marathon endpoint (note, this is not a ICMP ping, but a rest api call)
func (r *marathonClient) Ping() (bool, error) {
	if err := r.apiGet(MARATHON_API_PING, nil, nil); err != nil {
		return false, err
	}
	return true, nil
}

func (r *marathonClient) encodeRequest(data interface{}) (string, error) {
	response, err := json.Marshal(data)
	if err != nil {
		return "", err
	}

	return string(response), err
}

func (r *marathonClient) decodeRequest(stream io.Reader, result interface{}) error {
	if err := json.NewDecoder(stream).Decode(result); err != nil {
		return err
	}

	return nil
}

func (r *marathonClient) buildPostData(data interface{}) (string, error) {
	if data == nil {
		return "", nil
	}
	content, err := r.encodeRequest(data)
	if err != nil {
		return "", err
	}

	return content, nil
}

func (r *marathonClient) apiGet(uri string, post, result interface{}) error {
	return r.apiOperation("GET", uri, post, result)
}

func (r *marathonClient) apiPut(uri string, post, result interface{}) error {
	return r.apiOperation("PUT", uri, post, result)
}

func (r *marathonClient) apiPost(uri string, post, result interface{}) error {
	return r.apiOperation("POST", uri, post, result)
}

func (r *marathonClient) apiDelete(uri string, post, result interface{}) error {
	return r.apiOperation("DELETE", uri, post, result)
}

func (r *marathonClient) apiOperation(method, uri string, post, result interface{}) error {
	content, err := r.buildPostData(post)
	if err != nil {
		return err
	}

	_, _, err = r.apiCall(method, uri, content, result)

	return err
}

func (r *marathonClient) apiCall(method, uri, body string, result interface{}) (int, string, error) {
	glog.V(DEBUG_LEVEL).Infof("[api]: method: %s, uri: %s, body: %s", method, uri, body)

	status, content, _, err := r.httpRequest(method, uri, body)
	if err != nil {
		return 0, "", err
	}

	glog.V(DEBUG_LEVEL).Infof("[api] result: status: %d, content: %s\n", status, content)
	if status >= 200 && status <= 299 {
		if result != nil {
			if err := r.decodeRequest(strings.NewReader(content), result); err != nil {
				glog.V(DEBUG_LEVEL).Infof("failed to unmarshall the response from marathon, error: %s", err)
				return status, content, ErrInvalidResponse
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

	// step: lets decode into a error message
	var message struct {
		Message string `json:"message"`
	}

	if err := r.decodeRequest(strings.NewReader(content), &message); err != nil {
		return status, content, ErrInvalidResponse
	}

	errorMessage := "unknown error"
	if message.Message != "" {
		errorMessage = message.Message
	}

	return status, "", fmt.Errorf("%s", errorMessage)
}

func (r *marathonClient) httpRequest(method, uri, body string) (int, string, *http.Response, error) {
	var content string

	// step: get a member from the cluster
	marathon, err := r.cluster.GetMember()
	if err != nil {
		return 0, "", nil, err
	}

	url := fmt.Sprintf("%s/%s", marathon, uri)

	glog.V(DEBUG_LEVEL).Infof("[http] request: %s, uri: %s, url: %s", method, uri, url)
	// step: make the http request to marathon
	request, err := http.NewRequest(method, url, strings.NewReader(body))
	if err != nil {
		return 0, "", nil, err
	}

	// step: add any basic auth and the content headers
	if r.config.HttpBasicAuthUser != "" {
		request.SetBasicAuth(r.config.HttpBasicAuthUser, r.config.HttpBasicPassword)
	}
	request.Header.Add("Content-Type", "application/json")
	request.Header.Add("Accept", "application/json")

	response, err := r.httpClient.Do(request)
	if err != nil {
		r.cluster.MarkDown()
		// step: retry the request with another endpoint
		return r.httpRequest(method, uri, body)
	}

	if response.ContentLength != 0 {
		responseContent, err := ioutil.ReadAll(response.Body)
		if err != nil {
			return response.StatusCode, "", response, err
		}
		content = string(responseContent)
	}

	return response.StatusCode, content, response, nil
}
