package db

import (
	"io/ioutil"
	"os"
	"path"
	"path/filepath"
	"sync"

	"gopkg.in/yaml.v2"
)

// Storage is the main object exported by DBy. It consolidates together
// the Yaml Data and SQL
type Storage struct {
	sync.Mutex
	SQL  *SQL
	Data interface{}
	Path string
}

// NewStorageFactory for creating a new Storage
func NewStorageFactory(path string) (*Storage, error) {
	state := &Storage{
		SQL:  NewSQLFactory(),
		Path: path,
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
		state.Data = map[string]string{}
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

// Read for reading the local yaml file and importing it
// in memory
func (i *Storage) Read() error {
	f, err := ioutil.ReadFile(i.Path)
	if err != nil {
		return wrapErr(err, getFn())
	}

	i.Lock()
	defer i.Unlock()

	return yaml.Unmarshal(f, &i.Data)
}

// Write for writing memory content to the local yaml file
func (i *Storage) Write() error {
	i.Lock()
	defer i.Unlock()
	data, err := yaml.Marshal(&i.Data)
	if err != nil {
		return wrapErr(err, getFn())
	}

	wrkDir := path.Dir(i.Path)
	f, err := ioutil.TempFile(wrkDir, ".tx.*")
	if err != nil {
		return wrapErr(err, getFn())
	}

	_, err = f.Write(data)
	if err != nil {
		return wrapErr(err, getFn())
	}
	err = f.Close()
	if err != nil {
		return wrapErr(err, getFn())
	}

	return wrapErr(os.Rename(f.Name(), i.Path), getFn())
}

func (i *Storage) stateReload() error {
	err := i.Write()
	if err != nil {
		return wrapErr(err, getFn())
	}
	return wrapErr(i.Read(), getFn())
}
