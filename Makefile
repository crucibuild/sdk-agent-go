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

ci: build
	go test -v -race ./...
	go vet ./...
	go get golang.org/x/tools/cmd/cover
	go get github.com/mattn/goveralls
	go get github.com/go-playground/overalls
	"${GOPATH}/bin/overalls" -project=github.com/crucibuild/sdk-agent-go
	"${GOPATH}/bin/goveralls" -coverprofile=overalls.coverprofile -service=travis-ci
