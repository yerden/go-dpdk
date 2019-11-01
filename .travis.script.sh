#!/bin/bash

set -ev

export CGO_CFLAGS="-mssse3 -msse4.1 -msse4.2 `pkg-config --cflags libdpdk`"
export CGO_LDFLAGS=`pkg-config --libs libdpdk`

go test github.com/yerden/go-dpdk/common
go test github.com/yerden/go-dpdk/lcore
go test github.com/yerden/go-dpdk/eal
go test github.com/yerden/go-dpdk/ring
go test github.com/yerden/go-dpdk/mempool
go test github.com/yerden/go-dpdk/memzone
go test github.com/yerden/go-dpdk/port
