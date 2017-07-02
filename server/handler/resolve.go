package handler

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"

	"github.com/dpb587/bosh-compiled-releases/server/api"
	"github.com/dpb587/bosh-compiled-releases/server/repository"
)

type resolve struct {
	repo repository.Repository
}

var _ http.Handler = resolve{}

func NewResolve(repo repository.Repository) resolve {
	return resolve{
		repo: repo,
	}
}

func (h resolve) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		w.WriteHeader(405)
		w.Write([]byte("method not allowed"))

		return
	}

	rBytes, err := ioutil.ReadAll(r.Body)
	if err != nil {
		w.WriteHeader(400)
		w.Write([]byte(err.Error()))

		return
	}

	var data = api.ResolveRequest{}

	err = json.Unmarshal(rBytes, &data)
	if err != nil {
		w.WriteHeader(400)
		w.Write([]byte(err.Error()))

		return
	}

	log.Printf("%#+v", data)

	compiledRelease, err := h.repo.Find(
		data.Name,
		data.Version,
		repository.SourceRelease{
			Digest: data.Sha1,
		},
		repository.CompiledReleaseStemcell{
			OS:      data.Stemcell.OS,
			Version: data.Stemcell.Version,
		},
	)
	if err != nil {
		w.WriteHeader(500)
		w.Write([]byte(err.Error()))

		return
	} else if compiledRelease == nil {
		w.WriteHeader(404)
		w.Write([]byte("not found"))

		return
	}

	wBytes, err := json.Marshal(api.ResolveResponse{
		CompiledRelease: api.ResolveResponseCompiledRelease{
			Name:    data.Name,
			Version: data.Version,
			URL:     compiledRelease.URL,
			Sha1:    compiledRelease.Digest,
		},
	})
	if err != nil {
		w.WriteHeader(500)
		w.Write([]byte(err.Error()))

		return
	}

	w.WriteHeader(200)
	w.Write(wBytes)
}
