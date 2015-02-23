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

type Deployment struct {
	AffectedApps   []string          `json:"affectedApps"`
	ID             string            `json:"id"`
	Steps          []*DeploymentStep `json:"steps"`
	CurrentActions []*DeploymentStep `json:"currentActions"`
	CurrentStep    int               `json:"currentStep"`
	TotalSteps     int               `json:"totalSteps"`
	Version        string            `json:"version"`
}

type DeploymentID struct {
	DeploymentID string `json:"deploymentId"`
	Version      string `json:"version"`
}

type DeploymentStep struct {
	Action string `json:"action"`
	App    string `json:"app"`
}

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

func (client *Client) Deployments() ([]Deployment, error) {
	var deployments []Deployment
	if err := client.ApiGet(MARATHON_API_DEPLOYMENTS, "", &deployments); err != nil {
		return nil, err
	} else {
		return deployments, nil
	}
}

func (client *Client) DeleteDeployment(deployment Deployment, force bool) (Deployment, error) {
	var result Deployment
	if err := client.ApiDelete(fmt.Sprintf("%s/%s", MARATHON_API_DEPLOYMENTS, deployment.ID), nil, &result); err != nil {
		return Deployment{}, err
	} else {
		return result, nil
	}
}
