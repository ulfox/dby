package db

import (
	"io/ioutil"
	"path/filepath"
	"sync"

	"github.com/pkg/errors"
	"gopkg.in/yaml.v2"
)

// Storage for managing RW of our
// state yaml file
type Storage struct {
	sync.Mutex
	SQL  SQL
	Data interface{}
	Path string
}

// NewStorageFactory for creating a new Storage struct.
func NewStorageFactory(path string) (*Storage, error) {
	state := &Storage{
		SQL:  SQL{},
		Path: path,
	}

	stateDir := filepath.Dir(path)
	err := makeDirs(stateDir, 0700)
	if err != nil {
		return nil, errors.Wrap(err, "NewStorageFactory")
	}

	stateExists, err := fileExists(path)
	if err != nil {
		return nil, errors.Wrap(err, "NewStorageFactory")
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
		return nil, errors.Wrap(err, "NewStorageFactory")
	}

	return state, nil
}

// Read for reading the local yaml file and importing it
// in memory
func (i *Storage) Read() error {
	f, err := ioutil.ReadFile(i.Path)
	if err != nil {
		return errors.Wrap(err, "Read")
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
		return errors.Wrap(err, "Write")
	}

	return ioutil.WriteFile(i.Path, data, 0600)
}

func (i *Storage) stateReload() error {
	err := i.Write()
	if err != nil {
		return errors.Wrap(err, "stateReload")
	}
	return errors.Wrap(i.Read(), "stateReload")
}
