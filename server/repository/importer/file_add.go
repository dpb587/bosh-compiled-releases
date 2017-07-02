package main

import (
	"log"
	"os"

	"github.com/dpb587/bosh-compiled-releases/server/repository"
)

func main() {
	if len(os.Args) != 9 {
		log.Fatal("expected args: repository release-name release-version source-digest stemcell-os stemcell-version compiled-digest compiled-url")
	}

	repo := repository.Repository{}

	err := repository.ImportFileRepository(&repo, os.Args[1])
	if err != nil {
		log.Fatal("importing file repository: ", err)
	}

	repo.Add(repository.CompiledRelease{
		Name:    os.Args[2],
		Version: os.Args[3],
		Source: repository.SourceRelease{
			Digest: os.Args[4],
		},
		Stemcell: repository.CompiledReleaseStemcell{
			OS:      os.Args[5],
			Version: os.Args[6],
		},
		Tarball: repository.CompiledReleaseTarball{
			Digest: os.Args[7],
			URL:    os.Args[8],
		},
	})

	err = repository.ExportFileRepository(repo, os.Args[1])
	if err != nil {
		log.Fatal("exporting file repository: ", err)
	}
}
