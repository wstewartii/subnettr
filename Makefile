PATH := $(GOPATH)/bin:$(PATH)

all: build

build: subnettr test

test:
	go test -v

subnettr:
	go build -i -o build/subnettr

install:
	cp build/subnettr /usr/local/bin/subnettr
