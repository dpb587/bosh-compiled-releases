include "./concourse/pipeline-helpers";

{
  "jobs": [
    import_repo_release("github.com/dpb587/openvpn-bosh-release"),
    import_repo_release("github.com/dpb587/ssoca-bosh-release"),

    compile_boshio_release("cloudfoundry/bosh"; "ubuntu-trusty"),
    compile_boshio_release("concourse/concourse"; "ubuntu-trusty"),
    compile_boshio_release("cloudfoundry/garden-runc-release"; "ubuntu-trusty"),
    compile_boshio_release("cloudfoundry/uaa-release"; "ubuntu-trusty"),
    compile_boshio_release("cloudfoundry-community/docker-registry-boshrelease"; "ubuntu-trusty"),
    compile_boshio_release("cloudfoundry/syslog-release"; "ubuntu-trusty"),
    compile_boshio_release("cloudfoundry/bosh-dns-release"; "ubuntu-trusty"),

    {
      "name": "import-cf-deployment",
      "plan": [
        {
          "get": "cf-deployment",
          "trigger": true
        },
        {
          "get": "repo"
        },
        {
          "task": "import-cf-deployment",
          "file": "repo/concourse/tasks/import-cf-deployment/config.yml",
          "params": {
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
    },

    {
      "name": "push-app",
      "serial": true,
      "plan": [
        {
          "get": "repo",
          "trigger": true
        },
        {
          "task": "build-app",
          "file": "repo/concourse/tasks/build-app/config.yml"
        },
        {
          "put": "app",
          "params": {
            "manifest": "app/manifest.yml"
          }
        }
      ]
    }
  ],
  "resources": [
    {
      "name": "repo",
      "type": "git",
      "source": {
        "branch": "master",
        "private_key": "((git_private_key))",
        "uri": "git@github.com:dpb587/bosh-compiled-releases.git"
      }
    },
    {
      "name": "release-compiler",
      "type": "git",
      "source": {
        "uri": "https://github.com/dpb587/bosh-release-compiler.git"
      }
    },
    {
      "name": "ubuntu-trusty",
      "type": "metalink-repository",
      "source": {
        "uri": "https://dpb587.github.io/upstream-blob-mirror/repository/bosh.io/stemcell/bosh-warden-boshlite-ubuntu-trusty-go_agent/index.xml"
      }
    },

    repo_release("github.com/dpb587/openvpn-bosh-release"; {"uri": "https://github.com/dpb587/openvpn-bosh-release.git"}),
    repo_release("github.com/dpb587/ssoca-bosh-release"; {"uri": "https://github.com/dpb587/ssoca-bosh-release.git"}),

    boshio_release("cloudfoundry/bosh"),
    boshio_release("concourse/concourse"),
    boshio_release("cloudfoundry/garden-runc-release"),
    boshio_release("cloudfoundry/uaa-release"),
    boshio_release("cloudfoundry-community/docker-registry-boshrelease"),
    boshio_release("cloudfoundry/syslog-release"),
    boshio_release("cloudfoundry/bosh-dns-release"),

    {
      "name": "cf-deployment",
      "type": "git",
      "source": {
        "uri": "https://github.com/cloudfoundry/cf-deployment.git",
        "paths": [
          "operations/use-compiled-releases.yml"
        ]
      }
    },

    {
      "name": "app",
      "type": "cf",
      "source": {
        "api": "((cf_api))",
        "username": "((cf_username))",
        "password": "((cf_password))",
        "organization": "((cf_organization))",
        "space": "((cf_space))"
      }
    }
  ],
  "resource_types": [
    {
      "name": "metalink-repository",
      "type": "docker-image",
      "source": {
        "repository": "dpb587/metalink-repository-resource"
      }
    }
  ]
}
