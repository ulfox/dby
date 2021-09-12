package db

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"os"
	"path"
	"path/filepath"
	"strings"
	"sync"

	e "github.com/ulfox/dby/errors"
	"gopkg.in/yaml.v2"
)

type erf = func(e interface{}, p ...interface{}) error

var wrapErr erf = e.WrapErr

// Storage is the main object exported by DBy. It consolidates together
// the Yaml Data and SQL
type Storage struct {
	sync.Mutex
	*state
	SQL  *SQL
	Path string
	mem  bool
}

// NewStorageFactory for creating a new Storage
func NewStorageFactory(p ...interface{}) (*Storage, error) {
	var path string = "local/dby.yaml"
	var inMem bool = true

	if len(p) > 0 {
		switch i := p[0].(type) {
		case string:
			path = i
			inMem = false
		case bool:
			inMem = i
		}
	}

	state := &Storage{
		SQL:   NewSQLFactory(),
		state: newStateFactory(),
		Path:  path,
		mem:   inMem,
	}

	err := state.dbinit()
	if err != nil {
		return nil, wrapErr(err)
	}

	return state, nil
}

func (s *Storage) Close() error {
	if !s.mem {
		err := s.Write()
		if err != nil {
			return wrapErr(err)
		}
	}

	s.Clear()
	s.SQL.Clear()
	return nil
}

func (s *Storage) dbinit() error {
	if s.mem {
		s.PushData(emptyMap())
		s.SetAD(0)
		return nil
	}

	stateDir := filepath.Dir(s.Path)
	err := makeDirs(stateDir, 0700)
	if err != nil {
		return wrapErr(err)
	}

	stateExists, err := fileExists(s.Path)
	if err != nil {
		return wrapErr(err)
	}

	if !stateExists {
		s.PushData(emptyMap())
		s.SetAD(0)
		err = s.Write()
		if err != nil {
			return wrapErr(err)
		}
	}

	err = s.Read()
	if err != nil {
		return wrapErr(err)
	}

	return nil
}

// SetNames can set names automatically to the documents
// that have the queried paths.
// input(f) is the first path that will be quieried
// input(l) is the last path
//
// If a document has both paths, a name will be generated
// and will be mapped with the document's index
func (s *Storage) SetNames(f, l string) error {
	for i := range s.GetAllData() {
		s.SetAD(i)
		kind, err := s.GetPath(strings.ToLower(f))
		if err != nil {
			continue
		}
		name, err := s.GetPath(strings.ToLower(l))
		if err != nil {
			continue
		}

		sKind, ok := kind.(string)
		if !ok {
			wrapErr(fieldNotString, strings.ToLower(f), kind)
		}

		sName, ok := name.(string)
		if !ok {
			wrapErr(fieldNotString, strings.ToLower(l), name)
		}
		err = s.addDoc(
			fmt.Sprintf(
				"%s/%s",
				strings.ToLower(sKind),
				strings.ToLower(sName),
			),
			i,
		)
		if err != nil {
			return wrapErr(err)
		}
	}

	return nil
}

// SetName adds a name for a document and maps with it the given doc index
func (s *Storage) SetName(n string, i int) error {
	err := s.Switch(i)
	if err != nil {
		return wrapErr(err)
	}
	err = s.addDoc(n, i)
	if err != nil {
		return wrapErr(err)
	}

	return nil
}

// DeleteDoc will the document with the given index
func (s *Storage) DeleteDoc(i int) error {
	err := s.DeleteData(i)
	if err != nil {
		return wrapErr(err)
	}

	return nil
}

// Switch will change Active Document (AD) to the given index
func (s *Storage) Switch(i int) error {
	err := s.SetAD(i)
	if err != nil {
		return wrapErr(err)
	}
	return nil
}

// AddDoc will add a new document to the stack and will switch
// Active Document index to that document
func (s *Storage) AddDoc() error {
	s.PushData(emptyMap())
	s.SetAD(len(s.GetAllData()) - 1)
	return s.stateReload()
}

