package repository

import (
	"bufio"
	"encoding/json"
	"log"
	"os"
	"path/filepath"
)

func ImportFileRepository(repo *Repository, path string) error {
	if _, err := os.Stat(path); os.IsNotExist(err) {
		return nil
	}

	file, err := os.Open(path)
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

		repo.Add(compiledRelease)
	}
	if err := scanner.Err(); err != nil {
		log.Panic("scanning file repository: ", err)
	}

	return nil
}

func ExportFileRepository(repo Repository, path string) error {
	err := os.MkdirAll(filepath.Dir(path), 0755)
	if err != nil {
		log.Panic("creating file repository directory: ", err)
	}

	file, err := os.Create(path)
	if err != nil {
		log.Panic("opening file repository: ", err)
	}

	defer file.Close()

	for _, compiledRelease := range repo {
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
