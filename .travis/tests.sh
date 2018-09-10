#!/usr/bin/env bash

# because things are never simple.
# See https://github.com/codecov/example-go#caveat-multiple-files

set -e
echo "" > coverage.txt

for d in $(go list ./... | grep -v vendor); do
    go test -race -coverprofile=profile.out -covermode=atomic $d
    if [ -f profile.out ]; then
        cat profile.out >> coverage.txt
        rm profile.out
    fi
done

# check that we are not breaking some projects that depend on us. Remove this after moving to
# Go versioned modules, see https://github.com/insomniacslk/dhcp/issues/123
go get github.com/rtr7/router7/cmd/...
cd "${GOPATH}/src/github.com/rtr7/router7"
go build github.com/rtr7/router7/cmd/...
