package fileutil

import (
	"net/url"
	"os"
	"path"
)

func CreateDirIfNotExist(location string) error {
	// Determine the protocol
	u, err := url.Parse(location)
	if err != nil {
		return err
	}

	if err = os.MkdirAll(path.Join(u.Host, u.Path), 0755); err != nil && !os.IsExist(err) {
		return err
	}

	return nil
}
