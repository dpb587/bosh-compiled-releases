#!/bin/bash

set -eu -o pipefail

release_dir="$1"
data_dir="$2"
repository="$3"

for compiled_metalink in $( set -eu ; cd "$release_dir" ; find releases -mindepth 4 -name '*.meta4' ); do
  compiled_metalink="$release_dir/$compiled_metalink"

  stemcell_version="$( echo "$( basename "$( dirname "$compiled_metalink" )" )" )"
  stemcell_os="$( echo "$( basename "$( dirname "$( dirname "$compiled_metalink" )" )" )" )"

  release_name="$( echo "$( basename "$( dirname "$( dirname "$( dirname "$compiled_metalink" )" )" )" )" )"
  release_version="$( meta4 file-version --metalink "$compiled_metalink" )"

  compiled_digest="$( meta4 file-hash --metalink "$compiled_metalink" sha-1 )"
  compiled_url="$( meta4 file-urls --metalink "$compiled_metalink" | head -n1 )"

  source_metalink="$release_dir/releases/$release_name/$release_name-$release_version.meta4"
  source_digest="$( meta4 file-hash --metalink "$source_metalink" sha-1 )"

  bcr file-add-compiled-release \
    "$data_dir/bcr.json" \
    "$release_name" \
    "$release_version" \
    "$source_digest" \
    "$stemcell_os" \
    "$stemcell_version" \
    "$compiled_digest" \
    "$compiled_url"
done

if [[ -z "$( git status --porcelain )" ]]; then
  exit
fi

git add .

git commit -m "$repository: update bcr.json"
