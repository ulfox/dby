package db

import (
	"io/ioutil"
	"path/filepath"
	"sync"

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

func NewStorageFactory(path string) (*Storage, error) {
	state := &Storage{
		SQL:  SQL{},
		Path: path,
	}

	stateDir := filepath.Dir(path)
	err := makeDirs(stateDir, 0700)
	if err != nil {
		return nil, err
	}

	stateExists, err := fileExists(path)
	if err != nil {
		return nil, err
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
		return nil, err
	}

	return state, nil
}

func (i *Storage) Read() error {
	f, err := ioutil.ReadFile(i.Path)
	if err != nil {
		return err
	}

	i.Lock()
	defer i.Unlock()

	return yaml.Unmarshal(f, &i.Data)
}

func (i *Storage) Write() error {
	i.Lock()
	defer i.Unlock()
	data, err := yaml.Marshal(&i.Data)
	if err != nil {
		return err
	}

	return ioutil.WriteFile(i.Path, data, 0600)
}

func (i *Storage) stateReload() error {
	err := i.Write()
	if err != nil {
		return err
	}
	return i.Read()
}
