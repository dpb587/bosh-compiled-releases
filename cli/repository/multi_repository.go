package repository

type MultiRepository interface {
	ReadableRepository

	Attach(ReadableRepository)
}

type multiRepository struct {
	repositories []ReadableRepository
}

var _ MultiRepository = &multiRepository{}

func NewMultiRepository() MultiRepository {
	return &multiRepository{
		repositories: []ReadableRepository{},
	}
}

func (r *multiRepository) Find(name string, version string, source SourceRelease, stemcell CompiledReleaseStemcell) (*CompiledReleaseTarball, error) {
	for _, repository := range r.repositories {
		compiledRelease, err := repository.Find(name, version, source, stemcell)
		if err != nil {
			return nil, err
		} else if compiledRelease != nil {
			return compiledRelease, nil
		}
	}

	return nil, nil
}

func (r *multiRepository) List() ([]CompiledRelease, error) {
	var compiledReleases []CompiledRelease

	for _, repository := range r.repositories {
		compiledReleases, err := repository.List()
		if err != nil {
			return nil, err
		}

		for _, compiledRelease := range compiledReleases {
			compiledReleases = append(compiledReleases, compiledRelease)
		}
	}

	return compiledReleases, nil
}

func (r *multiRepository) Attach(repository ReadableRepository) {
	r.repositories = append(r.repositories, repository)
}
