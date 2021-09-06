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

	"gopkg.in/yaml.v2"
)

// Storage is the main object exported by DBy. It consolidates together
// the Yaml Data and SQL
type Storage struct {
	sync.Mutex
	SQL  *SQL
	Data []interface{}
	Lib  map[string]int
	AD   int
	Path string
}

// NewStorageFactory for creating a new Storage
func NewStorageFactory(path string) (*Storage, error) {
	state := &Storage{
		SQL:  NewSQLFactory(),
		Path: path,
		Data: make([]interface{}, 0),
		Lib:  make(map[string]int),
	}

	stateDir := filepath.Dir(path)
	err := makeDirs(stateDir, 0700)
	if err != nil {
		return nil, wrapErr(err, getFn())
	}

	stateExists, err := fileExists(path)
	if err != nil {
		return nil, wrapErr(err, getFn())
	}

	if !stateExists {
		state.Data = append(state.Data, map[string]string{})
		state.AD = 0
		state.Write()
	}

	// Here we read the state. If we just created the file
	// then this is done to ensure everything was encoded and
	// written the right way.
	err = state.Read()
	if err != nil {
		return nil, wrapErr(err, getFn())
	}

	return state, nil
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
			wrapErr(fmt.Errorf(fieldNotString, strings.ToLower(f), kind), getFn())
		}

		sName, ok := name.(string)
		if !ok {
			wrapErr(fmt.Errorf(fieldNotString, strings.ToLower(l), name), getFn())
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
		return wrapErr(err, getFn())
	}
	s.Lib[strings.ToLower(n)] = i

	return nil
}

// Switch will change Active Document (AD) to the given index
func (s *Storage) Switch(i int) error {
	if i > len(s.Data)-1 {
		return wrapErr(fmt.Errorf(libOutOfIndex), getFn())
	}
	s.AD = i
	return nil
}

// AddDoc will add a new document to the stack and will switch
// Active Document index to that document
func (s *Storage) AddDoc() error {
	s.AD++
	s.Data = append(s.Data, make(map[interface{}]interface{}))
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
		return wrapErr(fmt.Errorf(docNotExists, strings.ToLower(n)), getFn())
	}
	s.AD = i
	return nil
}

// ImportDocs for importing documents
func (s *Storage) ImportDocs(path string) error {
	impf, err := ioutil.ReadFile(path)
	if err != nil {
		return wrapErr(err, getFn())
	}

	var dataArray []interface{}
	var counter int
	var data interface{}

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
		return wrapErr(err, getFn())
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

// Read for reading the local yaml file and importing it
// in memory
func (s *Storage) Read() error {
	f, err := ioutil.ReadFile(s.Path)
	if err != nil {
		return wrapErr(err, getFn())
	}

	s.Lock()
	defer s.Unlock()

	s.Data = nil
	s.Data = make([]interface{}, 0)

	var counter int
	var data interface{}
	dec := yaml.NewDecoder(bytes.NewReader(f))
	for {
		s.Data = append(s.Data, data)
		err := dec.Decode(&s.Data[counter])
		if err == nil {
			counter++
			data = nil
			continue
		}

		if err.Error() == "EOF" {
			break
		}
		return wrapErr(err, getFn())
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
		return wrapErr(err, getFn())
	}

	var buf bytes.Buffer
	enc := yaml.NewEncoder(&buf)

	for _, j := range s.Data {
		if j == nil {
			continue
		} else if v, ok := j.(map[interface{}]interface{}); ok && len(v) == 0 {
			continue
		}

		err := enc.Encode(j)
		if err != nil {
			return wrapErr(err, getFn())
		}
	}

	_, err = f.Write(buf.Bytes())
	if err != nil {
		return wrapErr(err, getFn())
	}
	err = f.Close()
	if err != nil {
		return wrapErr(err, getFn())
	}

	return wrapErr(os.Rename(f.Name(), s.Path), getFn())
}

func (s *Storage) stateReload() error {
	err := s.Write()
	if err != nil {
		return wrapErr(err, getFn())
	}

	return wrapErr(s.Read(), getFn())
}
