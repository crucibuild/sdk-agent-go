default: all

SCRIPTS_PATH:=$(CURDIR)/../scripts-build-go/script

.PHONY: all check dependencies

dependencies:
	go get -v -d -u "github.com/crucibuild/scripts-build-go"
	go get -v -d -t ./...

build: dependencies
	go build ./...

test: build
	go test ./...

check:
	$(SCRIPTS_PATH)/check.sh

coverage:
	$(SCRIPTS_PATH)/gen_coverage.sh

all: test check coverage
	$(SCRIPTS_PATH)/push_coverage.sh

ci: all
	true
