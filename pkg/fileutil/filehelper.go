package fileutil

import (
	"os"
)

func CreateDirIfNotExist(location string) error {
	dirLocation, err := GetPathFromUrl(location, false)
	if err != nil {
		return err
	}

	if err = os.MkdirAll(dirLocation, 0755); err != nil && !os.IsExist(err) {
		return err
	}

	return nil
}
