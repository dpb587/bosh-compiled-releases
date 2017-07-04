package cmd

import (
	"fmt"
	"io/ioutil"
	"log"
	"path/filepath"
	"strconv"
	"strings"

	yaml "gopkg.in/yaml.v2"

	"github.com/dpb587/bosh-compiled-releases/cli/repository"
	"github.com/jessevdk/go-flags"
)

type RewriteManifest struct {
	Server []string `long:"server" description:"Remote server to query"`
	Local  []string `long:"local" description:"Local path to query"`

	Args RewriteManifestArgs `positional-args:"true"`
}

var _ flags.Commander = RewriteManifest{}

type RewriteManifestArgs struct {
	Manifest string `positional-arg-name:"MANIFEST-PATH" description:"Manifest path to parse"`
}

func (c RewriteManifest) Execute(args []string) error {
	repo := repository.NewMultiRepository()

	for _, local := range c.Local {
		files, err := filepath.Glob(local)
		if err != nil {
			log.Fatal("globbing file repositories: ", err)
		}

		for _, file := range files {
			repo.Attach(repository.NewFileRepository(file))
		}
	}

	for _, server := range c.Server {
		repo.Attach(repository.NewServerRepository(server))
	}

	// start

	var manifest map[string]interface{}

	bytes, err := ioutil.ReadFile(c.Args.Manifest)
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

		compiledRelease, err := repo.Find(
			release["name"].(string),
			version,
			repository.SourceRelease{
				Digest: release["sha1"].(string),
			},
			repository.CompiledReleaseStemcell{
				OS:      stemcellOS,
				Version: stemcellVersion,
			},
		)
		if err != nil {
			log.Fatal(err)
		} else if compiledRelease == nil {
			continue
		}

		release["url"] = compiledRelease.URL
		release["sha1"] = compiledRelease.Digest

		manifest["releases"].([]interface{})[releaseIdx] = release
	}

	manifestBytes, err := yaml.Marshal(manifest)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println(string(manifestBytes))

	return nil
}
