FROM golang:1.17-buster

RUN dpkg --add-architecture arm64 \
    && apt update \
    && apt install -y --no-install-recommends \
        gcc-aarch64-linux-gnu \
        libc6-dev-arm64-cross \
        pkg-config \
        libhamlib2:arm64 \
        libhamlib-dev:arm64 \
    && rm -rf /var/lib/apt/lists/*

RUN dpkg --add-architecture armhf \
    && apt-get update \
    && apt-get install -y --no-install-recommends \
       gcc-arm-linux-gnueabihf \
       libc6-dev-armhf-cross \
       libhamlib2:armhf \
       libhamlib-dev:armhf \
    && rm -rf /var/lib/apt/lists/*

RUN apt-get update \
    && apt-get install -y --no-install-recommends \
       gcc \
       libhamlib2 \
       libhamlib-dev \
    && rm -rf /var/lib/apt/lists/*

ARG TAG
RUN mkdir -p /go/src/github.com/tzneal/
RUN cd /go/src/github.com/tzneal/ && git clone -b $TAG https://github.com/tzneal/ham-go.git
RUN cd /go/src/github.com/tzneal/ham-go/release-build && make write-version
RUN cd /go/src/github.com/tzneal/ham-go && git log HEAD~1


ENV GOOS=linux
ENV CGO_ENABLED=1

ENV GOARCH=arm64
ENV CC=aarch64-linux-gnu-gcc
ENV PATH="/go/bin/${GOOS}_${GOARCH}:${PATH}"
ENV PKG_CONFIG_PATH=/usr/lib/aarch64-linux-gnu/pkgconfig

RUN GOARCH=${GOARCH} cd /go/src/github.com/tzneal/ham-go/cmd/termlog && go build
RUN mv /go/src/github.com/tzneal/ham-go/cmd/termlog/termlog /go/src/github.com/tzneal/ham-go/cmd/termlog/termlog.${GOARCH}

ENV GOARCH=arm
ENV CC=arm-linux-gnueabihf-gcc
ENV PATH="/go/bin/${GOOS}_${GOARCH}:${PATH}"
ENV PKG_CONFIG_PATH=/usr/lib/arm-linux-gnueabihf/pkgconfig

RUN GOARCH=${GOARCH} cd /go/src/github.com/tzneal/ham-go/cmd/termlog && go build
RUN mv /go/src/github.com/tzneal/ham-go/cmd/termlog/termlog /go/src/github.com/tzneal/ham-go/cmd/termlog/termlog.${GOARCH}

ENV GOARCH=amd64
ENV GOARM=5
ENV CC=gcc
ENV PATH="/go/bin/${GOOS}_${GOARCH}:${PATH}"
ENV PKG_CONFIG_PATH=/usr/lib/x86_64-linux-gnu/pkgconfig

RUN GOARCH=${GOARCH} cd /go/src/github.com/tzneal/ham-go/cmd/termlog && go build
RUN mv /go/src/github.com/tzneal/ham-go/cmd/termlog/termlog /go/src/github.com/tzneal/ham-go/cmd/termlog/termlog.${GOARCH}
