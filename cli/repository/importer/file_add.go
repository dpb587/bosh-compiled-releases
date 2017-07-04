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

	repo := repository.NewFileRepository(os.Args[1])

	err := repo.Add(repository.CompiledRelease{
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
	if err != nil {
		log.Fatal("adding compiled release: ", err)
	}

	err = repo.Save()
	if err != nil {
		log.Fatal("exporting file repository: ", err)
	}
}
