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
)

func TestApplications(t *testing.T) {
	client := NewFakeMarathonEndpoint(t)
	applications, err := client.Applications()
	assertOnError(err, t)
	assertOnNull(applications, t)
	assertOnInteger(len(applications.Apps), 2, t)
}

func TestListApplications(t *testing.T) {
	client := NewFakeMarathonEndpoint(t)
	applications, err := client.ListApplications()
	assertOnError(err, t)
	assertOnNull(applications, t)
	assertOnInteger(len(applications), 2, t)
	assertOnString(applications[0], FAKE_APP_NAME, t)
	assertOnString(applications[1], FAKE_APP_NAME_BROKEN, t)
}

func TestApplicationVersions(t *testing.T) {
	client := NewFakeMarathonEndpoint(t)
	versions, err := client.ApplicationVersions(FAKE_APP_NAME)
	assertOnError(err, t)
	assertOnNull(versions, t)
	assertOnNull(versions.Versions, t)
	assertOnInteger(len(versions.Versions), 1, t)
	assertOnString(versions.Versions[0], "2014-04-04T06:25:31.399Z", t)
	/* check we get an error on app not there */
	versions, err = client.ApplicationVersions("/not/there")
	assertOnNoError(err, t)
}

func TestSetApplicationVersion(t *testing.T) {
	client := NewFakeMarathonEndpoint(t)
	deployment, err := client.SetApplicationVersion(FAKE_APP_NAME, &ApplicationVersion{Version: "2014-04-04T06:25:31.399Z"})
	assertOnError(err, t)
	assertOnNull(deployment, t)
	assertOnNull(deployment.Version, t)
	assertOnNull(deployment.DeploymentID, t)
	assertOnString(deployment.Version, "2014-04-04T06:25:31.399Z", t)
	assertOnString(deployment.DeploymentID, "83b215a6-4e26-4e44-9333-5c385eda6438", t)

	_, err = client.SetApplicationVersion("/not/there", &ApplicationVersion{Version: "2014-04-04T06:25:31.399Z"})
	assertOnNoError(err, t)
}

func TestHasApplicationVersion(t *testing.T) {
	client := NewFakeMarathonEndpoint(t)
	found, err := client.HasApplicationVersion(FAKE_APP_NAME, "2014-04-04T06:25:31.399Z")
	assertOnError(err, t)
	assertOnBool(found, true, t)
	found, err = client.HasApplicationVersion(FAKE_APP_NAME, "###2015-04-04T06:25:31.399Z")
	assertOnError(err, t)
	assertOnBool(found, false, t)
}

func TestApplicationOK(t *testing.T) {
	client := NewFakeMarathonEndpoint(t)
	ok, err := client.ApplicationOK(FAKE_APP_NAME)
	assertOnError(err, t)
	assertOnBool(ok, true, t)
	ok, err = client.ApplicationOK(FAKE_APP_NAME_BROKEN)
	assertOnError(err, t)
	assertOnBool(ok, false, t)
}

func TestListApplication(t *testing.T) {
	client := NewFakeMarathonEndpoint(t)
	application, err := client.Application(FAKE_APP_NAME)
	assertOnError(err, t)
	assertOnNull(application, t)
	assertOnString(application.ID, FAKE_APP_NAME, t)
	assertOnNull(application.HealthChecks, t)
	assertOnNull(application.Tasks, t)
	assertOnInteger(len(application.HealthChecks), 1, t)
	assertOnInteger(len(application.Tasks), 2, t)
}

func TestHasApplication(t *testing.T) {
	client := NewFakeMarathonEndpoint(t)
	found, err := client.HasApplication(FAKE_APP_NAME)
	assertOnError(err, t)
	assertOnBool(found, true, t)
}
