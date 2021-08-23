package db

import (
	"errors"
	"strings"
)

func (s *Storage) Upsert(k string, i interface{}) error {
	err := s.SQL.upsert(k, i, s.Data)
	if err != nil {
		return err
	}

	return s.stateReload()
}

func (s *Storage) GetFirst(k string) (interface{}, error) {
	obj, err := s.SQL.getFirst(k, s.Data)
	if err != nil {
		return nil, err
	}

	return obj, nil
}

func (s *Storage) Get(k string) ([]string, error) {
	obj, err := s.SQL.get(k, s.Data)
	if err != nil {
		return nil, err
	}

	return obj, nil
}

func (s *Storage) GetPath(k string) (interface{}, bool) {
	keys := strings.Split(k, ".")
	obj, exists := s.SQL.getPath(keys, s.Data)
	if exists {
		return obj, exists
	}

	return nil, false
}

func (s *Storage) Delete(k string) error {
	deleted := s.SQL.delPath(k, s.Data)
	if deleted {
		return s.Write()
	}

	return errors.New(KeyDoesNotExist)
}

func (s *Storage) MergeDBs(path string) error {
	err := s.SQL.mergeDBs(path, s.Data)
	if err != nil {
		return err
	}

	return s.Write()
}
