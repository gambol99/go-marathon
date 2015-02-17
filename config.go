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

type Config struct {
	/* the url for marathon */
	URL string
	/* event handler port */
	EventsPort int
	/* the interface we should be listening on for events */
	EventsInterface string
	/* switch on debugging */
	Debug bool
	/* the timeout for requests */
	RequestTimeout int
}

func NewDefaultConfig() Config {
	return Config{
		URL:             "http://localhost:8080",
		EventsPort:      DEFAULT_EVENTS_PORT,
		EventsInterface: DEFAULT_EVENTS_BIND,
		Debug:           false,
		RequestTimeout:  5}
}
