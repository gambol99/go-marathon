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
	"io"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"net/url"
	"strings"
	"sync"
	"testing"

	"github.com/donovanhide/eventsource"
	yaml "gopkg.in/yaml.v2"
)

const (
	fakeMarathonURL   = "http://127.0.0.1:3000,127.0.0.1:3000,127.0.0.1:3000"
	fakeGroupName     = "/test"
	fakeGroupName1    = "/qa/product/1"
	fakeAppName       = "/fake_app"
	fakeTaskID        = "fake_app.12345"
	fakeAppNameBroken = "/fake_app_broken"
	fakeDeploymentID  = "867ed450-f6a8-4d33-9b0e-e11c5513990b"
	fakeAPIFilename   = "./tests/rest-api/methods.yml"
	fakeAPIPort       = 3000
)

type restMethod struct {
	// the uri of the method
	URI string `yaml:"uri,omitempty"`
	// the http method type (GET|PUT etc)
	Method string `yaml:"method,omitempty"`
	// the content i.e. response
	Content string `yaml:"content,omitempty"`
}

type fakeServer struct {
	io.Closer

	eventSrv *eventsource.Server
	httpSrv  *httptest.Server
}

type endpoint struct {
	io.Closer

	Server fakeServer
	Client Marathon
	URL    string
}

type fakeEvent struct {
	data string
}

var uris map[string]*string
var once sync.Once

func getTestURL(urlString string) string {
	parsedURL, _ := url.Parse(urlString)
	return fmt.Sprintf("%s://%s", parsedURL.Scheme, strings.Join([]string{parsedURL.Host, parsedURL.Host, parsedURL.Host}, ","))
}

func newFakeMarathonEndpoint(t *testing.T, config *Config) *endpoint {
	once.Do(func() {
		// step: open and read in the methods yaml
		contents, err := ioutil.ReadFile(fakeAPIFilename)
		if err != nil {
			t.Fatalf("unable to read in the methods yaml file: %s", fakeAPIFilename)
		}
		// step: unmarshal the yaml
		var methods []*restMethod
		err = yaml.Unmarshal([]byte(contents), &methods)
		if err != nil {
			t.Fatalf("Unable to unmarshal the methods yaml, error: %s", err)
		}

		// step: construct a hash from the methods
		uris = make(map[string]*string, 0)
		for _, method := range methods {
			uris[fmt.Sprintf("%s:%s", method.Method, method.URI)] = &method.Content
		}
	})

	eventSrv := eventsource.NewServer()

	mux := http.NewServeMux()
	mux.HandleFunc("/v2/events", eventSrv.Handler("event"))
	mux.HandleFunc("/", func(writer http.ResponseWriter, reader *http.Request) {
		key := fmt.Sprintf("%s:%s", reader.Method, reader.RequestURI)
		content, found := uris[key]
		if found {
			writer.Header().Add("Content-Type", "application/json")
			writer.Write([]byte(*content))
			return
		}

		http.Error(writer, "not found", 404)
	})

	httpSrv := httptest.NewServer(mux)

	defaultConfig := NewDefaultConfig()
	if config == nil {
		config = &defaultConfig
	}

	if config.URL == defaultConfig.URL {
		config.URL = getTestURL(httpSrv.URL)
	}

	client, err := NewClient(*config)
	if err != nil {
		t.Fatalf("Failed to create the fake client, %s, error: %s", config.URL, err)
	}

	return &endpoint{
		Server: fakeServer{
			eventSrv: eventSrv,
			httpSrv:  httpSrv,
		},
		Client: client,
		URL:    config.URL,
	}
}

func (t fakeEvent) Id() string {
	return "0"
}

func (t fakeEvent) Event() string {
	return "MarathonEvent"
}

func (t fakeEvent) Data() string {
	return t.data
}

func (s *fakeServer) PublishEvent(event string) {
	s.eventSrv.Publish([]string{"event"}, fakeEvent{event})
}

func (s *fakeServer) Close() error {
	s.eventSrv.Close()
	s.httpSrv.Close()
	return nil
}

func (e *endpoint) Close() error {
	return e.Server.Close()
}
