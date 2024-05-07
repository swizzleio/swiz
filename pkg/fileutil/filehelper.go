package fileutil

import (
	"os"

	"github.com/spf13/afero"
)

//go:generate mockery --name FileHelper --filename filehelper_mock.go --output ../../mocks/pkg/fileutil --outpkg mockfileutil
type FileHelper interface {
	CreateDirIfNotExist(location string) error
}

type FileHelp struct {
	fh    FileUrlHelper
	appFs afero.Fs
}

func NewFileHelper(appFs afero.Fs) FileHelper {
	if appFs == nil {
		appFs = afero.NewOsFs()
	}
	return &FileHelp{
		fh:    NewFileUrlHelper(appFs),
		appFs: appFs,
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
