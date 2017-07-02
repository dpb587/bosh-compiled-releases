package record

import (
	"github.com/dpb587/metalink"
)

type Metalink struct {
	Release  MetalinkRelease  `xml:"https://dpb587.github.io/bosh-compiled-release/schema-0.1.0.xsd release,,omitempty" json:"release"`
	Stemcell MetalinkStemcell `xml:"https://dpb587.github.io/bosh-compiled-release/schema-0.1.0.xsd stemcell,,omitempty" json:"stemcell"`
	metalink.Metalink
}

type MetalinkRelease struct {
	Name     string `xml:"name,,omitempty" json:"name"`
	Metalink metalink.Metalink
}

type MetalinkStemcell struct {
	OS      string `xml:"os,,omitempty" json:"os"`
	Version string `xml:"version,,omitempty" json:"version"`
}
