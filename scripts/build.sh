#!/bin/sh
export VERSION=1.5.3
docker build --build-arg VERSION=${VERSION} --no-cache -f Dockerfile -t api:${VERSION} .
