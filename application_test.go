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
	app.DependsOn("fake_app1", "fake_app2")
	assert.Equal(t, 3, len(app.Dependencies))
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

func TestApplicationStorage(t *testing.T) {
	app := NewDockerApplication()
	assert.Equal(t, 0.0, app.Disk)
	app.Storage(0.10)
	assert.Equal(t, 0.10, app.Disk)
}

func TestApplicationName(t *testing.T) {
	app := NewDockerApplication()
	assert.Equal(t, "", app.ID)
	app.Name(fakeAppName)
	assert.Equal(t, fakeAppName, app.ID)
}

func TestApplicationCPU(t *testing.T) {
	app := NewDockerApplication()
	assert.Equal(t, 0.0, app.CPUs)
	app.CPU(0.1)
	assert.Equal(t, 0.1, app.CPUs)
}

func TestApplicationArgs(t *testing.T) {
	app := NewDockerApplication()
	assert.Equal(t, 0, len(app.Args))
	app.Arg("-p").Arg("option", "-v")
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

func TestApplicationLabels(t *testing.T) {
	app := NewDockerApplication()
	assert.Equal(t, 0, len(app.Labels))
	app.AddLabel("hello", "world")
	assert.Equal(t, 1, len(app.Labels))
}

func TestHasHealthChecks(t *testing.T) {
	app := NewDockerApplication()
	assert.False(t, app.HasHealthChecks())
	app.Container.Docker.Container("quay.io/gambol99/apache-php:latest").Expose(80)
	_, err := app.CheckTCP(80, 10)
	assert.NoError(t, err)
	assert.True(t, app.HasHealthChecks())
}

func TestApplicationCheckTCP(t *testing.T) {
	app := NewDockerApplication()
	assert.False(t, app.HasHealthChecks())
	_, err := app.CheckTCP(80, 10)
	assert.Error(t, err)
	assert.False(t, app.HasHealthChecks())
	app.Container.Docker.Container("quay.io/gambol99/apache-php:latest").Expose(80)
	_, err = app.CheckTCP(80, 10)
	assert.NoError(t, err)
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
	assert.Error(t, err)
	assert.False(t, app.HasHealthChecks())
	app.Container.Docker.Container("quay.io/gambol99/apache-php:latest").Expose(80)
	_, err = app.CheckHTTP("/health", 80, 10)
	assert.NoError(t, err)
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
	endpoint := newFakeMarathonEndpoint(t, nil)
	defer endpoint.Close()

	application := NewDockerApplication()
	application.ID = "/fake_app"
	app, err := endpoint.Client.CreateApplication(application)
	assert.NoError(t, err)
	assert.NotNil(t, app)
	assert.Equal(t, application.ID, "/fake_app")
	assert.Equal(t, app.Deployments[0]["id"], "f44fd4fc-4330-4600-a68b-99c7bd33014a")
}

func TestUpdateApplication(t *testing.T) {
	endpoint := newFakeMarathonEndpoint(t, nil)
	defer endpoint.Close()

	application := NewDockerApplication()
	application.ID = "/fake_app"
	id, err := endpoint.Client.UpdateApplication(application)
	assert.NoError(t, err)
	assert.Equal(t, id.DeploymentID, "83b215a6-4e26-4e44-9333-5c385eda6438")
	assert.Equal(t, id.Version, "2014-08-26T07:37:50.462Z")
}

func TestApplications(t *testing.T) {
	endpoint := newFakeMarathonEndpoint(t, nil)
	defer endpoint.Close()

	applications, err := endpoint.Client.Applications(nil)
	assert.NoError(t, err)
	assert.NotNil(t, applications)
	assert.Equal(t, len(applications.Apps), 2)
}

func TestListApplications(t *testing.T) {
	endpoint := newFakeMarathonEndpoint(t, nil)
	defer endpoint.Close()

	applications, err := endpoint.Client.ListApplications(nil)
	assert.NoError(t, err)
	assert.NotNil(t, applications)
	assert.Equal(t, len(applications), 2)
	assert.Equal(t, applications[0], fakeAppName)
	assert.Equal(t, applications[1], fakeAppNameBroken)
}

func TestApplicationVersions(t *testing.T) {
	endpoint := newFakeMarathonEndpoint(t, nil)
	defer endpoint.Close()

	versions, err := endpoint.Client.ApplicationVersions(fakeAppName)
	assert.NoError(t, err)
	assert.NotNil(t, versions)
	assert.NotNil(t, versions.Versions)
	assert.Equal(t, len(versions.Versions), 1)
	assert.Equal(t, versions.Versions[0], "2014-04-04T06:25:31.399Z")
	/* check we get an error on app not there */
	versions, err = endpoint.Client.ApplicationVersions("/not/there")
	assert.Error(t, err)
}

func TestRestartApplication(t *testing.T) {
	endpoint := newFakeMarathonEndpoint(t, nil)
	defer endpoint.Close()

	id, err := endpoint.Client.RestartApplication(fakeAppName, false)
	assert.NoError(t, err)
	assert.NotNil(t, id)
	assert.Equal(t, "83b215a6-4e26-4e44-9333-5c385eda6438", id.DeploymentID)
	assert.Equal(t, "2014-08-26T07:37:50.462Z", id.Version)
	id, err = endpoint.Client.RestartApplication("/not/there", false)
	assert.Error(t, err)
	assert.Nil(t, id)
}

func TestSetApplicationVersion(t *testing.T) {
	endpoint := newFakeMarathonEndpoint(t, nil)
	defer endpoint.Close()

	deployment, err := endpoint.Client.SetApplicationVersion(fakeAppName, &ApplicationVersion{Version: "2014-08-26T07:37:50.462Z"})
	assert.NoError(t, err)
	assert.NotNil(t, deployment)
	assert.NotNil(t, deployment.Version)
	assert.NotNil(t, deployment.DeploymentID)
	assert.Equal(t, deployment.Version, "2014-08-26T07:37:50.462Z")
	assert.Equal(t, deployment.DeploymentID, "83b215a6-4e26-4e44-9333-5c385eda6438")
	_, err = endpoint.Client.SetApplicationVersion("/not/there", &ApplicationVersion{Version: "2014-04-04T06:25:31.399Z"})
	assert.Error(t, err)
}

func TestHasApplicationVersion(t *testing.T) {
	endpoint := newFakeMarathonEndpoint(t, nil)
	defer endpoint.Close()

	found, err := endpoint.Client.HasApplicationVersion(fakeAppName, "2014-04-04T06:25:31.399Z")
	assert.NoError(t, err)
	assert.True(t, found)
	found, err = endpoint.Client.HasApplicationVersion(fakeAppName, "###2015-04-04T06:25:31.399Z")
	assert.NoError(t, err)
	assert.False(t, found)
}

func TestDeleteApplication(t *testing.T) {
	endpoint := newFakeMarathonEndpoint(t, nil)
	defer endpoint.Close()

	id, err := endpoint.Client.DeleteApplication(fakeAppName)
	assert.NoError(t, err)
	assert.NotNil(t, id)
	assert.Equal(t, "83b215a6-4e26-4e44-9333-5c385eda6438", id.DeploymentID)
	assert.Equal(t, "2014-08-26T07:37:50.462Z", id.Version)
	id, err = endpoint.Client.DeleteApplication("no_such_app")
	assert.Error(t, err)
}

func TestApplicationOK(t *testing.T) {
	endpoint := newFakeMarathonEndpoint(t, nil)
	defer endpoint.Close()

	ok, err := endpoint.Client.ApplicationOK(fakeAppName)
	assert.NoError(t, err)
	assert.True(t, ok)
	ok, err = endpoint.Client.ApplicationOK(fakeAppNameBroken)
	assert.NoError(t, err)
	assert.False(t, ok)
}

func TestApplication(t *testing.T) {
	endpoint := newFakeMarathonEndpoint(t, nil)
	defer endpoint.Close()

	application, err := endpoint.Client.Application(fakeAppName)
	assert.NoError(t, err)
	assert.NotNil(t, application)
	assert.Equal(t, application.ID, fakeAppName)
	assert.NotNil(t, application.HealthChecks)
	assert.NotNil(t, application.Tasks)
	assert.Equal(t, len(application.HealthChecks), 1)
	assert.Equal(t, len(application.Tasks), 2)

	_, err = endpoint.Client.Application("no_such_app")
	assert.Equal(t, ErrDoesNotExist, err)

	config := NewDefaultConfig()
	config.URL = "http://non-existing-marathon-host.local:5555"
	endpoint = newFakeMarathonEndpoint(t, &config)
	defer endpoint.Close()

	_, err = endpoint.Client.Application(fakeAppName)
	assert.NotEqual(t, ErrDoesNotExist, err)
	assert.Error(t, err)
}
