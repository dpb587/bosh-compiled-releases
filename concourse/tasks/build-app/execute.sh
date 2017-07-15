#!/bin/bash

set -eu -o pipefail

task_dir="$PWD"

export GOPATH="$PWD/gopath"

cd gopath/src/github.com/dpb587/bosh-compiled-releases

./app/build

cp -rp app/* "$task_dir/app/"
