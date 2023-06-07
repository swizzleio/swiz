package fileutil

import (
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
	"path"
	"path/filepath"
)

func GetPathFromUrl(location string, preserveFilename bool) (string, error) {
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
		homePath, err = expandUserPath(u.Host)
		if err != nil {
			return "", err
		}

		return path.Join(homePath, dir), nil
	}

	return "", fmt.Errorf("unsupported protocol: %s", u.Scheme)
}

func OpenUrl(location string) ([]byte, error) {
	return OpenUrlWithBaseDir("", location)
}

func OpenUrlWithBaseDir(baseDir string, location string) ([]byte, error) {
	fullLocation, err := UrlWithBaseDir(baseDir, location)
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
		return fileGet(fullLocation)
	case "https":
		return httpGet(fullLocation)
	}

	return nil, fmt.Errorf("unsupported protocol: %s", u.Scheme)
}

func UrlWithBaseDir(baseDir string, location string) (string, error) {
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

func WriteUrl(location string, data []byte) error {
	return WriteUrlWithBaseDir("", location, data)
}

func WriteUrlWithBaseDir(baseDir string, location string, data []byte) error {
	fullLocation, err := UrlWithBaseDir(baseDir, location)
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
		return fileSave(fullLocation, data)
	}

	return fmt.Errorf("unsupported protocol: %s", u.Scheme)
}

func GetScheme(location string) (string, error) {
	// Determine the protocol
	u, err := url.Parse(location)
	if err != nil {
		return "", err
	}

	return u.Scheme, nil
}

func expandUserPath(filePath string) (string, error) {
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

func httpGet(location string) ([]byte, error) {
	// Send HTTP GET request to the file URL
	response, err := http.Get(location)
	if err != nil {
		return nil, err
	}
	defer response.Body.Close()

	// Copy data from the response to byte[]
	return io.ReadAll(response.Body)
}

func fileGet(location string) ([]byte, error) {
	fullPath, err := GetPathFromUrl(location, true)
	if err != nil {
		return nil, err
	}

	// Open a file for reading
	file, err := os.Open(fullPath)
	if err != nil {
		return nil, err
	}

	// Close the file when we are done
	defer file.Close()

	return io.ReadAll(file)
}

func fileSave(location string, data []byte) error {
	fullPath, err := GetPathFromUrl(location, true)
	if err != nil {
		return err
	}

	// Open a file for writing
	file, err := os.Create(fullPath)
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
