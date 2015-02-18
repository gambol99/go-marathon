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
var marathon_interface string
var marathon_port int

func init() {
	flag.StringVar(&marathon_url, "url", "http://127.0.0.1:8080", "the url for the marathon endpoint")
	flag.StringVar(&marathon_interface, "interface", "eth0", "the interface we should use for events")
	flag.IntVar(&marathon_port, "port", 10001, "the port the events service should run on")
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
	config.Debug = false
	config.EventsPort = marathon_port
	config.EventsInterface = marathon_interface
	client, err := marathon.NewClient(config)
	if err != nil {
		glog.Fatalf("Failed to create a client for marathon, error: %s", err)
	}

	/* step: lets register for events */
	update := make(marathon.EventsChannel, 5)
	if err := client.AddEventsListener(update, marathon.EVENTS_APPLICATIONS); err != nil {
		glog.Fatalf("Failed to register for subscriptions, %s", err)
	} else {
		for {
			event := <-update
			glog.Infof("EVENT: %s", event)
		}
	}
}
