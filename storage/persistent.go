package storage

import (
	"errors"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"sync"

	"github.com/google/uuid"
)

type PersistentStore struct {
	sync.Mutex
	uploadPath string
}

func NewPersistentStore(uploadPath string) (*PersistentStore, error) {
	if err := os.MkdirAll(uploadPath, 0700); err == nil {
		return &PersistentStore{uploadPath: uploadPath}, nil
	} else {
		return nil, errors.New("unable to create base folder to store uploaded files")
	}
}

func (p *PersistentStore) Upload(fileContent []byte) (string, error) {
	p.Mutex.Lock()
	defer p.Mutex.Unlock()
	fileName := uuid.New().String()
	imagePath := filepath.Join(p.uploadPath, fileName+".zip")

	newFile, err := os.Create(imagePath)
	if err != nil {
		return "", errors.New("unable to write file to persistent storage")
	}
	defer newFile.Close()
	if _, err := newFile.Write(fileContent); err != nil || newFile.Close() != nil {
		return "", errors.New("unable to write file to persistent storage")
	}
	log.Println("UPLOADING:", imagePath)
	return fileName, nil
}

func (p *PersistentStore) Download(fileName string) ([]byte, error) {
	p.Mutex.Lock()
	defer p.Mutex.Unlock()
	if _, err := os.Stat(filepath.Join(p.uploadPath, fileName) + ".zip"); os.IsNotExist(err) {
		return nil, errors.New("file not found")
	}
	if fc, err := ioutil.ReadFile(filepath.Join(p.uploadPath, fileName) + ".zip"); err != nil {
		return nil, errors.New("cannot read file")
	} else {
		return fc, nil
	}
}
