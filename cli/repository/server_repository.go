package repository

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"strings"

	v1 "github.com/dpb587/bosh-compiled-releases/cli/server/v1/model"
)

type ServerRepository interface {
	ReadableRepository
}

type serverRepository struct {
	endpoint string
}

var _ ServerRepository = &serverRepository{}

func NewServerRepository(endpoint string) ServerRepository {
	return &serverRepository{
		endpoint: endpoint,
	}
}

func (r *serverRepository) Find(name string, version string, source SourceRelease, stemcell CompiledReleaseStemcell) (*CompiledReleaseTarball, error) {
	wBytes, err := json.Marshal(v1.ResolveRequest{
		Name:    name,
		Version: version,
		Sha1:    source.Digest,
		Stemcell: v1.ResolveRequestStemcell{
			OS:      stemcell.OS,
			Version: stemcell.Version,
		},
	})
	if err != nil {
		log.Fatal(err)
	}

	req, err := http.NewRequest("GET", fmt.Sprintf("%s/v1/resolve", r.endpoint), strings.NewReader(string(wBytes)))
	if err != nil {
		log.Fatal(err)
	}

	res, err := http.DefaultClient.Do(req)
	if err != nil {
		log.Fatal(err)
	} else if res.StatusCode != 200 {
		return nil, nil
	}

	rBytes, err := ioutil.ReadAll(res.Body)
	if err != nil {
		log.Fatal(err)
	}

	var resolved v1.ResolveResponse

	err = json.Unmarshal(rBytes, &resolved)
	if err != nil {
		log.Fatal(err)
	}

	return &CompiledReleaseTarball{
		Digest: resolved.CompiledRelease.Sha1,
		URL:    resolved.CompiledRelease.URL,
	}, nil
}

func (r *serverRepository) List() ([]CompiledRelease, error) {
	return nil, nil
}
