package db

import (
	"strings"
)

// Upsert is a SQL wrapper for adding/updating map structures
func (s *Storage) Upsert(k string, i interface{}) error {
	data, err := s.SQL.toInterfaceMap(&i)
	if err != nil {
		return wrapErr(err)
	}
	err = s.SQL.upsertRecursive(strings.Split(k, "."), s.getData(), data)
	if err != nil {
		return wrapErr(err)
	}

	return s.stateReload()
}

// UpsertGlobal is a SQL wrapper for adding/updating map structures
// in all documents. This will change all existing paths to the given
// structure and add new if the path is missing for a document
func (s *Storage) UpsertGlobal(k string, i interface{}) error {
	data, err := s.SQL.toInterfaceMap(&i)
	if err != nil {
		return wrapErr(err)
	}

	c := s.GetAD()
	for j := range s.GetAllData() {
		s.SetAD(j)
		err := s.SQL.upsertRecursive(strings.Split(k, "."), s.getData(), data)
		if err != nil {
			return wrapErr(err)
		}
	}

	s.SetAD(c)

	return s.stateReload()
}

// UpdateGlobal is a SQL wrapper for adding/updating map structures
// in all documents. This will change all existing paths to the given
// structure (if any)
func (s *Storage) UpdateGlobal(k string, i interface{}) error {
	data, err := s.SQL.toInterfaceMap(&i)
	if err != nil {
		return wrapErr(err)
	}

	c := s.GetAD()
	for j := range s.GetAllData() {
		s.SetAD(j)
		if _, err := s.GetPath(k); err != nil {
			continue
		}

		err := s.SQL.upsertRecursive(strings.Split(k, "."), s.getData(), data)
		if err != nil {
			return wrapErr(err)
		}
	}

	s.SetAD(c)

	return s.stateReload()
}

// GetFirst is a SQL wrapper for finding the first key in the
// yaml hierarchy. If two keys are on the same level but under
// different paths, then the selection will be random
func (s *Storage) GetFirst(k string) (interface{}, error) {
	obj, err := s.SQL.getFirst(k, s.getData())
	if err != nil {
		return nil, wrapErr(err)
	}

	return *obj, nil
}

// GetFirstGlobal does the same as GetFirst but for all docs.
// Instead of returning an interface it returns a map with keys
// the index of the doc that a key was found and value the value of the key
func (s *Storage) GetFirstGlobal(k string) map[int]interface{} {
	found := make(map[int]interface{})

	c := s.GetAD()
	for j := range s.GetAllData() {
		s.SetAD(j)

		obj, err := s.SQL.getFirst(k, s.getData())
		if err != nil {
			continue
		}
		found[j] = *obj
	}

	s.SetAD(c)

	return found
}

// Get is alias of FindKeys. This function will be replaced
// by FindKeys in the future.
// For now we keep both for compatibility
func (s *Storage) Get(k string) ([]string, error) {
	issueWarning(deprecatedFeature, "Get()", "FindKeys()")

	obj, err := s.SQL.findKeys(k, s.getData())
	if err != nil {
		return nil, wrapErr(err)
	}

	return obj, nil
}

// FindKeys is a SQL wrapper that finds all the paths for a given
// e.g. ["key-1.test", "key-2.key-3.test"] will be returned
func (s *Storage) FindKeys(k string) ([]string, error) {
	obj, err := s.SQL.findKeys(k, s.getData())
	if err != nil {
		return nil, wrapErr(err)
	}

	return obj, nil
}

// FindKeysGlobal does the same as FindKeys but for all docs.
// Instead of returning a list of keys it returns a map with indexes
// from the docs and value an array of paths that was found
func (s *Storage) FindKeysGlobal(k string) map[int][]string {
	found := make(map[int][]string)

	c := s.GetAD()
	for j := range s.GetAllData() {
		s.SetAD(j)

		obj, err := s.SQL.findKeys(k, s.getData())
		if err != nil || len(obj) == 0 {
			continue
		}
		found[j] = obj
	}

	s.SetAD(c)

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
	obj, err := s.SQL.getPath(keys, s.getData())
	if err != nil {
		return nil, wrapErr(err)
	}

	return *obj, nil
}

// GetPathGlobal does the same as GetPath but globally for all
// docs
func (s *Storage) GetPathGlobal(k string) map[int]interface{} {
	found := make(map[int]interface{})
	keys := strings.Split(k, ".")

	c := s.GetAD()
	for j := range s.GetAllData() {
		s.SetAD(j)
		obj, err := s.SQL.getPath(keys, s.getData())
		if err != nil {
			continue
		}
		found[j] = *obj
	}

	s.SetAD(c)

	return found
}

// Delete is a SQL wrapper that deletes the last key from a given
// path. For example, Delete("key-1.key-2.key-3") would first
// validate that the path exists, then it would export the value of
// GetPath("key-1.key-2") and delete the object that matches key-3
func (s *Storage) Delete(k string) error {
	err := s.SQL.delPath(k, s.getData())
	if err != nil {
		return wrapErr(err)
	}

	return s.stateReload()
}

// DeleteGlobal is the same as Delete but will try to delete
// the path on all docs (if found)
func (s *Storage) DeleteGlobal(k string) error {
	c := s.GetAD()
	for j := range s.GetAllData() {
		s.SetAD(j)
		err := s.Delete(k)
		if err != nil {
			continue
		}
	}
	s.SetAD(c)

	return s.stateReload()
}

// MergeDBs is a SQL wrapper that merges a source yaml file
// with the DBy local yaml file.
func (s *Storage) MergeDBs(path string) error {
	err := s.SQL.mergeDBs(path, s.getData())
	if err != nil {
		return wrapErr(err)
	}

	return s.stateReload()
}
