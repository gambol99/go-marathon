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
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestApplicationDependsOn(t *testing.T) {
	app := NewDockerApplication()
	app.DependsOn("fake_app")
	app.DependsOn("fake_app1")
	assert.Equal(t, 2, len(app.Dependencies))
}

func TestApplicationMemory(t *testing.T) {
	app := NewDockerApplication()
	app.Memory(50.0)
	assert.Equal(t, 50.0, app.Mem)
}

func TestApplicationCount(t *testing.T) {
	app := NewDockerApplication()
	assert.Equal(t, 0, app.Instances)
	app.Count(1)
	assert.Equal(t, 1, app.Instances)
}

func TestApplicationArgs(t *testing.T) {
	app := NewDockerApplication()
	assert.Equal(t, 0, len(app.Args))
	app.Arg("-p").Arg("option").Arg("-v")
	assert.Equal(t, 3, len(app.Args))
	assert.Equal(t, "-p", app.Args[0])
	assert.Equal(t, "option", app.Args[1])
	assert.Equal(t, "-v", app.Args[2])
}

func TestApplicationEnvs(t *testing.T) {
	app := NewDockerApplication()
	assert.Equal(t, 0, len(app.Env))
	app.AddEnv("hello", "world")
	assert.Equal(t, 1, len(app.Env))
}

func TestHasHealthChecks(t *testing.T) {
	app := NewDockerApplication()
	assert.False(t, app.HasHealthChecks())
	app.Container.Docker.Container("quay.io/gambol99/apache-php:latest").Expose(80)
	_, err := app.CheckTCP(80, 10)
	assert.Nil(t, err)
	assert.True(t, app.HasHealthChecks())
}

func TestApplicationCheckTCP(t *testing.T) {
	app := NewDockerApplication()
	assert.False(t, app.HasHealthChecks())
	_, err := app.CheckTCP(80, 10)
	assert.NotNil(t, err)
	assert.False(t, app.HasHealthChecks())
	app.Container.Docker.Container("quay.io/gambol99/apache-php:latest").Expose(80)
	_, err = app.CheckTCP(80, 10)
	assert.Nil(t, err)
	assert.True(t, app.HasHealthChecks())
	check := app.HealthChecks[0]
	if check == nil {
		t.FailNow()
	}
	assert.Equal(t, "TCP", check.Protocol)
	assert.Equal(t, 10, check.IntervalSeconds)
	assert.Equal(t, 0, check.PortIndex)
}

func TestApplicationCheckHTTP(t *testing.T) {
	app := NewDockerApplication()
	assert.False(t, app.HasHealthChecks())
	_, err := app.CheckHTTP("/", 80, 10)
	assert.NotNil(t, err)
	assert.False(t, app.HasHealthChecks())
	app.Container.Docker.Container("quay.io/gambol99/apache-php:latest").Expose(80)
	_, err = app.CheckHTTP("/health", 80, 10)
	assert.Nil(t, err)
	assert.True(t, app.HasHealthChecks())
	check := app.HealthChecks[0]
	if check == nil {
		t.FailNow()
	}
	assert.Equal(t, "HTTP", check.Protocol)
	assert.Equal(t, 10, check.IntervalSeconds)
	assert.Equal(t, "/health", check.Path)
	assert.Equal(t, 0, check.PortIndex)
}

func TestCreateApplication(t *testing.T) {
	client := NewFakeMarathonEndpoint(t)
	application := NewDockerApplication()
	application.ID = "/fake_app"
	app, err := client.CreateApplication(application, false)
	assert.Nil(t, err)
	assert.NotNil(t, app)
	assert.Equal(t, application.ID, "/fake_app")
	assert.Equal(t, app.DeploymentID[0]["id"], "f44fd4fc-4330-4600-a68b-99c7bd33014a")
}

func TestApplications(t *testing.T) {
	client := NewFakeMarathonEndpoint(t)
	applications, err := client.Applications(nil)
	assert.Nil(t, err)
	assert.NotNil(t, applications)
	assert.Equal(t, len(applications.Apps), 2)
}

func TestListApplications(t *testing.T) {
	client := NewFakeMarathonEndpoint(t)
	applications, err := client.ListApplications(nil)
	assert.Nil(t, err)
	assert.NotNil(t, applications)
	assert.Equal(t, len(applications), 2)
	assert.Equal(t, applications[0], FAKE_APP_NAME)
	assert.Equal(t, applications[1], FAKE_APP_NAME_BROKEN)
}

func TestApplicationVersions(t *testing.T) {
	client := NewFakeMarathonEndpoint(t)
	versions, err := client.ApplicationVersions(FAKE_APP_NAME)
	assert.Nil(t, err)
	assert.NotNil(t, versions)
	assert.NotNil(t, versions.Versions)
	assert.Equal(t, len(versions.Versions), 1)
	assert.Equal(t, versions.Versions[0], "2014-04-04T06:25:31.399Z")
	/* check we get an error on app not there */
	versions, err = client.ApplicationVersions("/not/there")
	assert.NotNil(t, err)
}

func TestSetApplicationVersion(t *testing.T) {
	client := NewFakeMarathonEndpoint(t)
	deployment, err := client.SetApplicationVersion(FAKE_APP_NAME, &ApplicationVersion{Version: "2014-08-26T07:37:50.462Z"})
	assert.Nil(t, err)
	assert.NotNil(t, deployment)
	assert.NotNil(t, deployment.Version)
	assert.NotNil(t, deployment.DeploymentID)
	assert.Equal(t, deployment.Version, "2014-08-26T07:37:50.462Z")
	assert.Equal(t, deployment.DeploymentID, "83b215a6-4e26-4e44-9333-5c385eda6438")

	_, err = client.SetApplicationVersion("/not/there", &ApplicationVersion{Version: "2014-04-04T06:25:31.399Z"})
	assert.NotNil(t, err)
}

func TestHasApplicationVersion(t *testing.T) {
	client := NewFakeMarathonEndpoint(t)
	found, err := client.HasApplicationVersion(FAKE_APP_NAME, "2014-04-04T06:25:31.399Z")
	assert.Nil(t, err)
	assert.True(t, found)
	found, err = client.HasApplicationVersion(FAKE_APP_NAME, "###2015-04-04T06:25:31.399Z")
	assert.Nil(t, err)
	assert.False(t, found)
}

func TestApplicationOK(t *testing.T) {
	client := NewFakeMarathonEndpoint(t)
	ok, err := client.ApplicationOK(FAKE_APP_NAME)
	assert.Nil(t, err)
	assert.True(t, ok)
	ok, err = client.ApplicationOK(FAKE_APP_NAME_BROKEN)
	assert.Nil(t, err)
	assert.False(t, ok)
}

func TestListApplication(t *testing.T) {
	client := NewFakeMarathonEndpoint(t)
	application, err := client.Application(FAKE_APP_NAME)
	assert.Nil(t, err)
	assert.NotNil(t, application)
	assert.Equal(t, application.ID, FAKE_APP_NAME)
	assert.NotNil(t, application.HealthChecks)
	assert.NotNil(t, application.Tasks)
	assert.Equal(t, len(application.HealthChecks), 1)
	assert.Equal(t, len(application.Tasks), 2)
}

func TestHasApplication(t *testing.T) {
	client := NewFakeMarathonEndpoint(t)
	found, err := client.HasApplication(FAKE_APP_NAME)
	assert.Nil(t, err)
	assert.True(t, found)
}
