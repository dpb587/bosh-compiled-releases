package main

import (
	"fmt"
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
	logger.Formatter = &logrus.JSONFormatter{}
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

	logger.Infof("loaded %d compiled releases", len(repo))

	http.HandleFunc("/resolve", handler.NewResolve(logger, repo).ServeHTTP)

	bind := fmt.Sprintf("%s:%s", "0.0.0.0", os.Getenv("PORT"))

	logger.WithFields(logrus.Fields{
		"server.local_addr": bind,
	}).Info("server is ready")

	log.Fatal(http.ListenAndServe(bind, nil))
}
