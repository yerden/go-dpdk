# Go bindings for DPDK framework.
[![Documentation](https://godoc.org/github.com/yerden/go-dpdk?status.svg)](http://godoc.org/github.com/yerden/go-dpdk) [![Build Status](https://github.com/yerden/go-dpdk/actions/workflows/unit.yml/badge.svg)](https://github.com/yerden/go-dpdk/actions/workflows/unit.yml) [![codecov](https://codecov.io/gh/yerden/go-dpdk/branch/master/graph/badge.svg?token=1XW04KL02S)](https://codecov.io/gh/yerden/go-dpdk)

# Building apps

Starting from DPDK 21.05, `pkg-config` becomes the only official way to build DPDK apps. Because of it `go-dpdk` uses `#cgo pkg-config` directive to link against your DPDK distribution.

Go compiler may fail to accept some C compiler flags. You can fix it by submitting those flags to environment:
```
export CGO_CFLAGS_ALLOW="-mrtm"
export CGO_LDFLAGS_ALLOW="-Wl,--(?:no-)?whole-archive"
```

# Caveats
* Only dynamic linking is viable at this point.
* If you isolate CPU cores with `isolcpus` kernel parameter then `GOMAXPROCS` should be manually specified to reflect the actual number of logical cores in CPU mask. E.g. if `isolcpus=12-95` on a 96-core machine then default value for `GOMAXPROCS` would be 12 but it should be at least 84.
