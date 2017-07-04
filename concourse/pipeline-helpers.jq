def slug($string):
  $string | split("/") | join("-") | gsub("github.com"; "github")
;

def boshio_release($release):
  {
    "name": ("bosh-io-" + slug($release)),
    "type": "bosh-io-release",
    "source": {
      "repository": $release
    }
  }
;

def repo_release($release; $source):
  {
    "name": slug($release),
    "type": "git",
    "source": $source
  }
;

def compile_boshio_release($release; $stemcellref):
  {
    "name": ("bosh-io-" + slug($release) + "-" + slug($stemcellref)),
    "serial_groups": [
      "compiler"
    ],
    "plan": [
      {
        "aggregate": [
          {
            "get": "release",
            "resource": ("bosh-io-" + slug($release)),
            "trigger": true
          },
          {
            "get": "stemcell",
            resource: slug($stemcellref),
            trigger: true
          }
        ]
      },
      {
        "get": "release-compiler"
      },
      {
        "task": "compile-release",
        "file": "release-compiler/concourse/execute-local-bosh.yml",
        "privileged": true
      },
      {
        "get": "repo"
      },
      {
        "task": "finalize-compiled-release",
        "file": "repo/concourse/tasks/finalize-compiled-release/config.yml",
        "params": {
          "repository": ("github.com/" + $release),
          "s3_host": "((s3_host))",
          "s3_bucket": "((s3_bucket))",
          "s3_access_key": "((s3_access_key))",
          "s3_secret_key": "((s3_secret_key))",
          "git_commit_message": ("github.com/" + $release + ": add compiled release"),
          "git_user_email": "((maintainer_email))",
          "git_user_name": "((maintainer_name))"
        }
      },
      {
        "task": "import-metalinks",
        "file": "repo/concourse/tasks/import-metalinks/config.yml",
        "params": {
          "repository": ("github.com/" + $release),
          "git_commit_message": ("github.com/" + $release + ": update bcr.json"),
          "git_user_email": "((maintainer_email))",
          "git_user_name": "((maintainer_name))"
        }
      },
      {
        "put": "repo",
        "params": {
          "rebase": true,
          "repository": "repo"
        }
      }
    ]
  }
;


def import_repo_release($release):
  {
    "name": slug($release),
    "plan": [
      {
        "get": "release",
        "resource": slug($release),
        "trigger": true
      },
      {
        "get": "repo"
      },
      {
        "task": "import-metalinks",
        "file": "repo/concourse/tasks/import-metalinks/config-release.yml",
        "params": {
          "repository": $release,
          "git_commit_message": ($release + ": update bcr.json"),
          "git_user_email": "((maintainer_email))",
          "git_user_name": "((maintainer_name))"
        }
      },
      {
        "put": "repo",
        "params": {
          "rebase": true,
          "repository": "repo"
        }
      }
    ]
  }
;
