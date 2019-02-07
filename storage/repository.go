package storage

import (
	"sync"
)

type Storage struct {
	sync.RWMutex

	Data map[string]string
}

func NewRepository(storage *Storage) Repository {
	return &MemoryStorageRepository{
		Storage: storage,
	}
}

type Repository interface {
	Get(key string) (string, error)
	Set(key string, value string) error
	GetAll() map[string]string
}

type MemoryStorageRepository struct {
	Storage *Storage
}

func (r *MemoryStorageRepository) Get(key string) (string, error) {
	result, _ := r.Storage.Data[key]
	return result, nil
}

func (r *MemoryStorageRepository) Set(key string, value string) error {
	r.Storage.Lock()
	defer r.Storage.Unlock()
	r.Storage.Data[key] = value
	return nil
}

func (r *MemoryStorageRepository) GetAll() map[string]string {
	r.Storage.RLock()
	defer r.Storage.RUnlock()
	return r.Storage.Data
}
