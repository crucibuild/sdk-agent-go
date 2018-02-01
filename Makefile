default: check

.PHONY: check ci dependencies

SCRIPTS_PATH:=$(CURDIR)/../scripts-build-go/script

dependencies:
	go get -v -d -u "github.com/crucibuild/scripts-build-go"
	go get -v -d ./...

check: dependencies
	$(SCRIPTS_PATH)/check.sh

coverage:
	go get golang.org/x/tools/cmd/cover
	go get github.com/go-playground/overalls
	"${GOPATH}/bin/overalls" -project=github.com/crucibuild/sdk-agent-go

ci: check coverage
	go get github.com/mattn/goveralls
	"${GOPATH}/bin/goveralls" -coverprofile=overalls.coverprofile -service=travis-ci
