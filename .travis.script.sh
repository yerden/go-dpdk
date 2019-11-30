#!/bin/bash

set -ev

TAGS=$1

export CGO_CFLAGS_ALLOW=".*"
export CGO_LDFLAGS_ALLOW=".*"
export CGO_LDFLAGS="-Wl,--dynamic-linker=/lib64/ld-linux-x86-64.so.2"

DIRS="common lcore eal ring mempool memzone"
echo "Testing $TAGS"
for subdir in $DIRS; do
  go test -tags $TAGS github.com/yerden/go-dpdk/$subdir
done
