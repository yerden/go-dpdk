# go-dpdk
[![Documentation](https://godoc.org/github.com/yerden/go-dpdk?status.svg)](http://godoc.org/github.com/yerden/go-dpdk) [![Go Report Card](https://goreportcard.com/badge/github.com/yerden/go-dpdk)](https://goreportcard.com/report/github.com/yerden/go-dpdk) [![Build Status](https://travis-ci.com/yerden/go-dpdk.svg?branch=master)](https://travis-ci.com/yerden/go-dpdk)
Go bindings for DPDK library.

# Build prereqs
If you have your own DPDK distribution build then do:
```
# set path to dpdk installation
export RTE_SDK=~/work/dpdk

# example to Linux@x86_x64
export RTE_TARGET=x86_64-native-linux-gcc
export CGO_CFLAGS="-m64 -pthread -O3 -march=native -I$RTE_SDK/$RTE_TARGET/include"
export CGO_LDFLAGS="-L$RTE_SDK/$RTE_TARGET/lib -ldpdk -lz -lrt -lnuma -ldl -lm"
```

If you use libdpdk-dev from Ubuntu then do:
```
sudo apt install libdpdk-dev libnuma-dev
export CGO_CFLAGS="-m64 -pthread -O3 -march=native -I/usr/include/dpdk"
export CGO_LDFLAGS="-L/usr/lib/x86_64-linux-gnu -ldpdk -lz -lrt -lnuma -ldl -lm"
```
