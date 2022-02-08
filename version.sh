#!/bin/sh

version=$(git log --date=iso --pretty=format:"%cd @%H" -1)
if [ $? -ne 0 ]; then
  version="${version}"
fi

describe=$(git describe --tags 2>/dev/null)
if [ $? -eq 0 ]; then
  version="${version} @${describe}"
fi

compile=$(go version)
if [ $? -ne 0 ]; then
  compile="unknown go version"
fi

build=$(date +"%F %T %z")
if [ $? -ne 0 ]; then
  compile="unknown build time"
fi

cat <<EOF >internal/version.go
package internal

const (
    Build   = "$build"
    Compile = "$compile"
    Version = "$version"
)
EOF