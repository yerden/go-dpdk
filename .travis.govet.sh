#!/bin/bash

cd "$(dirname $0)"

. ./contrib/env.sh

set -e
DIRS="common eal ethdev lcore mempool ring port memzone"
# Add subdirectories here as we clean up golint on each.
for subdir in $DIRS; do
  pushd $subdir
  go vet
  popd
done
