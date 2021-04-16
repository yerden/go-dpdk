ARG DIST
ARG DPDK_VER
FROM nedrey/dpdk-build:${DIST}-${DPDK_VER}-sandybridge
ARG GO_VERSION
ENV DEBIAN_FRONTEND=noninteractive
ENV TZ=Asia/Almaty

RUN apt-get -y update && apt-get -y install \
		git \
		zlib1g-dev \
		pkg-config \
		gcc \
		curl \
		ibverbs-providers \
		libibverbs-dev \
		libmnl-dev \
		libjansson-dev \
		libnuma-dev \
		libpcap-dev \
		libisal-dev \
		libfdt-dev

RUN curl -SL https://dl.google.com/go/go${GO_VERSION}.linux-amd64.tar.gz |\
		tar -C /usr/local -xzf -

ENV GO111MODULE on
ENV GOPATH /go
ENV PATH /usr/local/go/bin:$GOPATH/bin:$PATH
ENV CGO_CFLAGS_ALLOW .*
ENV CGO_LDFLAGS_ALLOW .*
ENV PKG_CONFIG_PATH /usr/local/share/pkgconfig:$PKG_CONFIG_PATH

VOLUME /repo
WORKDIR /repo
