#!/bin/bash

cd "$(dirname $0)"

export CGO_CFLAGS="-mssse3 -msse4.1 -msse4.2 `pkg-config --cflags libdpdk`"
export CGO_LDFLAGS=`pkg-config --libs libdpdk`

set -e
DIRS="common eal ethdev lcore mempool ring port memzone"
# Add subdirectories here as we clean up golint on each.
for subdir in $DIRS; do
  pushd $subdir
  go vet
  popd
done
