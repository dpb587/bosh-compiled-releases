package main

import (
	"encoding/json"
	"encoding/xml"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"path/filepath"

	"github.com/dpb587/bosh-compiled-releases/api"
	"github.com/dpb587/bosh-compiled-releases/record"
	"github.com/dpb587/bosh-compiled-releases/server"
)

func main() {
	files, err := filepath.Glob(os.Args[1])
	if err != nil {
		log.Fatal(err)
	}

	repo := server.Repository{}

	for _, file := range files {
		var rec record.Metalink

		log.Print(file)

		recBytes, err := ioutil.ReadFile(file)
		if err != nil {
			log.Fatal(err)
		}

		err = xml.Unmarshal(recBytes, &rec)
		if err != nil {
			log.Fatal(err)
		}

		item := server.RepositoryItem{
			Name:            rec.Release.Name,
			Version:         rec.Release.Metalink.Files[0].Version,
			Digests:         []string{},
			StemcellOS:      rec.Stemcell.OS,
			StemcellVersion: rec.Stemcell.Version,
			CompiledURL:     rec.Files[0].URLs[0].URL,
		}

		for _, hash := range rec.Release.Metalink.Files[0].Hashes {
			if hash.Type == "sha-1" {
				item.Digests = append(item.Digests, hash.Hash)
			}
		}

		for _, hash := range rec.Metalink.Files[0].Hashes {
			if hash.Type == "sha-1" {
				item.CompiledDigest = hash.Hash
			}
		}

		repo = append(repo, item)
	}

	http.HandleFunc("/resolve", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != "POST" {
			w.WriteHeader(405)
			w.Write([]byte("method not allowed"))

			return
		}

		rBytes, err := ioutil.ReadAll(r.Body)
		if err != nil {
			w.WriteHeader(400)
			w.Write([]byte(err.Error()))

			return
		}

		var data = api.ResolveRequest{}

		err = json.Unmarshal(rBytes, &data)
		if err != nil {
			w.WriteHeader(400)
			w.Write([]byte(err.Error()))

			return
		}

		log.Printf("%#+v", data)

		for _, item := range repo {
			if item.Name != data.Name || item.Version != data.Version {
				continue
			} else if item.StemcellOS != data.Stemcell.OS || item.StemcellVersion != data.Stemcell.Version {
				continue
			}

			for _, digest := range item.Digests {
				if digest == data.Sha1 {
					wBytes, err := json.Marshal(api.ResolveResponse{
						CompiledRelease: api.ResolveResponseCompiledRelease{
							Name:    item.Name,
							Version: item.Version,
							URL:     item.CompiledURL,
							Sha1:    item.CompiledDigest,
						},
					})
					if err != nil {
						w.WriteHeader(500)
						w.Write([]byte(err.Error()))

						return
					}

					w.WriteHeader(200)
					w.Write(wBytes)

					return
				}
			}
		}

		w.WriteHeader(404)
		w.Write([]byte("not found"))
	})

	log.Fatal(http.ListenAndServe(":12345", nil))
}
