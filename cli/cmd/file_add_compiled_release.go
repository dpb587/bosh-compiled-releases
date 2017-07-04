package cmd

import (
	"log"

	"github.com/dpb587/bosh-compiled-releases/cli/repository"
	"github.com/jessevdk/go-flags"
)

type FileAddCompiledRelease struct {
	Args fileAddCompiledReleaseArgs `positional-args:"true" required:"true"`
}

var _ flags.Commander = FileAddCompiledRelease{}

type fileAddCompiledReleaseArgs struct {
	// repository release-name release-version source-digest stemcell-os stemcell-version compiled-digest compiled-url
	Repository      string `positional-arg-name:"FILE" description:"Local JSON index file"`
	ReleaseName     string `positional-arg-name:"RELEASE-NAME" description:"Release name"`
	ReleaseVersion  string `positional-arg-name:"RELEASE-VERSION" description:"Release version"`
	SourceDigest    string `positional-arg-name:"SOURCE-DIGEST" description:"Source digest/sha1"`
	StemcellOS      string `positional-arg-name:"STEMCELL-OS" description:"Stemcell OS (e.g. ubuntu-trusty)"`
	StemcellVersion string `positional-arg-name:"STEMCELL-VERSION" description:"Stemcell Version (e.g. 3421.11)"`
	CompiledDigest  string `positional-arg-name:"COMPILED-DIGEST" description:"Compiled digest/sha1"`
	CompiledURL     string `positional-arg-name:"COMPILED-URL" description:"Compiled tarball URL"`
}

func (c FileAddCompiledRelease) Execute(args []string) error {
	repo := repository.NewFileRepository(c.Args.Repository)

	err := repo.Add(repository.CompiledRelease{
		Name:    c.Args.ReleaseName,
		Version: c.Args.ReleaseVersion,
		Source: repository.SourceRelease{
			Digest: c.Args.SourceDigest,
		},
		Stemcell: repository.CompiledReleaseStemcell{
			OS:      c.Args.StemcellOS,
			Version: c.Args.StemcellVersion,
		},
		Tarball: repository.CompiledReleaseTarball{
			Digest: c.Args.CompiledDigest,
			URL:    c.Args.CompiledURL,
		},
	})
	if err != nil {
		log.Fatal("adding compiled release: ", err)
	}

	err = repo.Save()
	if err != nil {
		log.Fatal("saving file repository: ", err)
	}

	return nil
}
