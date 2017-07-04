package api

type ResolveRequest struct {
	Name     string                 `json:"name"`
	Version  string                 `json:"version"`
	Sha1     string                 `json:"sha1"`
	Stemcell ResolveRequestStemcell `json:"stemcell"`
}

type ResolveRequestStemcell struct {
	OS      string `json:"os"`
	Version string `json:"version"`
}

type ResolveResponse struct {
	CompiledRelease ResolveResponseCompiledRelease `json:"compiled_release"`
}

type ResolveResponseCompiledRelease struct {
	Sha1 string `json:"sha1"`
	URL  string `json:"url"`
}
