#!/bin/bash

cd "$(dirname $0)"

export CGO_LDFLAGS="-Wl,--dynamic-linker=/lib64/ld-linux-x86-64.so.2"

set -e
DIRS="common eal ethdev lcore mempool ring port memzone ethdev/flow"
# Add subdirectories here as we clean up golint on each.
for subdir in $DIRS; do
  pushd $subdir
  go vet -tags shared
  popd
done
