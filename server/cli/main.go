package main

import (
	"log"
	"net/http"
	"os"
	"path/filepath"

	"github.com/dpb587/bosh-compiled-releases/server/handler"
	"github.com/dpb587/bosh-compiled-releases/server/repository"
	"github.com/sirupsen/logrus"
)

func main() {
	logger := logrus.New()
	logger.Level = logrus.DebugLevel
	// logger.Formatter = &logrus.JSONFormatter{}
	logger.Out = os.Stdout

	files, err := filepath.Glob(os.Args[1])
	if err != nil {
		log.Fatal("globbing file repositories: ", err)
	}

	repo := repository.Repository{}

	for _, file := range files {
		err = repository.ImportFileRepository(&repo, file)
		if err != nil {
			log.Fatal("importing file repository: ", err)
		}
	}

	logger.Infof("Repository loaded with %d entries", len(repo))

	http.HandleFunc("/resolve", handler.NewResolve(logger, repo).ServeHTTP)

	logger.WithFields(logrus.Fields{
		"server.local_addr": "127.0.0.1:12345",
	}).Info("Server is ready")

	log.Fatal(http.ListenAndServe(":12345", nil))
}
