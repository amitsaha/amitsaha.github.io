Title: Quick and dirty debian packages for your Golang application 
Date: 2018-01-25 15:00
Category: golang
Status: Draft

In this post, we will learn about a quick and easy workflow for
building and deploying your golang applications as Debian packages.


**Assumptions**

I have been using [dep]() for dependency management, and I assume that
you are doing the same. Other dependency management solutions should
work with only the specific bits of the workflow swapped out to suit
the one you may be using.

I also assume that you have `make` and a recent `golang` toolset
installed.

If you want to integrate my workflow into an existing project, 
please skip ahead to the second use case.

**Use case #1: New golang application project**

Create a new directory which will be the home of our new project.
Since we are going to use [dep](), it has to be somewhere in our
`GOPATH`. In my case, I will assume it is in 
`$GOPATH/src/github.com/amitsaha/packaging-demo`. The first file,
I will create is a `main.go` which looks like this:

```
package main

import (
	log "github.com/sirupsen/logrus"
)

func main() {
	log.Info("I love logrus!")
}
```

This is a simple program, but it uses a thirdy party package
[logrus](https://github.com/sirupsen/logrus) (which is awesome btw).


**Workflow - Step #1**

Now, we come to the first step of our workflow - create a file 
called `Makefile` with the following contents:

```
GOPATH := $(shell go env GOPATH)
GODEP_BIN := $(GOPATH)/bin/dep
GOLINT := $(GOPATH)/bin/golint
VERSION := $(shell cat VERSION)-$(shell git rev-parse --short HEAD)

packages = $$(go list ./... | egrep -v '/vendor/')
files = $$(find . -name '*.go' | egrep -v '/vendor/')

ifeq "$(HOST_BUILD)" "yes"
	# Use host system for building
	BUILD_SCRIPT =./build-deb-host.sh
else
	# Use docker for building
	BUILD_SCRIPT = ./build-deb-docker.sh
endif


.PHONY: all
all: lint vet test build 

$(GODEP):
	go get -u github.com/golang/dep/cmd/dep

Gopkg.toml: $(GODEP)
	$(GODEP_BIN) init

vendor:         ## Vendor the packages using dep
vendor: $(GODEP) Gopkg.toml Gopkg.lock
	@ echo "No vendor dir found. Fetching dependencies now..."
	GOPATH=$(GOPATH):. $(GODEP_BIN) ensure

version:
	@ echo $(VERSION)

build:          ## Build the binary
build: vendor
	test $(BINARY_NAME)
	go build -o $(BINARY_NAME) -ldflags "-X main.Version=$(VERSION)" 

build-deb:      ## Build DEB package (needs other tools)
	test $(BINARY_NAME)
	test $(DEB_PACKAGE_NAME)
	test "$(DEB_PACKAGE_DESCRIPTION)"
	exec ${BUILD_SCRIPT}
	
test: vendor
	go test -race $(packages)

vet:            ## Run go vet
vet: vendor
	go tool vet -printfuncs=Debug,Debugf,Debugln,Info,Infof,Infoln,Error,Errorf,Errorln $(files)

lint:           ## Run go lint
lint: vendor $(GOLINT)
	$(GOLINT) -set_exit_status $(packages)
$(GOLINT):
	go get -u github.com/golang/lint/golint

clean:
	test $(BINARY_NAME)
	rm -f $(BINARY_NAME) 

help:           ## Show this help
	@fgrep -h "##" $(MAKEFILE_LIST) | fgrep -v fgrep | sed -e 's/\\$$//' | sed -e 's/##//'
```

**Workflow - Step #2**

Next, we will create a file called `VERSION` in the project 
directory and write a string such as `0.1` into it:

```
$ echo "0.1" > VERSION
```

This will be our major version.

**Workflow - Step #3**

Initialize a git repository in the application directory:

```
$ git init
```

And we will make a first commit:

```
$ git add -A .
$ git commit -m "Initial commit"
```

**Workflow - Step #4**

Let's first try and see what our `Makefile` allows us to do:

```
$ make help
vendor:          Vendor the packages using dep
build:           Build the binary
build-deb:       Build DEB package (needs other tools)
vet:             Run go vet
lint:            Run go lint
help:            Show this help
```




## Resources

- [Daily Dep](https://golang.github.io/dep/docs/daily-dep.html)


