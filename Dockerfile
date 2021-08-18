ARG DPDK_DIST
FROM nedrey/dpdk-build:${DPDK_DIST}
ARG GO_VERSION
ENV GOPATH /go
ENV PATH /usr/local/go/bin:$GOPATH/bin:$PATH
ENV CGO_CFLAGS_ALLOW .*
ENV CGO_LDFLAGS_ALLOW .*
ENV PKG_CONFIG_PATH /usr/local/share/pkgconfig:$PKG_CONFIG_PATH
ENV GO111MODULE on

RUN dnf install -y epel-release dnf-plugins-core && \
        dnf config-manager --set-enabled powertools && \
        dnf install -y libibverbs rdma-core-devel \
                jansson-devel zlib-devel gcc make git curl pkg-config \
                libpcap-devel numactl-devel && \
        (curl -SL https://dl.google.com/go/go${GO_VERSION}.linux-amd64.tar.gz |\
                tar -C /usr/local -xzf -)
