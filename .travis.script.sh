#!/bin/bash

set -ev

TAGS=$1

DIRS="common lcore eal ring mempool memzone ethdev port ethdev/flow"
echo "Testing $TAGS"
for subdir in $DIRS; do
  go test -tags $TAGS github.com/yerden/go-dpdk/$subdir
done
