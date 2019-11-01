#!/bin/bash

cd "$(dirname $0)"

CGO_CFLAGS=`pkg-config --cflags libdpdk`
CGO_LDFLAGS=`pkg-config --libs libdpdk`

go get golang.org/x/lint/golint
DIRS="lcore eal ring port mempool memzone"
# Add subdirectories here as we clean up golint on each.
for subdir in $DIRS; do
  pushd $subdir
  if golint|grep .; then
    echo "golint $subdir failed"
    exit 1
  fi
  popd
done
