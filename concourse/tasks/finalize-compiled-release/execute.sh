#!/bin/bash

set -eu -o pipefail

export AWS_ACCESS_KEY_ID="$s3_access_key"
export AWS_SECRET_ACCESS_KEY="$s3_secret_key"

git clone --quiet "file://$PWD/repo" repo-output

tar -xzf compiled-release/*.tgz $( tar -tzf compiled-release/*.tgz | grep release.MF$ )
version=$( grep '^version:' release.MF | awk '{print $2}' | tr -d "\"'" )

cd repo-output/


#
# source release
#

tarball_real=../release/release.tgz

tar -xzf "$tarball_real" $( tar -tzf "$tarball_real" | grep release.MF$ )
release_name=$( grep '^name:' release.MF | awk '{print $2}' | tr -d "\"'" )
release_version=$( grep '^version:' release.MF | awk '{print $2}' | tr -d "\"'" )
rm release.MF

tarball_nice="$release_name-$release_version.tgz"

metalink_path="data/$repository/releases/$release_name/$( basename "$tarball_nice" | sed 's/.tgz$//' ).meta4"

mkdir -p "$( dirname "$metalink_path" )"

meta4 create --metalink="$metalink_path"
meta4 set-published --metalink="$metalink_path" "$( date -u +%Y-%m-%dT%H:%M:%SZ )"
meta4 import-file --metalink="$metalink_path" --file="$tarball_nice" --version="$version" "$tarball_real"
meta4 file-set-url --metalink="$metalink_path" --file="$tarball_nice" "$( cat ../release/url )"


#
# compiled release
#

tarball_real=$( echo "../compiled-release/$release_name"-*.tgz )
tarball_nice="$( basename "$( echo "$tarball_real" | sed -E 's/-compiled-1.+.tgz/.tgz/' )" )"

tar -xzf "$tarball_real" $( tar -tzf "$tarball_real" | grep release.MF$ )
version=$( grep '^version:' release.MF | awk '{print $2}' | tr -d "\"'" )
stemcell=$( grep 'stemcell:' release.MF | head -n1 | awk '{print $2}' | tr -d "\"'" )
stemcell_os=$( echo "$stemcell" | cut -d/ -f1 )
stemcell_version=$( echo "$stemcell" | cut -d/ -f2 )
rm release.MF

metalink_path="data/$repository/releases/$release_name/$stemcell_os/$stemcell_version/$release_name-$version.meta4"

mkdir -p "$( dirname "$metalink_path" )"

meta4 create --metalink="$metalink_path"
meta4 set-published --metalink="$metalink_path" "$( date -u +%Y-%m-%dT%H:%M:%SZ )"
meta4 import-file --metalink="$metalink_path" --file="$tarball_nice" --version="$version" "$tarball_real"
meta4 file-upload --metalink="$metalink_path" --file="$tarball_nice" "$tarball_real" "s3://$s3_host/$s3_bucket/$release_name/$( meta4 file-hash --metalink="$metalink_path" sha-256 )"


#
# commit
#

if [[ -z "$( git status --porcelain )" ]]; then
  exit
fi

git config --global user.email "$git_user_email"
git config --global user.name "$git_user_name"

export GIT_COMMITTER_NAME="Concourse"
export GIT_COMMITTER_EMAIL="concourse.ci@localhost"

git add .

git commit -m "$repository: $stemcell"
