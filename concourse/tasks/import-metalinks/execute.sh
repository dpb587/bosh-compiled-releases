#!/bin/bash

set -eu -o pipefail

release_dir="${release_dir:-data/$repository}"

git clone --quiet "file://$PWD/repo" repo-output

cd repo-output

for compiled_metalink in $( find "$release_dir/compiled_releases" -name '*.meta4' ); do
  release_name="$( basename "$( dirname "$compiled_metalink" )" )"
  release_version="$( meta4 file-version --metalink "$compiled_metalink" )"

  stemcell_os="$( echo "$( basename "$compiled_metalink" )" | sed -E 's/.+-on-(.+)-stemcell-.+/\1/' )"
  stemcell_version="$( echo "$( basename "$compiled_metalink" )" | sed -E 's/.+-stemcell-(.+)\.meta4/\1/' )"

  compiled_digest="$( meta4 file-hash --metalink "$compiled_metalink" sha-1 )"
  compiled_url="$( meta4 file-urls --metalink "$compiled_metalink" | head -n1 )"

  source_metalink="$release_dir/releases/$release_name/$release_name-$release_version.meta4"
  source_digest="$( meta4 file-hash --metalink "$source_metalink" sha-1 )"

  bcr file-add-compiled-release \
    "data/$repository/bcr.json" \
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

git config --global user.email "$git_user_email"
git config --global user.name "$git_user_name"

export GIT_COMMITTER_NAME="Concourse"
export GIT_COMMITTER_EMAIL="concourse.ci@localhost"

git add .

git commit -m "$git_commit_message"
