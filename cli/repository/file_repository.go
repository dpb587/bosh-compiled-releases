package repository

import (
	"bufio"
	"encoding/json"
	"log"
	"os"
	"path/filepath"
)

type FileRepository interface {
	Repository

	Save() error
}

type fileRepository struct {
	path       string
	repository InMemoryRepository
}

var _ FileRepository = &fileRepository{}

func NewFileRepository(path string) FileRepository {
	return &fileRepository{
		path: path,
	}
}

func (r *fileRepository) Find(name string, version string, source SourceRelease, stemcell CompiledReleaseStemcell) (*CompiledReleaseTarball, error) {
	if err := r.requireRepository(); err != nil {
		return nil, err
	}

	return r.repository.Find(name, version, source, stemcell)
}

func (r *fileRepository) Add(compiledRelease CompiledRelease) error {
	if err := r.requireRepository(); err != nil {
		return err
	}

	return r.repository.Add(compiledRelease)
}

func (r *fileRepository) List() ([]CompiledRelease, error) {
	if err := r.requireRepository(); err != nil {
		return nil, err
	}

	return r.repository.List()
}

func (r *fileRepository) requireRepository() error {
	if r.repository != nil {
		return nil
	}

	r.repository = NewInMemoryRepository()

	if _, err := os.Stat(r.path); os.IsNotExist(err) {
		return nil
	}

	file, err := os.Open(r.path)
	if err != nil {
		log.Panic("opening file repository: ", err)
	}

	defer file.Close()

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		var compiledRelease CompiledRelease

		err := json.Unmarshal(scanner.Bytes(), &compiledRelease)
		if err != nil {
			log.Panic("unmarshalling json: ", err)
		}

		r.repository.Add(compiledRelease)
	}
	if err := scanner.Err(); err != nil {
		log.Panic("scanning file repository: ", err)
	}

	return nil
}

func (r *fileRepository) Save() error {
	if r.repository == nil {
		return nil
	}

	err := os.MkdirAll(filepath.Dir(r.path), 0755)
	if err != nil {
		log.Panic("creating file repository directory: ", err)
	}

	file, err := os.Create(r.path)
	if err != nil {
		log.Panic("opening file repository: ", err)
	}

	defer file.Close()

	compiledReleases, err := r.repository.List()
	if err != nil {
		log.Panic("listing compiled releases: ", err)
	}

	for _, compiledRelease := range compiledReleases {
		bytes, err := json.Marshal(compiledRelease)
		if err != nil {
			log.Panic("marshalling json: ", err)
		}

		bytes = append(bytes, byte('\n'))

		_, err = file.Write(bytes)
		if err != nil {
			log.Panic("writing file repository: ", err)
		}
	}

	return nil
}
