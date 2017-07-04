package handler

import (
	"net/http"

	"github.com/sirupsen/logrus"
)

type staticAsset struct {
	asset  string
	logger logrus.FieldLogger
}

var _ http.Handler = staticAsset{}

func NewStaticAsset(logger logrus.FieldLogger, asset string) staticAsset {
	return staticAsset{
		asset:  asset,
		logger: logger,
	}
}

func (h staticAsset) ServeHTTP(w http.ResponseWriter, r *http.Request) {
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

	http.ServeFile(w, r, h.asset)

	logger.WithField("response.status", 200).Info("served static asset")
}
