#!/bin/bash

set -ev

go test -v github.com/yerden/go-dpdk/lcore
go test -v github.com/yerden/go-dpdk/eal
go test -v github.com/yerden/go-dpdk/ring
go test -v github.com/yerden/go-dpdk/mempool
go test -v github.com/yerden/go-dpdk/port
