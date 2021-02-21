package file

import (
	"io"
	"os"
	"sync"
)

// File is used to represent a file in disk, marshall and unmarshal should be set by the caller with the apropriate type
type File struct {
	marshal   func(interface{}) (io.Reader, error)
	unmarshal func(io.Reader, interface{}) error
	mutex     sync.Mutex
	path      string
}

// NewFile creates the File struct
func NewFile(p string,
	m func(interface{}) (io.Reader, error),
	u func(io.Reader, interface{}) error) *File {
	return &File{
		path:      p,
		marshal:   m,
		unmarshal: u,
	}
}

// Create file in disk
func (f *File) Create() error {
	pf, err := os.Create(f.path)
	if err != nil {
		return err
	}
	pf.Close()
	return nil
}

// Save file in disk
func (f *File) Save(d interface{}) error {
	f.mutex.Lock()
	defer f.mutex.Unlock()
	fp, err := os.OpenFile(f.path, os.O_APPEND|os.O_WRONLY, os.ModeAppend)
	if err != nil {
		return err
	}
	defer fp.Close()
	b, err := f.marshal(d)
	if err != nil {
		return err
	}
	_, err = io.Copy(fp, b)
	return err
}

// Load the file from disk
func (f *File) Load(d interface{}) error {
	f.mutex.Lock()
	defer f.mutex.Unlock()
	fp, err := os.Open(f.path)
	if err != nil {
		return err
	}
	return f.unmarshal(fp, d)
}
