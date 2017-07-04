package repository

type Repository interface {
	ReadableRepository
	WritableRepository
}

type ReadableRepository interface {
	Find(name string, version string, source SourceRelease, stemcell CompiledReleaseStemcell) (*CompiledReleaseTarball, error)
	List() ([]CompiledRelease, error)
}

type WritableRepository interface {
	Add(compiledRelease CompiledRelease) error
}
