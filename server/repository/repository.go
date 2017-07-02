package repository

type Repository []CompiledRelease

func (r Repository) Find(name string, version string, source SourceRelease, stemcell CompiledReleaseStemcell) (*CompiledReleaseTarball, error) {
	for _, compiled := range r {
		if compiled.Name != name || compiled.Version != version {
			continue
		} else if compiled.Stemcell.OS != stemcell.OS || compiled.Stemcell.Version != stemcell.Version {
			continue
		} else if compiled.Source.Digest != source.Digest {
			continue
		}

		return &compiled.Tarball, nil
	}

	return nil, nil
}

func (r *Repository) Add(compiledRelease CompiledRelease) {
	if found, _ := r.Find(compiledRelease.Name, compiledRelease.Version, compiledRelease.Source, compiledRelease.Stemcell); found != nil {
		return
	}

	*r = append(*r, compiledRelease)
}

type SourceRelease struct {
	Digest string `json:"digest"`
}

type CompiledRelease struct {
	Name    string `json:"name"`
	Version string `json:"version"`

	Source   SourceRelease           `json:"source"`
	Stemcell CompiledReleaseStemcell `json:"stemcell"`
	Tarball  CompiledReleaseTarball  `json:"tarball"`
}

type CompiledReleaseStemcell struct {
	OS      string `json:"os"`
	Version string `json:"version"`
}

type CompiledReleaseTarball struct {
	Digest string `json:"digest"`
	URL    string `json:"url"`
}
