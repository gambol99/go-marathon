[![Build Status](https://travis-ci.org/gambol99/go-marathon.svg?branch=master)](https://travis-ci.org/gambol99/go-marathon)
[![GoDoc](http://godoc.org/github.com/gambol99/go-marathon?status.png)](http://godoc.org/github.com/gambol99/go-marathon)

#### **Go-Marathon**
-----

Go-marathon is a api library for working with [Marathon](https://mesosphere.github.io/marathon/). It currently supports 

  > - Application and group deployment
  > - Helper filters for pulling the status, configuration and tasks
  > - Multiple Endpoint support for HA deployments
  > - Marathon Subscriptions and Event callbacks
  
> ##### **Examples**

Check out the examples directory for more code examples

> ##### **Creating a client**

    import (
    	"flag"
    
    	marathon "github.com/gambol99/go-marathon"
    	"github.com/golang/glog"
    	"time"
    )
  
    marathon_url := http://10.241.1.71:8080
  	config := marathon.NewDefaultConfig()
  	config.URL = marathon_url
  	config.Debug = true
  	if client, err := marathon.NewClient(config); err != nil {
  		glog.Fatalf("Failed to create a client for marathon, error: %s", err)
  	} else {
  		applications, err := client.Applications()
  		...

> Note, you can also specify multiple endpoint for Marathon (i.e. you have setup Marathon in HA mode and having multiple running)

	marathon := "http://10.241.1.71:8080,10.241.1.72:8080,10.241.1.73:8080"
	
The first one specified will be used, if that goes offline the member is marked as *"unavailable"* and a background process will continue to ping the member until it's back online.

> ##### **Listing the applications**

	if applications, err := client.Applications(); err != nil 
		glog.Errorf("Failed to list applications")
	} else {
		glog.Infof("Found %d application running", len(applications.Apps))
		for _, application := range applications.Apps {
			glog.Infof("Application: %s", application)
			details, err := client.Application(application.ID)
			Assert(err)
			if details.Tasks != nil && len(details.Tasks) > 0 {
				for _, task := range details.Tasks {
					glog.Infof("task: %s", task)
				}
				// check the health of the application
				health, err := client.ApplicationOK(details.ID)
				glog.Infof("Application: %s, healthy: %t", details.ID, health)
			}
		}	


> ##### **Creating a new application**

	
		glog.Infof("Deploying a new application")
		application := new(marathon.Application)
		application.Name("/product/name/frontend)
		application.CPU(0.1).Memory(64).Storage(0.0).Count(2)
		application.Arg("/usr/sbin/apache2ctl").Arg("-D").Arg("FOREGROUND")
		application.AddEnv("NAME", "frontend_http")
		application.AddEnv("SERVICE_80_NAME", "test_http")
		// add the docker container
		application.Container = marathon.NewDockerContainer()
		application.Container.Docker.Container("quay.io/gambol99/apache-php:latest").Expose(80).Expose(443)

		if err := client.CreateApplication(application); err != nil {
			glog.Errorf("Failed to create application: %s, error: %s", application, err)
		} else {
			glog.Infof("Created the application: %s", application)
		}
		
> ##### **Scale Application**

Change the number of instance of the application to 4 

    glog.Infof("Scale to 4 instances")
		if err := client.ScaleApplicationInstances(application.ID, 4); err != nil {
			glog.Errorf("Failed to delete the application: %s, error: %s", application, err)
		} else {
			glog.Infof("Successfully scaled the application")
		}
	
