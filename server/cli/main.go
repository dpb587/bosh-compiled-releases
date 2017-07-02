package main

import (
	"log"
	"net/http"
	"os"
	"path/filepath"

	"github.com/dpb587/bosh-compiled-releases/server/handler"
	"github.com/dpb587/bosh-compiled-releases/server/repository"
)

func main() {
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

	http.HandleFunc("/resolve", handler.NewResolve(repo).ServeHTTP)

	log.Fatal(http.ListenAndServe(":12345", nil))
}
