package main

import (
	"encoding/xml"
	"fmt"
	"io/ioutil"
	"log"
	"os"

	"github.com/dpb587/bosh-compiled-releases/record"
	"github.com/dpb587/metalink"
)

func main() {
	if len(os.Args) != 6 {
		log.Fatal("expected args: release source.meta4 compiled.meta4 stemcell-os stemcell-version")
	}

	var rec record.Metalink = record.Metalink{
		Release: record.MetalinkRelease{
			Name: os.Args[1],
		},
		Stemcell: record.MetalinkStemcell{
			OS:      os.Args[4],
			Version: os.Args[5],
		},
	}

	sourceBytes, err := ioutil.ReadFile(os.Args[2])
	if err != nil {
		log.Fatal(err)
	}

	err = metalink.Unmarshal(sourceBytes, &rec.Release.Metalink)
	if err != nil {
		log.Fatal(err)
	}

	compiledBytes, err := ioutil.ReadFile(os.Args[3])
	if err != nil {
		log.Fatal(err)
	}

	err = metalink.Unmarshal(compiledBytes, &rec.Metalink)
	if err != nil {
		log.Fatal(err)
	}

	recBytes, err := xml.MarshalIndent(rec, "", "  ")
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println(string(recBytes))
}
