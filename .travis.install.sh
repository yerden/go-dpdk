#!/bin/bash

set -ev

CGO_CFLAGS=`pkg-config --cflags libdpdk`
CGO_LDFLAGS=`pkg-config --libs libdpdk`

go get golang.org/x/sys/unix
go get github.com/yerden/go-dpdk/eal
go get github.com/yerden/go-dpdk/lcore
