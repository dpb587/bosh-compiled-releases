package repository

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
