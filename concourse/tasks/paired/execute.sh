#!/bin/bash

set -eu

release_dir=$HOME/Projects/dpb587/openvpn-bosh-release

for compiled_metalink in $( find "$release_dir/compiled_releases" -name '*.meta4' ); do
  version=$( meta4 file-version --metalink "$compiled_metalink" )
  stemcell_os=$( echo "$( basename "$compiled_metalink" )" | sed -E 's/.+-on-(.+)-stemcell-.+/\1/' )
  stemcell_version=$( echo "$( basename "$compiled_metalink" )" | sed -E 's/.+-stemcell-(.+)\.meta4/\1/' )
  release_name=$( basename "$( dirname "$compiled_metalink" )" )
  source_metalink="$release_dir/releases/$release_name/$release_name-$version.meta4"

  go run importer/paired/main.go \
    "$release_name" \
    "$source_metalink" \
    "$compiled_metalink" \
    "$stemcell_os" \
    "$stemcell_version" \
    > /tmp/metalink

  mv /tmp/metalink "data/$release_name/$( shasum /tmp/metalink | awk '{ print $1 }' ).meta4"
done
