#!/bin/bash

set -ev

TAGS=$1

export CGO_LDFLAGS="-Wl,--dynamic-linker=/lib64/ld-linux-x86-64.so.2"

DIRS="common lcore eal ring mempool memzone port"
echo "Testing $TAGS"
for subdir in $DIRS; do
  go test -tags $TAGS github.com/yerden/go-dpdk/$subdir
done
