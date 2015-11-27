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
	"io"
	"io/ioutil"
)

// EventsTransport describes which transport should be used to deliver Marathon events
type EventsTransport int

// Config holds the settings and options for the client
type Config struct {
	// the url for marathon
	URL string
	// events transport: EventsTransportCallback or EventsTransportSSE
	EventsTransport EventsTransport
	// event handler port
	EventsPort int
	// the interface we should be listening on for events
	EventsInterface string
	// the timeout for requests
	RequestTimeout int
	// http basic auth
	HTTPBasicAuthUser string
	// http basic password
	HTTPBasicPassword string
	// custom callback url
	CallbackURL string
	// the output for debug log messages
	LogOutput io.Writer
}

// NewDefaultConfig create a default client config
func NewDefaultConfig() Config {
	return Config{
		URL:             "http://127.0.0.1:8080",
		EventsTransport: EventsTransportCallback,
		EventsPort:      10001,
		EventsInterface: "eth0",
		RequestTimeout:  5,
		LogOutput:       ioutil.Discard,
	}
}
