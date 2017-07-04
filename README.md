## Client

Use the client to resolve compiled releases from a compilation server...

    $ bosh-bcr https://dpb587-bosh-compiled-releases.cfapps.io manifest.yml

Inline, it might look like...

    $ bosh deploy <( bosh-bcr https://dpb587-bosh-compiled-releases.cfapps.io manifest.yml )

The commands make several assumptions...

 * `releases` - each release must include `name`, `sha1`, `url`, and `version`
 * `stemcells`/`resource_pools.stemcell` - each stemcell must include `os` and `version`


## Server

Use the server to provide a compilation lookup server...

    $ PORT=8080 go run server/cli/main.go "data/*/*/*/bcr.json"

Find a compiled release...

    $ echo '{"name":"openvpn","version":"4.0.0","sha1":"cc14b757e5ac9af99840167c10114845b51da41d","stemcell":{"os":"ubuntu-trusty","version":"3421.11"}}' \
      | curl -XGET -d@- http://localhost:8080/resolve
    {
      "compiled_release": {
        "sha1": "19e79e45b690bc933b0ff5d9e54574f25d0899b9",
        "url": "https://s3-external-1.amazonaws.com/dpb587-bosh-release-openvpn-us-east-1/compiled_releases/openvpn/openvpn-4.0.0-on-ubuntu-trusty-stemcell-3421.11-compiled-1.20170630134749.0.tgz"
      }
    }


## Pipeline

    fly -t dpb587-nightwatch-aws-use1 sp -p dpb587:bosh-compiled-releases -c <( bosh int -l concourse/secrets.yml concourse/pipeline.yml )


## App

    ./app/push dpb587-bosh-compiled-releases
