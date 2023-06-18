package fileutil

import (
	"github.com/spf13/afero"
	"os"
)

type FileHelper interface {
	CreateDirIfNotExist(location string) error
}

type FileHelp struct {
	fh    FileUrlHelper
	appFs afero.Fs
}

func NewFileHelper() FileHelper {
	return &FileHelp{
		fh:    NewFileUrlHelper(),
		appFs: afero.NewOsFs(),
	}
}

func (f FileHelp) CreateDirIfNotExist(location string) error {
	dirLocation, err := f.fh.GetPathFromUrl(location, false)
	if err != nil {
		return err
	}

	if err = f.appFs.MkdirAll(dirLocation, 0755); err != nil && !os.IsExist(err) {
		return err
	}

	return nil
}
