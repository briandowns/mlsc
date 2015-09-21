GOCMD=go
GOBUILD=$(GOCMD) build
GOCLEAN=$(GOCMD) clean
GOINSTALL=$(GOCMD) install
GOTEST=$(GOCMD) test

all: build

build:
	$(GOBUILD) -v -o bin/lslock-test

install:
	$(GOINSTALL)

clean:
	$(GOCLEAN) -n -i -x
	rm -f $(GOPATH)/bin/lslock-test
	rm -rf bin/lslock-test

.PHONY: all clean