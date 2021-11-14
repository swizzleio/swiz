package fshelper

import (
	"encoding/json"
	"github.com/spf13/afero"
	"os"
	"path/filepath"
)

type FsHelp struct {
	Fs afero.Fs
}

func NewFsHelper() FsHelper {
	return FsHelp{Fs: afero.NewOsFs()}
}

func (f FsHelp) CreatePath(filename string) error {
	dir := filepath.Dir(filename)
	err := os.MkdirAll(dir, os.ModePerm)
	return err
}

// ReadJson reads a file and converts it to JSON
func (f FsHelp) ReadJson(filename string, obj interface{}) error {
	d, err := afero.ReadFile(f.Fs, filename)
	if err != nil {
		return err
	}

	err = json.Unmarshal(d, obj)

	return err
}

// WriteJson pretty writes an object to JSON
func (f FsHelp) WriteJson(filename string, obj interface{}) error {

	err := f.CreatePath(filename)
	if err != nil {
		return err
	}

	b, err := json.MarshalIndent(obj, "", "  ") // Yeah we marshall with spaces instead of tabs... come at me bro :)
	if err != nil {
		return err
	}

	err = afero.WriteFile(f.Fs, filename, b, 0600)

	return err
}

// WriteString writes a string to a file
func (f FsHelp) WriteString(filename string, data string) error {
	return afero.WriteFile(f.Fs, filename, []byte(data), 0600)
}
