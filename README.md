# CLI

## `rewrite-manifest`

The `rewrite-manifest` command converts source release references into compiled release references by querying a remote server or local index. The updated manifest is sent to standard output.

    $ bcr rewrite-manifest --server=https://dpb587-bosh-compiled-releases.cfapps.io manifest.yml

The command make several assumptions...

 * `releases` - each release must include `name`, `sha1`, `url`, and `version`
 * `stemcells`/`resource_pools.stemcell` - each stemcell must include `os` and `version`


## `serve`

The `serve` command starts a simple HTTP server which can be queried to resolve source releases to compiled releases based on a locally-accessible indices.

    $ bcr serve --local=data/*/*/*/bcr.json*


### API


#### GET `/v1/resolve`

Convert a source release reference to a compiled release reference.

Request Body (`application/json`)

    {
      "name": string,
      "version": string,
      "sha1": string,
      "stemcell": {
        "os": string,
        "version": string
      }
    }

Response Body (`application/json`)

    {
      "compiled_release": {
        "sha1": string,
        "url": string
      }
    }

Example

    $ echo '{"name":"openvpn","version":"4.0.0","sha1":"cc14b757e5ac9af99840167c10114845b51da41d","stemcell":{"os":"ubuntu-trusty","version":"3421.11"}}' \
      | curl -XGET -d@- http://localhost:8080/resolve
    {
      "compiled_release": {
        "sha1": "19e79e45b690bc933b0ff5d9e54574f25d0899b9",
        "url": "https://s3-external-1.amazonaws.com/dpb587-bosh-release-openvpn-us-east-1/compiled_releases/openvpn/openvpn-4.0.0-on-ubuntu-trusty-stemcell-3421.11-compiled-1.20170630134749.0.tgz"
      }
    }


# Deployment

## Concourse

    fly set-pipeline -p dpb587:bosh-compiled-releases -c <( bosh int -l concourse/secrets.yml concourse/pipeline.yml )


## Cloud Foundry App

    ./app/push dpb587-bosh-compiled-releases
