#!/bin/bash

set -ev

. ./contrib/env.sh

go test -v github.com/yerden/go-dpdk/lcore
go test -v github.com/yerden/go-dpdk/eal
go test -v github.com/yerden/go-dpdk/ring
go test -v github.com/yerden/go-dpdk/mempool
go test -v github.com/yerden/go-dpdk/memzone
go test -v github.com/yerden/go-dpdk/port
