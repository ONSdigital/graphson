#!/bin/bash -eux

cwd=$(pwd)

pushd $cwd/graphson
  make build
popd