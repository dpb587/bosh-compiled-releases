#!/bin/bash

set -eu -o pipefail

cd repo

./app/build

cp -rp app/* ../app/
