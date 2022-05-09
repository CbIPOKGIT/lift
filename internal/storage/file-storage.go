package storage

import (
	"encoding/json"
	"log"
	"os"
)

type FileContent map[string]interface{}

type FileStorage struct {
	File    string
	Content FileContent
}

func (st *FileStorage) ReadContent() error {
	content, err := os.ReadFile(st.File)

	if err != nil {
		return err
	}

	if err := json.Unmarshal(content, &st.Content); err != nil {
		st.Content = make(FileContent)
	}
	return nil
}

func (st *FileStorage) Set(key string, val interface{}) error {
	st.Content[key] = val

	if content, err := json.Marshal(st.Content); err == nil {
		return os.WriteFile(st.File, content, 0644)
	} else {
		log.Println("Marshal")
		return err
	}
}

func (st *FileStorage) Get(key string) (string, error) {
	if val, has := st.Content[key]; has {
		if content, err := json.Marshal(val); err == nil {
			return string(content), nil
		} else {
			return "", err
		}
	} else {
		return "", nil
	}
}
