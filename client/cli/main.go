package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strconv"
	"strings"

	"github.com/dpb587/bosh-compiled-releases/server/api"

	yaml "gopkg.in/yaml.v2"
)

func main() {
	var manifest map[string]interface{}

	bytes, err := ioutil.ReadFile(os.Args[2])
	if err != nil {
		log.Fatal(err)
	}

	err = yaml.Unmarshal(bytes, &manifest)
	if err != nil {
		log.Fatal(err)
	}

	var stemcell map[interface{}]interface{}

	_, hasStemcells := manifest["stemcells"]
	_, hasResourcePools := manifest["resource_pools"]

	if hasStemcells && len(manifest["stemcells"].([]interface{})) == 1 {
		stemcell = manifest["stemcells"].([]interface{})[0].(map[interface{}]interface{})
	} else if hasResourcePools && len(manifest["resource_pools"].([]interface{})) == 1 {
		stemcell = manifest["resource_pools"].([]interface{})[0].(map[interface{}]interface{})["stemcell"].(map[interface{}]interface{})
	} else {
		log.Fatal("failed to identify stemcell")
	}

	var stemcellName string
	var stemcellOS string
	var stemcellVersion string

	if _, found := stemcell["name"]; found {
		stemcellName = stemcell["name"].(string)
	}

	if _, found := stemcell["os"]; found {
		stemcellOS = stemcell["os"].(string)
	}

	if _, found := stemcell["version"]; found {
		stemcellVersion = stemcell["version"].(string)
	}

	if stemcellOS == "" {
		if strings.Contains(stemcellName, "ubuntu-trusty") {
			stemcellOS = "ubuntu-trusty"
		} else if strings.Contains(stemcellName, "centos-7") {
			stemcellOS = "centos-7"
		}
	}

	for releaseIdx, releaseRaw := range manifest["releases"].([]interface{}) {
		release := releaseRaw.(map[interface{}]interface{})

		if _, found := release["name"]; !found {
			continue
		} else if _, found := release["version"]; !found {
			continue
		} else if _, found := release["sha1"]; !found {
			continue
		}

		var version string

		switch release["version"].(type) {
		case string:
			version = release["version"].(string)
		case int:
			version = strconv.Itoa(release["version"].(int))
		default:
			panic("unexpected type")
		}

		wBytes, err := json.Marshal(api.ResolveRequest{
			Name:    release["name"].(string),
			Version: version,
			Sha1:    release["sha1"].(string),
			Stemcell: api.ResolveRequestStemcell{
				OS:      stemcellOS,
				Version: stemcellVersion,
			},
		})
		if err != nil {
			log.Fatal(err)
		}

		res, err := http.Post(fmt.Sprintf("%s/resolve", os.Args[1]), "application/json", strings.NewReader(string(wBytes)))
		if err != nil {
			log.Fatal(err)
		} else if res.StatusCode != 200 {
			continue
		}

		rBytes, err := ioutil.ReadAll(res.Body)
		if err != nil {
			log.Fatal(err)
		}

		var resolved api.ResolveResponse

		err = json.Unmarshal(rBytes, &resolved)
		if err != nil {
			log.Fatal(err)
		}

		if resolved.CompiledRelease.Version == "" {
			continue
		}

		release["url"] = resolved.CompiledRelease.URL
		release["sha1"] = resolved.CompiledRelease.Sha1

		manifest["releases"].([]interface{})[releaseIdx] = release
	}

	manifestBytes, err := yaml.Marshal(manifest)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println(string(manifestBytes))
}
