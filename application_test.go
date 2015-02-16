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
	NewFakeMarathonEndpoint()
	applications, err := test_client.Applications()
	AssertOnError(err, t)
	AssertOnNull(applications, t)
	AssertOnInteger(len(applications.Apps), 2, t)
}

func TestListApplications(t *testing.T) {
	NewFakeMarathonEndpoint()
	applications, err := test_client.ListApplications()
	AssertOnError(err, t)
	AssertOnNull(applications, t)
	AssertOnInteger(len(applications), 2, t)
	AssertOnString(applications[0], FAKE_APP_NAME, t)
	AssertOnString(applications[1], FAKE_APP_NAME_BROKEN, t)
}

func TestApplicationVersions(t *testing.T) {
	NewFakeMarathonEndpoint()
	versions, err := test_client.ApplicationVersions(FAKE_APP_NAME)
	AssertOnError(err, t)
	AssertOnNull(versions, t)
	AssertOnNull(versions.Versions, t)
	AssertOnInteger(len(versions.Versions), 1, t)
	AssertOnString(versions.Versions[0], "2014-04-04T06:25:31.399Z", t)
	/* check we get an error on app not there */
	versions, err = test_client.ApplicationVersions("/not/there")
	AssertOnNoError(err, t)
}

func TestSetApplicationVersion(t *testing.T) {
	NewFakeMarathonEndpoint()
	deployment, err := test_client.SetApplicationVersion(FAKE_APP_NAME, &ApplicationVersion{Version: "2014-04-04T06:25:31.399Z"})
	AssertOnError(err, t)
	AssertOnNull(deployment, t)
	AssertOnNull(deployment.Version, t)
	AssertOnNull(deployment.DeploymentID, t)
	AssertOnString(deployment.Version, "2014-04-04T06:25:31.399Z", t)
	AssertOnString(deployment.DeploymentID, "83b215a6-4e26-4e44-9333-5c385eda6438", t)

	_, err = test_client.SetApplicationVersion("/not/there", &ApplicationVersion{Version: "2014-04-04T06:25:31.399Z"})
	AssertOnNoError(err, t)
}

func TestHasApplicationVersion(t *testing.T) {
	NewFakeMarathonEndpoint()
	found, err := test_client.HasApplicationVersion(FAKE_APP_NAME, "2014-04-04T06:25:31.399Z")
	AssertOnError(err, t)
	AssertOnBool(found, true, t)
	found, err = test_client.HasApplicationVersion(FAKE_APP_NAME, "###2015-04-04T06:25:31.399Z")
	AssertOnError(err, t)
	AssertOnBool(found, false, t)
}

func TestApplicationOK(t *testing.T) {
	NewFakeMarathonEndpoint()
	ok, err := test_client.ApplicationOK(FAKE_APP_NAME)
	AssertOnError(err, t)
	AssertOnBool(ok, true, t)
	ok, err = test_client.ApplicationOK(FAKE_APP_NAME_BROKEN)
	AssertOnError(err, t)
	AssertOnBool(ok, false, t)
}

func TestListApplication(t *testing.T) {
	NewFakeMarathonEndpoint()
	application, err := test_client.Application(FAKE_APP_NAME)
	AssertOnError(err, t)
	AssertOnNull(application, t)
	AssertOnString(application.ID, FAKE_APP_NAME, t)
	AssertOnNull(application.HealthChecks, t)
	AssertOnNull(application.Tasks, t)
	AssertOnInteger(len(application.HealthChecks), 1, t)
	AssertOnInteger(len(application.Tasks), 2, t)
}

func TestHasApplication(t *testing.T) {
	NewFakeMarathonEndpoint()
	found, err := test_client.HasApplication(FAKE_APP_NAME)
	AssertOnError(err, t)
	AssertOnBool(found, true, t)
}
