package storage

import (
	"errors"
	"sync"

	"github.com/google/uuid"
)

type VolatileStore struct {
	sync.Mutex
	store map[string][]byte
}

func NewVolatileStore() *VolatileStore {
	return &VolatileStore{store: make(map[string][]byte)}
}

func (v *VolatileStore) Upload(fileContent []byte) (string, error) {
	v.Mutex.Lock()
	defer v.Mutex.Unlock()
	fileName := uuid.New().String()
	v.store[fileName] = fileContent
	return fileName, nil
}

func (v *VolatileStore) Download(fileName string) ([]byte, error) {
	v.Mutex.Lock()
	defer v.Mutex.Unlock()
	if fc, ok := v.store[fileName]; ok {
		return fc, nil
	} else {
		return nil, errors.New("file not found")
	}
}
