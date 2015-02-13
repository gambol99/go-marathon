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

type Event struct {
	EventType string `json:"eventType"`
	AppID     string `json:"appId"`
}

type EventStatusUpdate struct {
	EventType  string `json:"eventType"`
	Timestamp  string `json:"timestamp,omitempty"`
	SlaveID    string `json:"slaveId,omitempty"`
	TaskID     string `json:"taskId"`
	TaskStatus string `json:"taskStatus"`
	AppID      string `json:"appId"`
	Host       string `json:"host"`
	Ports      []int  `json:"ports,omitempty"`
	Version    string `json:"version,omitempty"`
}

type EventHealthCheckChanged struct {
	EventType string `json:"eventType"`
	Timestamp string `json:"timestamp,omitempty"`
	AppID     string `json:"appId"`
	TaskID    string `json:"taskId"`
	Version   string `json:"version,omitempty"`
	Alive     bool   `json:"alive"`
}
