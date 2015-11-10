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
	"time"
)

// Deployment is a marathon deployment definition
type Deployment struct {
	ID             string              `json:"id"`
	Version        string              `json:"version"`
	CurrentStep    int                 `json:"currentStep"`
	TotalSteps     int                 `json:"totalSteps"`
	AffectedApps   []string            `json:"affectedApps"`
	Steps          [][]*DeploymentStep `json:"steps"`
	CurrentActions []*DeploymentStep   `json:"currentActions"`
}

// DeploymentID is the identifier for a application deployment
type DeploymentID struct {
	DeploymentID string `json:"deploymentId"`
	Version      string `json:"version"`
}

// DeploymentStep is a step in the application deployment plan
type DeploymentStep struct {
	Action string `json:"action"`
	App    string `json:"app"`
}

// DeploymentPlan is a collection of steps for application deployment
type DeploymentPlan struct {
	ID       string `json:"id"`
	Version  string `json:"version"`
	Original struct {
		Apps         []*Application `json:"apps"`
		Dependencies []string       `json:"dependencies"`
		Groups       []*Group       `json:"groups"`
		ID           string         `json:"id"`
		Version      string         `json:"version"`
	} `json:"original"`
	Steps  []*DeploymentStep `json:"steps"`
	Target struct {
		Apps         []*Application `json:"apps"`
		Dependencies []string       `json:"dependencies"`
		Groups       []*Group       `json:"groups"`
		ID           string         `json:"id"`
		Version      string         `json:"version"`
	} `json:"target"`
}

// Deployments retrieves a list of current deployments
func (r *marathonClient) Deployments() ([]*Deployment, error) {
	var deployments []*Deployment
	err := r.apiGet(marathonAPIDeployments, nil, &deployments)
	if err != nil {
		return nil, err
	}

	return deployments, nil
}

// DeleteDeployment delete a current deployment from marathon
// 	id:		the deployment id you wish to delete
// 	force:	whether or not to force the deletion
func (r *marathonClient) DeleteDeployment(id string, force bool) (*DeploymentID, error) {
	deployment := new(DeploymentID)
	err := r.apiDelete(fmt.Sprintf("%s/%s", marathonAPIDeployments, id), nil, deployment)
	if err != nil {
		return nil, err
	}

	return deployment, nil
}

// HasDeployment checks to see if a deployment exists
// 	id:		the deployment id you are looking for
func (r *marathonClient) HasDeployment(id string) (bool, error) {
	deployments, err := r.Deployments()
	if err != nil {
		return false, err
	}
	for _, deployment := range deployments {
		if deployment.ID == id {
			return true, nil
		}
	}
	return false, nil
}

// WaitOnDeployment waits on a deployment to finish
//  version:		the version of the application
// 	timeout:		the timeout to wait for the deployment to take, otherwise return an error
func (r *marathonClient) WaitOnDeployment(id string, timeout time.Duration) error {
	if found, err := r.HasDeployment(id); err != nil {
		return err
	} else if !found {
		return nil
	}

	nowTime := time.Now()
	stopTime := nowTime.Add(timeout)
	if timeout <= 0 {
		stopTime = nowTime.Add(time.Duration(900) * time.Second)
	}

	// step: a somewhat naive implementation, but it will work
	for {
		if time.Now().After(stopTime) {
			return ErrTimeoutError
		}
		found, err := r.HasDeployment(id)
		if err != nil {
			return err
		}
		if !found {
			return nil
		}
		time.Sleep(time.Duration(2) * time.Second)
	}
}
