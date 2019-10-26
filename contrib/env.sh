#!/bin/bash

# set path to dpdk installation
# DPDK=$1

CGO_CFLAGS=""
if [ "x$DPDK" != "x" ]; then
	DPDK_INCLUDE="$DPDK/include"
	DPDK_LIB="$DPDK/lib"
fi

LD_LIBRARY_PATH="$LD_LIBRARY_PATH:$DPDK_LIB"
CGO_CFLAGS="$CGO_CFLAGS -m64 -pthread -O3 -march=native -I$DPDK_INCLUDE"
CGO_LDFLAGS="$CGO_LDFLAGS -L$DPDK_LIB"

# to activate mempool drivers constructors
# we need to link every symbol
DRIVERS="-lrte_mempool_ring -lrte_mempool_stack"

# --whole-archive ... --no-whole-archive is for static libraries
# --no-as-needed ... --as-needed is for shared libraries
CGO_LDFLAGS="$CGO_LDFLAGS -Wl,--no-as-needed -Wl,--whole-archive $DRIVERS -Wl,--no-whole-archive -Wl,--as-needed"

# -ldpdk should go AFTER drivers
CGO_LDFLAGS="$CGO_LDFLAGS -ldpdk -lz -lrt -lnuma -ldl -lm"

export CGO_CFLAGS CGO_LDFLAGS LD_LIBRARY_PATH
