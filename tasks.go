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
	"encoding/json"
)

type Tasks struct {
	Tasks []Task `json:"tasks"`
}

var (

)

func (client *MarathonClient) AllTasks() (Tasks, error) {
	var tasks Tasks
	if tasks, _, err := client.ApiGet(MARATHON_API_TASKS, &tasks); err != nil {
		return nil, err
	} else {
		return tasks, nil
	}
}

func (r *MarathonClient) Tasks(application_id string) (tasks Tasks, err error) {
	var response string
	if response, err = r.Get(fmt.Sprintf("%s%s/tasks", MARATHON_API_APPS, application_id ) ); err != nil {
		return
	} else {
		/* step: we need to un-marshall the json response from marathon */
		if err = json.Unmarshal([]byte(response), &tasks); err != nil {
			return
		}
		return
	}
}

