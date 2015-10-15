#
#   Author: Rohith (gambol99@gmail.com)
#   Date: 2015-02-10 15:35:14 +0000 (Tue, 10 Feb 2015)
#
#  vim:ts=2:sw=2:et
#
NAME=go-marathon
AUTHOR=gambol99
HARDWARE=$(shell uname -m)
VERSION=$(shell awk '/const Version/ { print $$4 }' version.go | sed 's/"//g')
VETARGS?=-asmdecl -atomic -bool -buildtags -copylocks -methods -nilfunc -printf -rangeloops -shift -structtags -unsafeptr

.PHONY: test examples authors changelog

build:
	go build

authors:
	git log --format='%aN <%aE>' | sort -u > AUTHORS

deps:
	@echo "--> Installing build dependencies"
	@go get -d -v ./... $(DEPS)
	go get github.com/stretchr/testify/assert
	go get gopkg.in/yaml.v2

vet:
	@echo "--> Running go tool vet $(VETARGS) ."
	@go tool vet 2>/dev/null ; if [ $$? -eq 3 ]; then \
		go get golang.org/x/tools/cmd/vet; \
	fi
	@go tool vet $(VETARGS) .

cover:
	go list github.com/${AUTHOR}/${NAME} | xargs -n1 go test --cover

format:
	@echo "--> Running go fmt"
	@go fmt $(PACKAGES)

test: deps
	@echo "--> Running go tests"
	go test -v
	@$(MAKE) vet
	@$(MAKE) cover

changelog: release
	git log $(shell git tag | tail -n1)..HEAD --no-merges --format=%B > changelog
