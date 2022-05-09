package storage

import (
	"os"
)

const STORAGE_FILE = "storage.bin"

type Storage interface {
	// Записуємо значення з ключом key і значенням val
	Set(key string, val interface{}) error

	// Знаходимо значення по ключу
	Get(key string) (string, error)
}

func New() Storage {
	storage := new(FileStorage)
	storage.File = STORAGE_FILE

	_, err := os.Stat(storage.File)
	if err != nil && os.IsNotExist(err) {
		os.Create(storage.File)
	}

	storage.ReadContent()
	return storage
}
