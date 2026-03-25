package service

import (
	"encoding/csv"
	"errors"
	"io"
	"os"
	"strings"
	"sync"
)

type Taxonomy struct {
	mu   sync.RWMutex
	data map[string]string
}

func NewTaxonomy() *Taxonomy {
	return &Taxonomy{
		data: make(map[string]string),
	}
}

func (t *Taxonomy) LoadCSV(path string) error {
	f, err := os.Open(path)
	if err != nil {
		return err
	}
	defer f.Close()

	r := csv.NewReader(f)
	r.FieldsPerRecord = -1

	header, err := r.Read()
	if err != nil {
		return err
	}
	if len(header) < 2 {
		return errors.New("invalid csv header")
	}

	tmp := make(map[string]string)

	for {
		record, err := r.Read()
		if errors.Is(err, io.EOF) {
			break
		}
		if err != nil {
			return err
		}
		if len(record) < 2 {
			continue
		}

		id := strings.TrimSpace(record[0])
		name := strings.TrimSpace(record[1])

		if id == "" || name == "" {
			continue
		}

		tmp[id] = name
	}

	t.mu.Lock()
	t.data = tmp
	t.mu.Unlock()

	return nil
}

func (t *Taxonomy) Get(id string) (string, bool) {
	t.mu.RLock()
	defer t.mu.RUnlock()

	name, ok := t.data[id]
	return name, ok
}

func (t *Taxonomy) Count() int {
	t.mu.RLock()
	defer t.mu.RUnlock()

	return len(t.data)
}
