default: check

.PHONY: check ci dependencies

dependencies:
	go get -t -v ./...

check: dependencies
	! gofmt -d . 2>&1 | read
	! golint ./... 2>&1 | read
	go test -v -race ./...
	go vet ./...

coverage:
	go get golang.org/x/tools/cmd/cover
	go get github.com/go-playground/overalls
	"${GOPATH}/bin/overalls" -project=github.com/crucibuild/sdk-agent-go

ci: check coverage
	go get github.com/mattn/goveralls
	"${GOPATH}/bin/goveralls" -coverprofile=overalls.coverprofile -service=travis-ci
