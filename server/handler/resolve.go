package handler

import (
	"encoding/json"
	"io/ioutil"
	"net/http"

	"github.com/dpb587/bosh-compiled-releases/server/api"
	"github.com/dpb587/bosh-compiled-releases/server/repository"
	"github.com/sirupsen/logrus"
)

type resolve struct {
	repo   repository.Repository
	logger logrus.FieldLogger
}

var _ http.Handler = resolve{}

func NewResolve(logger logrus.FieldLogger, repo repository.Repository) resolve {
	return resolve{
		repo:   repo,
		logger: logger,
	}
}

func (h resolve) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	logger := h.logger.WithFields(logrus.Fields{
		"request.remote_addr": r.RemoteAddr,
		"request.method":      r.Method,
		"request.uri":         r.RequestURI,
	})

	if r.Method != "GET" {
		w.WriteHeader(405)
		w.Write([]byte("method not allowed"))

		logger.WithField("response.status", 405).Warn("method not allowed")

		return
	}

	rBytes, err := ioutil.ReadAll(r.Body)
	if err != nil {
		w.WriteHeader(400)
		w.Write([]byte(err.Error()))

		logger.WithField("response.status", 400).Warn("bad request body: %s", err)

		return
	}

	var data = api.ResolveRequest{}

	err = json.Unmarshal(rBytes, &data)
	if err != nil {
		w.WriteHeader(400)
		w.Write([]byte(err.Error()))

		logger.WithField("response.status", 400).Warnf("bad request body data: %s", err)

		return
	}

	logger = logger.WithFields(logrus.Fields{
		"handler.resolve.release_name":     data.Name,
		"handler.resolve.release_version":  data.Version,
		"handler.resolve.source_digest":    data.Sha1,
		"handler.resolve.stemcell_os":      data.Stemcell.OS,
		"handler.resolve.stemcell_version": data.Stemcell.Version,
	})

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

		logger.WithField("response.status", 500).Errorf("finding compiled releases: %s", err)

		return
	} else if compiledRelease == nil {
		w.WriteHeader(404)
		w.Write([]byte("not found"))

		logger.WithField("response.status", 404).Info("missing compiled release")

		return
	}

	wBytes, err := json.Marshal(api.ResolveResponse{
		CompiledRelease: api.ResolveResponseCompiledRelease{
			URL:  compiledRelease.URL,
			Sha1: compiledRelease.Digest,
		},
	})
	if err != nil {
		w.WriteHeader(500)
		w.Write([]byte(err.Error()))

		logger.WithField("response.status", 500).Error(err)

		return
	}

	w.WriteHeader(200)
	w.Write(wBytes)

	logger.WithField("response.status", 200).Info("found compiled release")
}
