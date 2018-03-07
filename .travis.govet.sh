#!/bin/bash

cd "$(dirname $0)"
DIRS=". dhcpv4 dhcpv6 iana"
set -e
for subdir in $DIRS; do
  pushd $subdir
  go vet
  popd
done
