#!/bin/bash

set -eu -o pipefail

git clone --quiet "file://$PWD/repo" repo-output

cd repo-output/

bosh interpolate ../cf-deployment/cf-deployment.yml \
  --path=/releases \
  > ../source-releases.yml

bosh interpolate ../cf-deployment/cf-deployment.yml \
  --ops-file=../cf-deployment/operations/use-compiled-releases.yml \
  --path=/releases \
  > ../compiled-releases.yml

for release_name in $( grep name: ../compiled-releases.yml | cut -c9- ); do
  repository="github.com/cloudfoundry/cf-deployment.$release_name"

  release_version=$( bosh interpolate ../compiled-releases.yml --path="/name=$release_name/version" )

  compiled_url=$( bosh interpolate ../compiled-releases.yml --path="/name=$release_name/url" )
  compiled_sha1=$( bosh interpolate ../compiled-releases.yml --path="/name=$release_name/sha1" )

  if [ ! -e data/github.com/cloudfoundry/cf-deployment.$release_name/bcr.json ]; then
    true # definitely do it
  elif grep "$compiled_sha1" data/github.com/cloudfoundry/cf-deployment/bcr.json | grep "$release_name" | grep -q "$release_version" ; then
    echo "skipping $release_name/$release_version"

    continue
  fi


  #
  # source release
  #

  source_url=$( bosh interpolate ../source-releases.yml --path="/name=$release_name/url" )
  source_sha1=$( bosh interpolate ../source-releases.yml --path="/name=$release_name/sha1" )

  curl -o ../source-release.tgz "$source_url"

  tarball_real=../source-release.tgz
  tarball_nice="$release_name-$release_version.tgz"

  metalink_path="data/$repository/releases/$release_name/$( basename "$tarball_nice" | sed 's/.tgz$//' ).meta4"

  mkdir -p "$( dirname "$metalink_path" )"

  meta4 create --metalink="$metalink_path"
  meta4 set-published --metalink="$metalink_path" "$( date -u +%Y-%m-%dT%H:%M:%SZ )"
  meta4 import-file --metalink="$metalink_path" --file="$tarball_nice" --version="$release_version" "$tarball_real"
  meta4 file-set-hash --metalink="$metalink_path" --file="$tarball_nice" sha-1 "$source_sha1"
  meta4 file-set-url --metalink="$metalink_path" --file="$tarball_nice" "$source_url"


  #
  # compiled release
  #

  curl -o ../compiled-release.tgz "$compiled_url"

  tarball_real=../compiled-release.tgz

  tar -xzf "$tarball_real" $( tar -tzf "$tarball_real" | grep release.MF$ )
  stemcell=$( set +o pipefail ; grep '^  stemcell:' release.MF | awk '{print $2}' | head -n1 | tr -d "\"'" )
  stemcell_os=$( echo "$stemcell" | cut -d/ -f1 )
  stemcell_version=$( echo "$stemcell" | cut -d/ -f2 )
  rm release.MF

  tarball_nice="$release_name-$release_version-on-$stemcell_os-stemcell-$stemcell_version.tgz"
  metalink_path="data/$repository/compiled_releases/$release_name/$( basename "$tarball_nice" | sed 's/.tgz$//' ).meta4"

  mkdir -p "$( dirname "$metalink_path" )"

  meta4 create --metalink="$metalink_path"
  meta4 set-published --metalink="$metalink_path" "$( date -u +%Y-%m-%dT%H:%M:%SZ )"
  meta4 import-file --metalink="$metalink_path" --file="$tarball_nice" --version="$release_version" "$tarball_real"
  meta4 file-set-hash --metalink="$metalink_path" --file="$tarball_nice" sha-1 "$compiled_sha1"
  meta4 file-set-url --metalink="$metalink_path" --file="$tarball_nice" "$compiled_url"


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

  git commit -m "$repository: add compiled release"

  ./concourse/tasks/regenerate-repository.sh "data/$repository" "data/$repository" "$repository"
done
