#!/bin/bash

go get -u github.com/client9/misspell/cmd/misspell
go get -u github.com/gordonklaus/ineffassign
go get -u github.com/fzipp/gocyclo/cmd/gocyclo@latest

DIRS="common eal ethdev lcore mempool ring port"
# Add subdirectories here as we clean up golint on each.
for subdir in $DIRS; do
  pushd $subdir
  misspell -error .
  gocyclo .
  ineffassign .
  popd
done
