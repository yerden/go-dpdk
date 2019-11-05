ARG DIST
ARG DPDK_VER
FROM nedrey/dpdk-build:${DIST}-${DPDK_VER}
ARG GO_VERSION

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
		libpcap-dev

RUN curl -SL https://dl.google.com/go/go${GO_VERSION}.linux-amd64.tar.gz |\
		tar -C /usr/local -xzf -

ENV GO111MODULE on
ENV GOPATH /go
ENV PATH /usr/local/go/bin:$GOPATH/bin:$PATH

COPY . /work
WORKDIR /work
