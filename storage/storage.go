package storage

import (
	"encoding/json"
	"os"
	"path/filepath"
	"sync"
)

type Storage struct {
	dataDir string
	mutex   sync.RWMutex
}

func NewStorage(dataDir string) (*Storage, error) {
	if err := os.MkdirAll(dataDir, 0755); err != nil {
		return nil, err
	}
	return &Storage{dataDir: dataDir}, nil
}

func (s *Storage) Save(filename string, data interface{}) error {
	s.mutex.Lock()
	defer s.mutex.Unlock()

	path := filepath.Join(s.dataDir, filename)
	file, err := os.Create(path)
	if err != nil {
		return err
	}
	defer file.Close()

	encoder := json.NewEncoder(file)
	encoder.SetIndent("", "  ")
	return encoder.Encode(data)
}

func (s *Storage) Load(filename string, data interface{}) error {
	s.mutex.RLock()
	defer s.mutex.RUnlock()

	path := filepath.Join(s.dataDir, filename)
	file, err := os.Open(path)
	if err != nil {
		if os.IsNotExist(err) {
			return nil // File doesn't exist yet
		}
		return err
	}
	defer file.Close()

	return json.NewDecoder(file).Decode(data)
}
