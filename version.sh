#!/bin/bash

set -e

cd "$(dirname "$0")"
RELEASE=0
VERSION="1.0"
if REV=$(git rev-parse --short HEAD); then
    VERSION="${VERSION}-${REV}"
fi

cat > version.go <<HERE
package main
//Generated by go generate DO NOT EDIT

const version = "${VERSION}"
HERE
