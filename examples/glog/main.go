/*
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

var marathonURL string

func init() {
	glog.CopyStandardLogTo("INFO") // redirect marathon logs to glog

	flag.StringVar(&marathonURL, "url", "http://127.0.0.1:8080", "the url for the marathon endpoint")
}

func main() {
	flag.Parse()
	config := marathon.NewDefaultConfig()
	config.URL = marathonURL
	client, err := marathon.NewClient(config)
	if err != nil {
		glog.Exitln(err)
	}

	applications, err := client.Applications(nil)
	if err != nil {
		glog.Exitln(err)
	}

	for _, a := range applications.Apps {
		glog.Infof("App ID: %v\n", a.ID)
	}
}
