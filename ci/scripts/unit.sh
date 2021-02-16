#!/bin/bash -eux

cwd=$(pwd)

pushd $cwd/graphson
  make test
popd