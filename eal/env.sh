#!/bin/bash

# set path to dpdk installation
export RTE_SDK=$1

# example to Linux@x86_x64
export RTE_TARGET=x86_64-native-linux-gcc
export CGO_CFLAGS="-m64 -pthread -O3 -march=native -I$RTE_SDK/$RTE_TARGET/include"
export CGO_LDFLAGS="-L$RTE_SDK/$RTE_TARGET/lib -ldpdk -lz -lrt -lnuma -ldl -lm"
