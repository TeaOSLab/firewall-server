#!/usr/bin/env bash

ROOT=$(dirname "$0")
NAME="edge-firewall-server"
DIST=$ROOT/"../dist/${NAME}"
OS=${1}
ARCH=${2}

if [ -z "$OS" ]; then
	echo "usage: build.sh OS ARCH"
	exit
fi
if [ -z "$ARCH" ]; then
	echo "usage: build.sh OS ARCH"
	exit
fi

echo "building ${OS}/${ARCH}"

env GOOS="${OS}" GOARCH="${ARCH}" go build -ldflags "-s -w" -trimpath -o "${DIST}-${OS}-${ARCH}" "${ROOT}/../cmd/edge-firewall-server"

echo "[done]"