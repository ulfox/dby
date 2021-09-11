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
	SQL  *SQL
	Data []interface{}
	Lib  map[string]int
	AD   int
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
		SQL:  NewSQLFactory(),
		Path: path,
		Data: make([]interface{}, 0),
		Lib:  make(map[string]int),
		mem:  inMem,
	}

	err := state.dbinit()
	if err != nil {
		return nil, wrapErr(err)
	}

	return state, nil
}

func (s *Storage) dbinit() error {
	if s.mem {
		s.Data = append(s.Data, emptyMap())
		s.AD = 0
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
		s.Data = append(s.Data, emptyMap())
		s.AD = 0
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
	for i := range s.Data {
		s.AD = i
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

		docName := fmt.Sprintf("%s/%s", strings.ToLower(sKind), strings.ToLower(sName))
		s.Lib[docName] = i
	}

	return nil
}

// SetName adds a name for a document and maps with it the given doc index
func (s *Storage) SetName(n string, i int) error {
	err := s.Switch(i)
	if err != nil {
		return wrapErr(err)
	}
	s.Lib[strings.ToLower(n)] = i

	return nil
}

// DeleteDoc will the document with the given index
func (s *Storage) DeleteDoc(i int) error {
	if i > len(s.Data)-1 && i != 0 {
		return wrapErr(libOutOfIndex)
	}

	s.Data = append(s.Data[:i], s.Data[i+1:]...)

	if s.AD == i {
		s.AD = 0
	}

	return nil
}

// Switch will change Active Document (AD) to the given index
func (s *Storage) Switch(i int) error {
	if i > len(s.Data)-1 && i != 0 {
		return wrapErr(libOutOfIndex)
	}
	s.AD = i
	return nil
}

// AddDoc will add a new document to the stack and will switch
// Active Document index to that document
func (s *Storage) AddDoc() error {
	if len(s.Data) == 0 {
		s.AD = 0
	} else {
		s.AD = len(s.Data)
	}
	s.Data = append(s.Data, emptyMap())
	return s.stateReload()
}

// ListDocs will return an array with all docs names
func (s *Storage) ListDocs() []string {
	var docs []string
	for i := range s.Lib {
		docs = append(docs, i)
	}
	return docs
}

// SwitchDoc for switching to a document using the documents name (if any)
func (s *Storage) SwitchDoc(n string) error {
	i, exists := s.Lib[strings.ToLower(n)]
	if !exists {
		return wrapErr(docNotExists, strings.ToLower(n))
	}
	s.AD = i
	return nil
}

// DeleteAll for removing all docs
func (s *Storage) DeleteAll(delete bool) *Storage {
	if delete {
		s.Data = make([]interface{}, 0)
	}
	return s
}

// ImportDocs for importing documents
func (s *Storage) ImportDocs(path string, o ...bool) error {
	impf, err := ioutil.ReadFile(path)
	if err != nil {
		return wrapErr(err)
	}

	var dataArray []interface{}
	var counter int
	var data interface{}

	if len(o) > 0 {
		issueWarning(deprecatedFeature, "ImportDocs(string, bool)", "Storage.DeleteAll(true).ImportDocs(path)")
		if o[0] {
			s.Data = nil
		}
	}

	data = nil
	dec := yaml.NewDecoder(bytes.NewReader(impf))
	for {
		dataArray = append(dataArray, data)
		err := dec.Decode(&dataArray[counter])
		if err == nil {
			counter++
			data = nil
			continue
		}

		if err.Error() == "EOF" {
			break
		}
		return wrapErr(err)
	}

	for _, j := range dataArray {
		if j == nil {
			continue
		}
		if len(j.(map[interface{}]interface{})) == 0 {
			continue
		}
		s.Data = append(s.Data, j)
	}
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

	s.Data = nil
	s.Data = make([]interface{}, 0)

	var data interface{}
	dec := yaml.NewDecoder(bytes.NewReader(f))
	for {
		err := dec.Decode(&data)
		if err == nil {
			s.Data = append(s.Data, data)
			data = nil
			continue
		}

		if err.Error() == "EOF" {
			break
		}
		return wrapErr(err)
	}

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

	for _, j := range s.Data {
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
