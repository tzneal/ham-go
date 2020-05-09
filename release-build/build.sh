#!/usr/bin/env bash

if [ $# -ne 1 ]; then
  echo 1>&2 "Usage: $0 tagname"
  exit 3
fi
docker build -f Dockerfile --build-arg TAG=$1 -t termlog-build .
for arch in amd64 arm arm64; do
  docker run --rm -v `pwd`:/tmp/dist termlog-build:latest cp /go/src/github.com/tzneal/ham-go/cmd/termlog/termlog.${arch} /tmp/dist/
done
