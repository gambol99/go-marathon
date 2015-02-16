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

package main

import (
	"flag"

	marathon "github.com/gambol99/go-marathon"
	"github.com/golang/glog"
)

var marathon_url string

func init() {
	flag.StringVar(&marathon_url, "url", "http://127.0.0.1:8080", "the url for the marathon endpoint")
}

func Assert(err error) {
	if err != nil {
		glog.Fatalf("Failed, error: %s", err)
	}
}

func main() {
	flag.Parse()
	config := marathon.NewDefaultConfig()
	config.URL = marathon_url
	config.Debug = true
	client, err := marathon.NewClient(config)
	if err != nil {
		glog.Fatalf("Failed to create a client for marathon, error: %s", err)
	}

	glog.Infof("Retrieving a list of groups")
	if groups, err := client.Groups(); err != nil {
		glog.Errorf("Failed to retrieve the groups from maratho, error: %s", err)
	} else {
		for _, group := range groups.Groups {
			glog.Infof("Found group: %s", group.ID)
		}
	}

	/* step: the frontend app */
	frontend := new(marathon.Application)
	frontend.Name("/product/group/frontend")
	frontend.CPU(0.1).Memory(64).Storage(0.0).Count(2)
	frontend.Arg("/usr/sbin/apache2ctl").Arg("-D").Arg("FOREGROUND")
	frontend.AddEnv("NAME", "frontend_http")
	frontend.AddEnv("SERVICE_80_NAME", "test_http")
	frontend.Container = marathon.NewDockerContainer()
	frontend.DependsOn("/product/group/cache")
	frontend.Container.Docker.Container("quay.io/gambol99/apache-php:latest").Expose(80).Expose(443)

	redis := new(marathon.Application)
	redis.Name("/product/group/cache")
	redis.CPU(0.1).Memory(64).Storage(0.0).Count(2)
	redis.Arg("/usr/sbin/apache2ctl").Arg("-D").Arg("FOREGROUND")
	redis.AddEnv("NAME", "group_cache")
	redis.AddEnv("SERVICE_6379_NAME", "test_redis")
	redis.Container = marathon.NewDockerContainer()
	redis.Container.Docker.Container("redis:latest").Expose(6379)

	group := marathon.NewApplicationGroup("/product/group")
	group.App(frontend).App(redis)

	if version, err := client.CreateGroup(group); err != nil {
		glog.Errorf("Failed to create the group: %s, error: %s", group.ID, err)
	} else {
		glog.Infof("Successfully created the group: %s, version: %s", group.ID, version)
	}
}
