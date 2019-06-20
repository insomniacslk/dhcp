#!/usr/bin/env bash

# because things are never simple.
# See https://github.com/codecov/example-go#caveat-multiple-files

set -e
echo "" > coverage.txt

# show the network configuration. This can help troubleshooting integration
# tests.
ip a

for d in $(go list ./... | grep -v vendor); do
    go test -race -coverprofile=profile.out -covermode=atomic $d
    if [ -f profile.out ]; then
        cat profile.out >> coverage.txt
        rm profile.out
    fi
    # integration tests
    go test -c -cover -tags=integration -race -covermode=atomic $d
    testbin="./$(basename $d).test"
    # only run it if it was built - i.e. if there are integ tests
    test -x "${testbin}" && sudo "./${testbin}" -test.coverprofile=profile.out
    if [ -f profile.out ]; then
        cat profile.out >> coverage.txt
        rm -f profile.out
    fi
done

# check that we are not breaking some projects that depend on us. Remove this after moving to
# Go versioned modules, see https://github.com/insomniacslk/dhcp/issues/123

# Skip go1.9 for this check. rtr7/router7 depends on miekg/dns, which does not
# support go1.9
if [[ "$TRAVIS_GO_VERSION" =~ ^1.(9|10|11)$ ]]
then
    exit 0
fi

go get github.com/rtr7/router7/cmd/...
cd "${GOPATH}/src/github.com/rtr7/router7"
go build github.com/rtr7/router7/cmd/...
