#!/bin/bash

set -ev

go test -v github.com/yerden/go-dpdk/eal
go test -v github.com/yerden/go-dpdk/lcore
