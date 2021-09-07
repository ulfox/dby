package db

import (
	"fmt"
	"strings"
)

// Upsert is a SQL wrapper for adding/updating map structures
func (s *Storage) Upsert(k string, i interface{}) error {
	err := s.SQL.upsertRecursive(strings.Split(k, "."), s.Data[s.AD], i)
	if err != nil {
		return wrapErr(err, getFn())
	}

	return s.stateReload()
}

// UpsertGlobal is a SQL wrapper for adding/updating map structures
// in all documents. This will change all existing paths to the given
// structure and add new if the path is missing for a document
func (s *Storage) UpsertGlobal(k string, i interface{}) error {
	c := s.AD
	for j := range s.Data {
		err := s.SQL.upsertRecursive(strings.Split(k, "."), s.Data[j], i)
		if err != nil {
			return wrapErr(err, getFn())
		}
	}

	s.AD = c

	return s.stateReload()
}

// UpdateGlobal is a SQL wrapper for adding/updating map structures
// in all documents. This will change all existing paths to the given
// structure and add new if the path is missing for a document
func (s *Storage) UpdateGlobal(k string, i interface{}) error {
	c := s.AD
	for j := range s.Data {
		s.AD = j

		if _, err := s.GetPath(k); err != nil {
			continue
		}

		err := s.SQL.upsertRecursive(strings.Split(k, "."), s.Data[s.AD], i)
		if err != nil {
			return wrapErr(err, getFn())
		}
	}

	s.AD = c

	return s.stateReload()
}

// GetFirst is a SQL wrapper for finding the first key in the
// yaml hierarchy. If two keys are on the same level but under
// different paths, then the selection will be random
func (s *Storage) GetFirst(k string) (interface{}, error) {
	obj, err := s.SQL.getFirst(k, s.Data[s.AD])
	if err != nil {
		return nil, wrapErr(err, getFn())
	}

	return obj, nil
}

// GetFirstGlobal does the same as GetFirst but for all docs.
// Instead of returning an interface it returns a map with keys
// the index of the doc that a key was found and value the value of the key
func (s *Storage) GetFirstGlobal(k string) map[int]interface{} {
	found := make(map[int]interface{})

	c := s.AD
	for j := range s.Data {
		s.AD = j

		obj, err := s.SQL.getFirst(k, s.Data[s.AD])
		if err != nil {
			continue
		}
		found[s.AD] = obj
	}

	s.AD = c

	return found
}

// Get is a SQL wrapper that finds all the paths for a given
// e.g. ["key-1.test", "key-2.key-3.test"] will be returned
// if "test" was the key asked from the following yaml
// ---------
//
// key-1:
//		test: someValue-1
// key-2:
//		key-3:
//			test: someValue-2
//
func (s *Storage) Get(k string) ([]string, error) {
	fmt.Println("Warn: Deprecated is Get(). Will be replaced by FindKeys() in the future.")
	obj, err := s.SQL.get(k, s.Data[s.AD])
	if err != nil {
		return nil, wrapErr(err, getFn())
	}

	return obj, nil
}

// FindKeys is alias of Get. This function will replace
// Get in the future since this name for finding keys
// makes more sense
// For now we keep both for compatibility
func (s *Storage) FindKeys(k string) ([]string, error) {
	obj, err := s.SQL.get(k, s.Data[s.AD])
	if err != nil {
		return nil, wrapErr(err, getFn())
	}

	return obj, nil
}

// FindKeysGlobal does the same as FindKeys but for all docs.
// Instead of returning a list of keys it returns a map with indexes
// from the docs and value an array of paths that was found
func (s *Storage) FindKeysGlobal(k string) map[int][]string {
	found := make(map[int][]string)

	c := s.AD
	for j := range s.Data {
		s.AD = j

		obj, err := s.SQL.get(k, s.Data[s.AD])
		if err != nil || len(obj) == 0 {
			continue
		}
		found[s.AD] = obj
	}

	s.AD = c

	return found
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
	obj, err := s.SQL.getPath(keys, s.Data[s.AD])
	if err != nil {
		return nil, wrapErr(err, getFn())
	}

	return obj, nil
}

// GetPathGlobal does the same as GetPath but globally for all
// docs
func (s *Storage) GetPathGlobal(k string) map[int]interface{} {
	found := make(map[int]interface{})
	keys := strings.Split(k, ".")

	c := s.AD
	for j := range s.Data {
		s.AD = j

		obj, err := s.SQL.getPath(keys, s.Data[s.AD])
		if err != nil {
			continue
		}
		found[s.AD] = obj
	}

	s.AD = c

	return found
}

// Delete is a SQL wrapper that deletes the last key from a given
// path. For example, Delete("key-1.key-2.key-3") would first
// validate that the path exists, then it would export the value of
// GetPath("key-1.key-2") and delete the object that matches key-3
func (s *Storage) Delete(k string) error {
	err := s.SQL.delPath(k, s.Data[s.AD])
	if err != nil {
		return wrapErr(err, getFn())
	}

	return s.Write()
}

// DeleteGlobal is the same as Delete but will try to delete
// the path on all docs (if found)
func (s *Storage) DeleteGlobal(k string) error {
	err := s.SQL.delPath(k, s.Data[s.AD])
	if err != nil {
		return wrapErr(err, getFn())
	}

	return s.Write()
}

// MergeDBs is a SQL wrapper that merges a source yaml file
// with the DBy local yaml file.
func (s *Storage) MergeDBs(path string) error {
	err := s.SQL.mergeDBs(path, s.Data[s.AD])
	if err != nil {
		return wrapErr(err, getFn())
	}

	return s.Write()
}
