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

type Applications struct {
	Apps []Application `json:"apps"`
}

func (client *MarathonClient) Applications() (Applications, error) {
	var apps Applications
	if err := client.ApiGet(MARATHON_API_APPS, &apps); err != nil {
		return Applications{}, err
	} else {
		return apps, nil
	}
}

func (client *MarathonClient) ListApplications() ([]string, error) {
	if applications, err := client.Applications(); err != nil {
		return nil, err
	} else {
		list := make([]string, 0)
		for _, application := range applications.Apps {
			list = append(list, application.ID)
		}
		return list, nil
	}
}
