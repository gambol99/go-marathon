/*
Copyright 2017 The go-marathon Authors All rights reserved.

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
	"fmt"
	"log"
	"time"

	marathon "github.com/gambol99/go-marathon"
)

var marathonURL string
var dcosToken string

func init() {
	flag.StringVar(&marathonURL, "url", "http://127.0.0.1:8080", "the url for the marathon endpoint")
	flag.StringVar(&dcosToken, "token", "", "DCOS token for auth")
}

func assert(err error) {
	if err != nil {
		log.Fatalf("Failed, error: %s", err)
	}
}

func waitOnDeployment(client marathon.Marathon, id *marathon.DeploymentID) {
	assert(client.WaitOnDeployment(id.DeploymentID, 0))
}

func createRawPod() *marathon.Pod {
	var containers []*marathon.PodContainer
	for i := 0; i < 2; i++ {
		container := &marathon.PodContainer{
			Name: fmt.Sprintf("container%d", i),
			Exec: &marathon.PodExec{
				Command: marathon.PodCommand{
					Shell: "echo Hello World && sleep 600",
				},
			},
			Image: &marathon.PodContainerImage{
				Kind:      "DOCKER",
				ID:        "nginx",
				ForcePull: true,
			},
			VolumeMounts: []*marathon.PodVolumeMount{
				&marathon.PodVolumeMount{
					Name:      "sharedvolume",
					MountPath: "/peers",
				},
			},
			Resources: &marathon.Resources{
				Cpus: 0.1,
				Mem:  128,
			},
			Env: map[string]string{
				"key": "value",
			},
		}

		containers = append(containers, container)
	}

	pod := &marathon.Pod{
		ID:         "/mypod",
		Containers: containers,
		Scaling: &marathon.PodScalingPolicy{
			Kind:      "fixed",
			Instances: 2,
		},
		Volumes: []*marathon.PodVolume{
			&marathon.PodVolume{
				Name: "sharedvolume",
				Host: "/tmp",
			},
		},
	}

	return pod
}

func createConveniencePod() *marathon.Pod {
	pod := marathon.NewPod()

	pod.Name("mypod").
		Count(2).
		AddVolume(marathon.NewPodVolume("sharedvolume", "/tmp"))

	for i := 0; i < 2; i++ {
		image := marathon.NewDockerPodContainerImage().
			SetID("nginx")

		container := marathon.NewPodContainer().
			SetName(fmt.Sprintf("container%d", i)).
			CPUs(0.1).
			Memory(128).
			SetImage(image).
			AddEnv("key", "value").
			AddVolumeMount(marathon.NewPodVolumeMount("sharedvolume", "/peers")).
			SetCommand("echo Hello World && sleep 600")

		pod.AddContainer(container)
	}

	return pod
}

func doPlayground(client marathon.Marathon, pod *marathon.Pod) {
	// Create a pod
	fmt.Println("Creating pod...")
	pod, err := client.CreatePod(pod)
	assert(err)

	// Check its status
	fmt.Println("Waiting on pod...")
	err = client.WaitOnPod(pod.ID, time.Minute*1)
	assert(err)

	// Scale it
	fmt.Println("Scaling pod...")
	pod.Count(5)
	pod, err = client.UpdatePod(pod, true)
	assert(err)

	// Get instances
	status, err := client.PodStatus(pod.ID)
	fmt.Printf("Pod status: %s\n", status.Status)
	assert(err)

	// Kill an instance
	fmt.Println("Deleting an instance...")
	_, err = client.DeletePodInstance(pod.ID, status.Instances[0].ID)
	assert(err)

	// Delete it
	fmt.Println("Deleting the pod")
	id, err := client.DeletePod(pod.ID, true)
	assert(err)

	waitOnDeployment(client, id)
}

func main() {
	flag.Parse()
	config := marathon.NewDefaultConfig()
	config.URL = marathonURL
	config.DCOSToken = dcosToken
	client, err := marathon.NewClient(config)
	assert(err)

	fmt.Println("Convenience Pods:")
	podC := createConveniencePod()
	doPlayground(client, podC)

	fmt.Println("Raw Pods:")
	podR := createRawPod()
	doPlayground(client, podR)
}
