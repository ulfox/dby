package db

import (
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

func (s *Storage) GetPath(k string) (interface{}, error) {
	keys := strings.Split(k, ".")
	obj, err := s.SQL.getPath(keys, s.Data)
	if err != nil {
		return nil, err
	}

	return obj, nil
}

func (s *Storage) Delete(k string) error {
	err := s.SQL.delPath(k, s.Data)
	if err != nil {
		return err
	}

	return s.Write()
}

func (s *Storage) MergeDBs(path string) error {
	err := s.SQL.mergeDBs(path, s.Data)
	if err != nil {
		return err
	}

	return s.Write()
}
