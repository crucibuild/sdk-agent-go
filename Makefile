PREFIX=github.com/crucibuild
NAME=sdk-agent-go
FULL_NAME=${PREFIX}/${NAME}
FULL_GOPATH=${GOPATH}/src/${FULL_NAME}

default: test

get:
	go get "github.com/omeid/go-resources/cmd/resources"
	go get -t -v ./...

build: get
	cd example/agent-ping && resources -output="resources.go" -var="Resources" -trim="../" resources/* ../schema/*
	cd example/agent-pong && resources -output="resources.go" -var="Resources" -trim="../" resources/* ../schema/*
	go build -v ./...

test: build
	go test -v ./...
