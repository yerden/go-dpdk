# go-dpdk
[![Documentation](https://godoc.org/github.com/yerden/go-dpdk?status.svg)](http://godoc.org/github.com/yerden/go-dpdk) [![Go Report Card](https://goreportcard.com/badge/github.com/yerden/go-dpdk)](https://goreportcard.com/report/github.com/yerden/go-dpdk) [![Build Status](https://travis-ci.com/yerden/go-dpdk.svg?branch=master)](https://travis-ci.com/yerden/go-dpdk)
Go bindings for DPDK library.

# Custom DPDK build
Currently, build is tested with DPDK 19.08.

If you have your own DPDK distribution build then do:
```
# set path to dpdk installation
export RTE_SDK=/path/to/dpdk

# build as shared lib if you need
# sed -i "s/\(CONFIG_RTE_BUILD_SHARED_LIB\)=\(y\|n\)/\1=y/" $RTE_SDK/config/common_base

# ... build as you want
export RTE_TARGET=x86_64-native-linuxapp-gcc
make config T=$RTE_TARGET O=mybuild
make O=mybuild
```

You may then use bundled `contrib/env.sh` script to setup DPDK environment.
You have to either specify path to your build via `DPDK` environment variable.
Alternatively, you may specify `DPDK_INCLUDE` and `DPDK_LIB` variable which point to
headers and library binaries. In this case, you should do:
```
DPDK=/path/to/dpdk/mybuild . ./contrib/env.sh
```

You may choose to link static executable, for example for containerized app:
```
go test --ldflags '-extldflags "-static"'
```

If you use libdpdk-dev from Ubuntu then do:
```
sudo apt install libdpdk-dev libnuma-dev
export CGO_CFLAGS="-m64 -pthread -O3 -march=native -I/usr/include/dpdk"
export CGO_LDFLAGS="-L/usr/lib/x86_64-linux-gnu -ldpdk -lz -lrt -lnuma -ldl -lm"
```

If you use dpdk-devel from CentOS then do:
```
sudo yum install zlib-devel numactl-devel dpdk-devel
export CGO_CFLAGS="-m64 -pthread -O3 -march=native -I/usr/include/dpdk"
export CGO_LDFLAGS="-L/usr/lib64 -ldpdk -lz -lrt -lnuma -ldl -lm"
```

# Meson build

go-dpdk supports compiling with the new DPDK building system based on
meson. This approach implies that DPDK libraries and drivers are
installed globally onto the system and the use of pkg-config is
encouraged to do static or dynamic linking.

To support this type of DPDK build, you should specify 'static' or
'shared' tag to Go compiler, for example `go build -tags static` for
static linking.

Due to Cgo usage considerations you should also allow any flag
returned by pkg-config. Also, on some systems you'd have to specify a
path to dynamic linker:
```
export CGO_CFLAGS_ALLOW=".*"
export CGO_LDFLAGS_ALLOW=".*"
export CGO_LDFLAGS="-Wl,--dynamic-linker=/lib64/ld-linux-x86-64.so.2"
```
