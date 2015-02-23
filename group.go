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

type Group struct {
	ID           string         `json:"id"`
	Apps         []*Application `json:"apps"`
	Dependencies []string       `json:"dependencies"`
	Groups       []*Group       `json:"groups"`
}

type Groups struct {
	ID           string         `json:"id"`
	Apps         []*Application `json:"apps"`
	Dependencies []string       `json:"dependencies"`
	Groups       []*Group       `json:"groups"`
}

// Create a new Application Group
// Params:
//		name:	the name of the group
func NewApplicationGroup(name string) *Group {
	return &Group{
		ID:           name,
		Apps:         make([]*Application, 0),
		Dependencies: make([]string, 0),
		Groups:       make([]*Group, 0),
	}
}

// Specify the name of the group
// Params:
// 		name:	the name of the group
func (group *Group) Name(name string) *Group {
	group.ID = name
	return group
}

// Add a application to the group in question
// Params:
// 		application:	a pointer to the Application
func (group *Group) App(application *Application) *Group {
	if group.Apps == nil {
		group.Apps = make([]*Application, 0)
	}
	group.Apps = append(group.Apps, application)
	return group
}

// Retrieve a list of all the groups from marathon
func (client *Client) Groups() (*Groups, error) {
	groups := new(Groups)
	if err := client.ApiGet(MARATHON_API_GROUPS, "", groups); err != nil {
		return nil, err
	}
	return groups, nil
}

// Retrieve the configuration of a specific group from marathon
// Params:
//		name:	the identifier for the group
func (client *Client) Group(name string) (*Group, error) {
	group := new(Group)
	if err := client.ApiGet(fmt.Sprintf("%s%s", MARATHON_API_GROUPS, name), "", group); err != nil {
		return nil, err
	}
	return group, nil
}

// Check if the group exists in marathon
// Params:
// 		name:	the identifier for the group
func (client *Client) HasGroup(name string) (bool, error) {
	uri := fmt.Sprintf("%s%s", MARATHON_API_GROUPS, name)
	status, _, err := client.ApiCall(HTTP_GET, uri, "", nil)
	fmt.Printf("STATUS: %d", status)
	if err == nil {
		fmt.Printf("GHELLO")
		return true, nil
	} else if status == 404 {
		fmt.Printf("404")
		return false, nil
	} else {
		fmt.Printf("NONONON")
		return false, err
	}
}

// Create a new group in marathon
// Params:
//		group:	a pointer the Group structure defining the group
func (client *Client) CreateGroup(group *Group) (*ApplicationVersion, error) {
	version := new(ApplicationVersion)
	if err := client.ApiPost(MARATHON_API_GROUPS, group, version); err != nil {
		return nil, err
	}
	return version, nil
}

// Delete a group from marathon
// Params:
// 		name:	the identifier for the group
func (client *Client) DeleteGroup(name string) (*ApplicationVersion, error) {
	version := new(ApplicationVersion)
	uri := fmt.Sprintf("%s%s", MARATHON_API_GROUPS, name)
	if err := client.ApiDelete(uri, nil, version); err != nil {
		return nil, err
	}
	return version, nil
}
