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

func TestDeployments(t *testing.T) {
	NewFakeMarathonEndpoint()
	deployments, err := test_client.Deployments()
	AssertOnError(err, t)
	AssertOnNull(deployments, t)
	AssertOnInteger(len(deployments), 1, t)
	deployment := deployments[0]
	AssertOnNull(deployment, t)
	AssertOnString(deployment.ID, "867ed450-f6a8-4d33-9b0e-e11c5513990b", t)
	AssertOnNull(deployment.Steps, t)
	AssertOnInteger(len(deployment.Steps), 1, t)
}

func TestDeleteDeployment(t *testing.T) {
	NewFakeMarathonEndpoint()
	id, err := test_client.DeleteDeployment(FAKE_DEPLOYMENT_ID, false)
	AssertOnError(err, t)
	AssertOnNull(id, t)
	AssertOnString(id.DeploymentID, "0b1467fc-d5cd-4bbc-bac2-2805351cee1e", t)
	AssertOnString(id.Version, "2014-08-26T08:20:26.171Z", t)
}


