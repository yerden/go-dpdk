#!/bin/bash

set -ev

go get golang.org/x/sys/unix
go get github.com/yerden/go-dpdk/eal
go get github.com/yerden/go-dpdk/lcore