// ListDocs will return an array with all docs names
func (s *Storage) ListDocs() []string {
	var docs []string
	for i := range s.Lib() {
		docs = append(docs, i)
	}
	return docs
}

// SwitchDoc for switching to a document using the documents name (if any)
func (s *Storage) SwitchDoc(n string) error {
	i, exists := s.LibIndex(n)
	if !exists {
		return wrapErr(docNotExists, strings.ToLower(n))
	}
	s.SetAD(i)
	return nil
}

// DeleteAll for removing all docs
func (s *Storage) DeleteAll(delete bool) *Storage {
	if delete {
		s.DeleteAllData()
		s.ClearLib()
	}
	return s
}

// ImportDocs for importing documents
func (s *Storage) ImportDocs(path string, o ...bool) error {
	impf, err := ioutil.ReadFile(path)
	if err != nil {
		return wrapErr(err)
	}

	var counter int
	var data interface{}
	s.UnsetBufferArray()

	dec := yaml.NewDecoder(bytes.NewReader(impf))
	for {
		err = dec.Decode(&data)
		if err == nil {
			s.PushBuffer(data)
			data = nil

			counter++
			continue
		}

		if err.Error() == "EOF" {
			break
		}
		s.UnsetBufferArray()
		return wrapErr(err)
	}

	if len(o) > 0 {
		issueWarning(deprecatedFeature, "ImportDocs(string, bool)", "Storage.DeleteAll(true).ImportDocs(path)")
		if o[0] {
			s.UnsetDataArray()
			s.ClearLib()
		}
	}

	for _, j := range s.GetAllBuffer() {
		if j == nil {
			continue
		}
		if len((*j).(map[interface{}]interface{})) == 0 {
			continue
		}
		s.PushData(*j)
	}
	s.UnsetBufferArray()
	return s.stateReload()
}

// InMem for configuring db to write only in memory
func (s *Storage) InMem(m bool) *Storage {
	s.mem = m
	return s
}

// Read for reading the local yaml file and importing it
// in memory
func (s *Storage) Read() error {
	f, err := ioutil.ReadFile(s.Path)
	if err != nil {
		return wrapErr(err)
	}

	s.Lock()
	defer s.Unlock()

	s.UnsetBufferArray()

	var data interface{}
	dec := yaml.NewDecoder(bytes.NewReader(f))
	for {
		err := dec.Decode(&data)
		if err == nil {
			s.PushBuffer(data)
			data = nil
			continue
		}

		if err.Error() == "EOF" {
			break
		}
		s.UnsetBufferArray()
		return wrapErr(err)
	}

	s.UnsetDataArray()

	for _, j := range s.GetAllBuffer() {
		if j == nil {
			continue
		}
		s.PushData(*j)
	}
	s.UnsetBufferArray()
	return nil
}

// Write for writing memory content to the local yaml file
func (s *Storage) Write() error {
	s.Lock()
	defer s.Unlock()

	wrkDir := path.Dir(s.Path)
	f, err := ioutil.TempFile(wrkDir, ".tx.*")
	if err != nil {
		return wrapErr(err)
	}

	var buf bytes.Buffer
	enc := yaml.NewEncoder(&buf)

	for _, j := range s.GetAllData() {
		if j == nil {
			continue
		}

		err := enc.Encode(j)
		if err != nil {
			return wrapErr(err)
		}
	}

	_, err = f.Write(buf.Bytes())
	if err != nil {
		return wrapErr(err)
	}
	err = f.Close()
	if err != nil {
		return wrapErr(err)
	}

	return wrapErr(os.Rename(f.Name(), s.Path))
}

func (s *Storage) stateReload() error {
	if s.mem {
		return nil
	}

	err := s.Write()
	if err != nil {
		return wrapErr(err)
	}

	return wrapErr(s.Read())
}
