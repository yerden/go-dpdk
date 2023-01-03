# Go bindings for DPDK framework.
[![Documentation](https://godoc.org/github.com/yerden/go-dpdk?status.svg)](http://godoc.org/github.com/yerden/go-dpdk) [![Go Report Card](https://goreportcard.com/badge/github.com/yerden/go-dpdk)](https://goreportcard.com/report/github.com/yerden/go-dpdk) [![Build Status](https://github.com/yerden/go-dpdk/actions/workflows/unit.yml/badge.svg)](https://github.com/yerden/go-dpdk/actions/workflows/unit.yml)

# Building apps

Starting from DPDK 21.05, `pkg-config` becomes the only official way to build DPDK apps. Because of it `go-dpdk` uses `#cgo pkg-config` directive to link against your DPDK distribution.

Go compiler may fail to accept some C compiler flags. You can fix it by submitting those flags to environment:
```
export CGO_CFLAGS_ALLOW=".*"
export CGO_LDFLAGS_ALLOW=".*"
```

# Caveats
Only dynamic linking is viable at this point.
