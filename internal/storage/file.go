package storage

import (
	"encoding/json"
	"os"
	"strconv"
	"sync"
)

type FileRecord struct {
	UUID        string `json:"uuid"`
	ShortURL    string `json:"short_url"`
	OriginalURL string `json:"original_url"`
}

type FileStorage struct {
	mu       sync.RWMutex
	data     map[string]FileRecord
	filePath string
	nextUUID int
}

func NewFileStorage(filePath string) *FileStorage {
	storage := &FileStorage{
		data:     make(map[string]FileRecord),
		filePath: filePath,
		nextUUID: 1,
	}

	storage.load()

	return storage
}

func (s *FileStorage) Save(id, url string) {
	s.mu.Lock()
	defer s.mu.Unlock()

	s.data[id] = FileRecord{
		UUID:        strconv.Itoa(s.nextUUID),
		ShortURL:    id,
		OriginalURL: url,
	}

	s.nextUUID++

	s.save()
}

func (s *FileStorage) Get(id string) (string, bool) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	record, ok := s.data[id]

	if !ok {
		return "", false
	}

	return record.OriginalURL, true
}

func (s *FileStorage) load() {
	data, err := os.ReadFile(s.filePath)
	if err != nil {
		if os.IsNotExist(err) {
			return
		}

		return
	}

	var records []FileRecord

	if err := json.Unmarshal(data, &records); err != nil {
		return
	}

	for _, record := range records {
		s.data[record.ShortURL] = record

		uuid, err := strconv.Atoi(record.UUID)
		if err == nil && uuid >= s.nextUUID {
			s.nextUUID = uuid + 1
		}
	}
}

func (s *FileStorage) save() {
	records := make([]FileRecord, 0, len(s.data))

	for _, record := range s.data {
		records = append(records, record)
	}

	data, err := json.MarshalIndent(records, "", "  ")
	if err != nil {
		return
	}

	_ = os.WriteFile(
		s.filePath,
		data,
		0644,
	)
}
