package db

import (
	"strings"
)

// Upsert is a SQL wrapper for adding/updating map structures
func (s *Storage) Upsert(k string, i interface{}) error {
	err := s.SQL.upsert(k, i, s.Data)
	if err != nil {
		return err
	}

	return s.stateReload()
}

// GetFirst is a SQL wrapper for finding the first key in the
// yaml hierarchy. If two keys are on the same level but under
// different paths, then the selection will be random
func (s *Storage) GetFirst(k string) (interface{}, error) {
	obj, err := s.SQL.getFirst(k, s.Data)
	if err != nil {
		return nil, err
	}

	return obj, nil
}

// Get is a SQL wrapper that finds all the paths for a given
// e.g. ["key-1.test", "key-2.key-3.test"] will be returned
// if "test" was the key asked from the following yaml
// ---------
// key-1:
//	test: someValue-1
// key-2:
//	key-3:
//		test: someValue-2
//
func (s *Storage) Get(k string) ([]string, error) {
	obj, err := s.SQL.get(k, s.Data)
	if err != nil {
		return nil, err
	}

	return obj, nil
}

// GetPath is a SQL wrapper that returns the value for a given
// path. Example, it would return "value-1" if "key-1.key-2" was
// the path asked from the following yaml
// ---------
// key-1:
//	key-2: value-1
//
func (s *Storage) GetPath(k string) (interface{}, error) {
	keys := strings.Split(k, ".")
	obj, err := s.SQL.getPath(keys, s.Data)
	if err != nil {
		return nil, err
	}

	return obj, nil
}

// Delete is a SQL wrapper that deletes the last key from a given
// path. For example, Delete("key-1.key-2.key-3") would first
// validate that the path exists, then it would export the value of
// GetPath("key-1.key-2") and delete the object that matches key-3
func (s *Storage) Delete(k string) error {
	err := s.SQL.delPath(k, s.Data)
	if err != nil {
		return err
	}

	return s.Write()
}

// MergeDBs is a SQL wrapper that merges a source yaml file
// with the DBy local yaml file.
func (s *Storage) MergeDBs(path string) error {
	err := s.SQL.mergeDBs(path, s.Data)
	if err != nil {
		return err
	}

	return s.Write()
}
