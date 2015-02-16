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

func NewApplicationGroup(name string) *Group {
	return &Group{
		ID:           name,
		Apps:         make([]*Application, 0),
		Dependencies: make([]string, 0),
		Groups:       make([]*Group, 0),
	}
}

func (group *Group) Name(name string) *Group {
	group.ID = name
	return group
}

func (group *Group) App(application *Application) *Group {
	if group.Apps == nil {
		group.Apps = make([]*Application, 0)
	}
	group.Apps = append(group.Apps, application)
	return group
}

func (client *Client) Groups() (*Groups, error) {
	groups := new(Groups)
	if err := client.ApiGet(MARATHON_API_GROUPS, "", groups); err != nil {
		return nil, err
	}
	return groups, nil
}

func (client *Client) Group(name string) (*Group, error) {
	group := new(Group)
	if err := client.ApiGet(fmt.Sprintf("%s%s", MARATHON_API_GROUPS, name), "", group); err != nil {
		return nil, err
	}
	return group, nil
}

func (client *Client) HasGroup(name string) (bool, error) {
	if groups, err := client.Groups(); err != nil {
		return false, err
	} else {
		for _, group := range groups.Groups {
			if group.ID == name {
				return true, nil
			}
		}
		return false, nil
	}
}

func (client *Client) CreateGroup(group *Group) (*ApplicationVersion, error) {
	version := new(ApplicationVersion)
	if err := client.ApiPost(MARATHON_API_GROUPS, group, version); err != nil {
		return nil, err
	}
	return version, nil
}

func (client *Client) DeleteGroup(name string) (*ApplicationVersion, error) {
	version := new(ApplicationVersion)
	uri := fmt.Sprintf("%s%s", MARATHON_API_GROUPS, name)
	if err := client.ApiDelete(uri, "", version); err != nil {
		return nil, err
	}
	return version, nil
}
