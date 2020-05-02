#!/usr/bin/env bash

#docker build --no-cache -f Dockerfile -t termlog-build .
docker build -f Dockerfile -t termlog-build .
#for arch in amd64 arm arm64; do
#  docker run --rm -v `pwd`:/tmp/dist termlog-build-arm64:latest cp /go/src/github.com/tzneal/ham-go/cmd/termlog/termlog.${arch} /tmp/dist/
#done
