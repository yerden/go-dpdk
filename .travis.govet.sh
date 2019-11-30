#!/bin/bash

cd "$(dirname $0)"

export CGO_CFLAGS_ALLOW=".*"
export CGO_LDFLAGS_ALLOW=".*"
export CGO_LDFLAGS="-Wl,--dynamic-linker=/lib64/ld-linux-x86-64.so.2"

set -e
DIRS="common eal ethdev lcore mempool ring port memzone"
# Add subdirectories here as we clean up golint on each.
for subdir in $DIRS; do
  pushd $subdir
  go vet -tags static
  popd
done
