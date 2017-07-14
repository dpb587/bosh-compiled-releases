#!/bin/bash

set -eu -o pipefail

release_dir="${release_dir:-data/$repository}"

git clone --quiet "file://$PWD/repo" repo-output

cd repo-output

git config --global user.email "$git_user_email"
git config --global user.name "$git_user_name"

export GIT_COMMITTER_NAME="Concourse"
export GIT_COMMITTER_EMAIL="concourse.ci@localhost"

./concourse/tasks/regenerate-repository.sh "${release_dir:-data/$repository}" "data/$repository" "$repository"
