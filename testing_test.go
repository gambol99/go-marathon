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
	"io/ioutil"
	"net/http"
	"sync"
	"testing"

	yaml "gopkg.in/yaml.v2"
)

const (
	fakeMarathonURL   = "http://127.0.0.1:3000,127.0.0.1:3000,127.0.0.1:3000"
	fakeGroupName     = "/test"
	fakeGroupName1    = "/qa/product/1"
	fakeAppName       = "/fake_app"
	fakeAppNameBroken = "/fake_app_broken"
	fakeDeploymentID  = "867ed450-f6a8-4d33-9b0e-e11c5513990b"
	fakeAPIFilename   = "./tests/rest-api/methods.yml"
	fakeAPIPort       = 3000
)

type RestMethod struct {
	// the uri of the method
	URI string `yaml:"uri,omitempty"`
	// the http method type (GET|PUT etc)
	Method string `yaml:"method,omitempty"`
	// the content i.e. response
	Content string `yaml:"content,omitempty"`
}

var testClient struct {
	sync.Once
	client Marathon
}

func NewFakeMarathonEndpoint(t *testing.T) Marathon {
	testClient.Once.Do(func() {

		// step: open and read in the methods yaml
		contents, err := ioutil.ReadFile(fakeAPIFilename)
		if err != nil {
			t.Fatalf("unable to read in the methods yaml file: %s", fakeAPIFilename)
		}
		// step: unmarshal the yaml
		var methods []*RestMethod
		err = yaml.Unmarshal([]byte(contents), &methods)
		if err != nil {
			t.Fatalf("Unable to unmarshall the methods yaml, error: %s", err)
		}

		// step: construct a hash from the methods
		uris := make(map[string]*string, 0)
		for _, method := range methods {
			uris[fmt.Sprintf("%s:%s", method.Method, method.URI)] = &method.Content
		}

		http.HandleFunc("/", func(writer http.ResponseWriter, reader *http.Request) {
			key := fmt.Sprintf("%s:%s", reader.Method, reader.RequestURI)
			content, found := uris[key]
			if found {
				writer.Header().Add("Content-Type", "application/json")
				writer.Write([]byte(*content))
			}
			if !found {
				http.Error(writer, "not found", 404)
			}
		})

		go http.ListenAndServe(fmt.Sprintf("127.0.0.1:%d", fakeAPIPort), nil)

		config := NewDefaultConfig()
		config.URL = fakeMarathonURL
		//config.LogOutput = os.Stdout
		if testClient.client, err = NewClient(config); err != nil {
			t.Fatalf("Failed to create the fake client, %s, error: %s", fakeMarathonURL, err)
		}
	})

	return testClient.client
}
