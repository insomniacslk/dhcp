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

# from https://github.com/rtr7/router7/blob/aa404c3c54d9ad655479d7978ed18e81fe6ca05c/.travis.yml#L14
# TODO: get rid of this once https://github.com/google/gopacket/pull/470 is merged
go get github.com/google/gopacket/pcapgo
(cd $GOPATH/src/github.com/google/gopacket && wget -qO- https://patch-diff.githubusercontent.com/raw/google/gopacket/pull/470.patch | patch -p1)

go get github.com/rtr7/router7/cmd/...
cd "${GOPATH}/src/github.com/rtr7/router7"
go build github.com/rtr7/router7/cmd/...
