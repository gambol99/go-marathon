#
#   Author: Rohith (gambol99@gmail.com)
#   Date: 2015-02-10 15:35:14 +0000 (Tue, 10 Feb 2015)
#
#  vim:ts=2:sw=2:et
#
NAME="go-marathon"
AUTHOR=gambol99
HARDWARE=$(shell uname -m)
VERSION=$(shell awk '/const Version/ { print $$4 }' version.go | sed 's/"//g')

.PHONY: ruby-centos
ruby-centos:
	echo "install"

.PHONY: ruby-ubuntu
ruby-ubuntu:
	apt-get install -y ruby bundler

.PHONY: ruby-install
ruby-install:
	[ -x /usr/bin/apt-get ] && make ruby-ubuntu || :
	[ -x /usr/bin/yum ]     && make ruby-centos || :

.PHONY: test
test: ruby-install
	(cd tests/rest-api && bundler install --deployment)
	echo "Starting the Rest API for testing"
	thin -d start -c tests/rest-api
	echo "Performing tests"
	sleep 5
	curl localhost:3000/v2/info
	thin stop -c tests/rest-api

.PHONY: changelog
changelog: release
	git log $(shell git tag | tail -n1)..HEAD --no-merges --format=%B > changelog

.PHONY: update
update:
	git pull
	make
