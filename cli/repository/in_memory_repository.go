package repository

type InMemoryRepository interface {
	Repository
}

type inMemoryRepository struct {
	compiledReleases []CompiledRelease
}

var _ InMemoryRepository = &inMemoryRepository{}

func NewInMemoryRepository() InMemoryRepository {
	return &inMemoryRepository{
		compiledReleases: []CompiledRelease{},
	}
}

func (r *inMemoryRepository) Find(name string, version string, source SourceRelease, stemcell CompiledReleaseStemcell) (*CompiledReleaseTarball, error) {
	for _, compiledReleases := range r.compiledReleases {
		if compiledReleases.Name != name || compiledReleases.Version != version {
			continue
		} else if compiledReleases.Stemcell.OS != stemcell.OS || compiledReleases.Stemcell.Version != stemcell.Version {
			continue
		} else if compiledReleases.Source.Digest != source.Digest {
			continue
		}

		return &compiledReleases.Tarball, nil
	}

	return nil, nil
}

func (r *inMemoryRepository) Add(compiledRelease CompiledRelease) error {
	if found, _ := r.Find(compiledRelease.Name, compiledRelease.Version, compiledRelease.Source, compiledRelease.Stemcell); found != nil {
		return nil
	}

	r.compiledReleases = append(r.compiledReleases, compiledRelease)

	return nil
}

func (r *inMemoryRepository) List() ([]CompiledRelease, error) {
	return r.compiledReleases, nil
}
