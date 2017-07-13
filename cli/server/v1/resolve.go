package v1

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/dpb587/bosh-compiled-releases/cli/repository"
	"github.com/dpb587/bosh-compiled-releases/cli/server/v1/model"
	"github.com/sirupsen/logrus"
)

type resolve struct {
	repo   repository.ReadableRepository
	logger logrus.FieldLogger
}

var _ http.Handler = resolve{}

func NewResolve(logger logrus.FieldLogger, repo repository.ReadableRepository) http.Handler {
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
		w.Header().Add("Allow", "GET")
		http.Error(w, fmt.Sprintf("405 Method Not Allowed"), http.StatusMethodNotAllowed)
		logger.WithField("response.status", http.StatusMethodNotAllowed).Warn("method not allowed")

		return
	}

	rBytes, err := ioutil.ReadAll(r.Body)
	if err != nil {
		http.Error(w, fmt.Sprintf("400 Bad Request\n%s", err), http.StatusBadRequest)
		logger.WithField("response.status", http.StatusBadRequest).Warn("bad request body: %s", err)

		return
	}

	var data = model.ResolveRequest{}

	err = json.Unmarshal(rBytes, &data)
	if err != nil {
		http.Error(w, fmt.Sprintf("400 Bad Request\n%s", err), http.StatusBadRequest)
		logger.WithField("response.status", http.StatusBadRequest).Warnf("bad request body data: %s", err)

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
		http.Error(w, fmt.Sprintf("500 Internal Server Error"), http.StatusInternalServerError)
		logger.WithField("response.status", http.StatusInternalServerError).Errorf("finding compiled release: %s", err)

		return
	} else if compiledRelease == nil {
		http.Error(w, fmt.Sprintf("404 Not Found"), http.StatusNotFound)
		logger.WithField("response.status", http.StatusNotFound).Info("missing compiled release")

		return
	}

	wBytes, err := json.Marshal(model.ResolveResponse{
		CompiledRelease: model.ResolveResponseCompiledRelease{
			URL:  compiledRelease.URL,
			Sha1: compiledRelease.Digest,
		},
	})
	if err != nil {
		http.Error(w, fmt.Sprintf("500 Internal Server Error"), http.StatusInternalServerError)
		logger.WithField("response.status", http.StatusInternalServerError).Error(err)

		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write(wBytes)

	logger.WithField("response.status", http.StatusOK).Info("found compiled release")
}
