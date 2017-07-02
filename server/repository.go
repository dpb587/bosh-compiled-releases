package server

type Repository []RepositoryItem

type RepositoryItem struct {
	Name            string
	Version         string
	Digests         []string
	StemcellOS      string
	StemcellVersion string
	CompiledURL     string
	CompiledDigest  string
}
