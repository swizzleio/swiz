package fileutil

import (
	"fmt"
	"github.com/spf13/afero"
	"io"
	"net/http"
	"net/url"
	"os"
	"path"
	"path/filepath"
)

//go:generate mockery --name FileUrlHelper --filename openurl_mock.go
type FileUrlHelper interface {
	GetPathFromUrl(location string, preserveFilename bool) (string, error)
	GetScheme(location string) (string, error)
	UrlWithBaseDir(baseDir string, location string) (string, error)
	OpenUrl(location string) ([]byte, error)
	OpenUrlWithBaseDir(baseDir string, location string) ([]byte, error)
	WriteUrl(location string, data []byte) error
	WriteUrlWithBaseDir(baseDir string, location string, data []byte) error
}

type FileUrlHelp struct {
	appFs afero.Fs
}

func NewFileUrlHelper() FileUrlHelper {
	return &FileUrlHelp{
		appFs: afero.NewOsFs(),
	}
}

func (f FileUrlHelp) GetPathFromUrl(location string, preserveFilename bool) (string, error) {
	// Determine the protocol
	u, err := url.Parse(location)
	if err != nil {
		return "", err
	}

	switch u.Scheme {
	case "file":
		dir := u.Path
		if !preserveFilename {
			dir, _ = filepath.Split(u.Path)
		}

		var homePath string
		homePath, err = f.expandUserPath(u.Host)
		if err != nil {
			return "", err
		}

		return path.Join(homePath, dir), nil
	}

	return "", fmt.Errorf("unsupported protocol: %s", u.Scheme)
}

func (f FileUrlHelp) OpenUrl(location string) ([]byte, error) {
	return f.OpenUrlWithBaseDir("", location)
}

func (f FileUrlHelp) OpenUrlWithBaseDir(baseDir string, location string) ([]byte, error) {
	fullLocation, err := f.UrlWithBaseDir(baseDir, location)
	if err != nil {
		return nil, err
	}

	// Determine the protocol
	u, err := url.Parse(fullLocation)
	if err != nil {
		return nil, err
	}

	switch u.Scheme {
	case "file":
		return f.fileGet(fullLocation)
	case "https":
		return f.httpGet(fullLocation)
	}

	return nil, fmt.Errorf("unsupported protocol: %s", u.Scheme)
}

func (f FileUrlHelp) UrlWithBaseDir(baseDir string, location string) (string, error) {
	// Determine the protocol
	u, err := url.Parse(location)
	if err != nil {
		return "", err
	}

	switch u.Scheme {
	case "file":
		return fmt.Sprintf("%v://%v", u.Scheme, path.Join(path.Join(baseDir, u.Host), u.Path)), nil
	}

	return location, nil
}

func (f FileUrlHelp) WriteUrl(location string, data []byte) error {
	return f.WriteUrlWithBaseDir("", location, data)
}

func (f FileUrlHelp) WriteUrlWithBaseDir(baseDir string, location string, data []byte) error {
	fullLocation, err := f.UrlWithBaseDir(baseDir, location)
	if err != nil {
		return err
	}

	// Determine the protocol
	u, err := url.Parse(fullLocation)
	if err != nil {
		return err
	}

	switch u.Scheme {
	case "file":
		return f.fileSave(fullLocation, data)
	}

	return fmt.Errorf("unsupported protocol: %s", u.Scheme)
}

func (f FileUrlHelp) GetScheme(location string) (string, error) {
	// Determine the protocol
	u, err := url.Parse(location)
	if err != nil {
		return "", err
	}

	return u.Scheme, nil
}

func (f FileUrlHelp) expandUserPath(filePath string) (string, error) {
	if filePath == "" || filePath[0] != '~' {
		return filePath, nil
	}

	var homeDir string
	if homeDir = os.Getenv("HOME"); homeDir == "" {
		homeDir = os.Getenv("USERPROFILE") // For Windows
		if homeDir == "" {
			return "", fmt.Errorf("user home directory not found")
		}
	}

	return filepath.Join(homeDir, filePath[1:]), nil
}

func (f FileUrlHelp) httpGet(location string) ([]byte, error) {
	// Send HTTP GET request to the file URL
	response, err := http.Get(location)
	if err != nil {
		return nil, err
	}
	defer response.Body.Close()

	// Copy data from the response to byte[]
	return io.ReadAll(response.Body)
}

func (f FileUrlHelp) fileGet(location string) ([]byte, error) {
	fullPath, err := f.GetPathFromUrl(location, true)
	if err != nil {
		return nil, err
	}

	// Open a file for reading
	file, err := f.appFs.Open(fullPath)
	if err != nil {
		return nil, err
	}

	// Close the file when we are done
	defer file.Close()

	return io.ReadAll(file)
}

func (f FileUrlHelp) fileSave(location string, data []byte) error {
	fullPath, err := f.GetPathFromUrl(location, true)
	if err != nil {
		return err
	}

	// Open a file for writing
	file, err := f.appFs.Create(fullPath)
	if err != nil {
		return err
	}

	// Close the file when we are done
	defer file.Close()

	// Write data to the file
	_, err = file.Write(data)
	if err != nil {
		return err
	}

	return nil
}
