PREFIX=github.com/crucibuild
NAME=sdk-agent-go
FULL_NAME=${PREFIX}/${NAME}
FULL_GOPATH=${GOPATH}/src/${FULL_NAME}

default: build

build:
	go get "github.com/omeid/go-resources/cmd/resources"
	cd example/agent-ping && resources -output="resources.go" -var="Resources" -trim="../" resources/* ../schema/*
	cd example/agent-pong && resources -output="resources.go" -var="Resources" -trim="../" resources/* ../schema/*
	go build "${FULL_NAME}/example/agent-ping"
	go build "${FULL_NAME}/example/agent-pong"

install: build
	go install "${FULL_NAME}/example/agent-ping" "${FULL_NAME}/example/agent-pong" "${FULL_NAME}"

clean: env
	go clean "${FULL_NAME}"
