default: check

.PHONY: check ci dependencies

dependencies:
	go get -t -v ./...

check: dependencies
	go get -u github.com/alecthomas/gometalinter
	gometalinter --install --update
	gometalinter -j2 --config "$(CURDIR)/gometalinter.json" ./...

coverage:
	go get golang.org/x/tools/cmd/cover
	go get github.com/go-playground/overalls
	"${GOPATH}/bin/overalls" -project=github.com/crucibuild/sdk-agent-go

ci: check coverage
	go get github.com/mattn/goveralls
	"${GOPATH}/bin/goveralls" -coverprofile=overalls.coverprofile -service=travis-ci
