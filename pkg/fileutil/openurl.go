package fileutil

import (
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
	"path"
)

func OpenUrl(location string) ([]byte, error) {
	// Determine the protocol
	u, err := url.Parse(location)
	if err != nil {
		return nil, err
	}

	switch u.Scheme {
	case "file":
		return fileGet(path.Join(u.Host, u.Path))
	case "http":
		return httpGet(location)
	}

	return nil, fmt.Errorf("unsupported protocol: %s", u.Scheme)
}

func WriteUrl(location string, data []byte) error {
	// Determine the protocol
	u, err := url.Parse(location)
	if err != nil {
		return err
	}

	switch u.Scheme {
	case "file":
		return fileSave(path.Join(u.Host, u.Path), data)
	}

	return fmt.Errorf("unsupported protocol: %s", u.Scheme)
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
	// Open a file for reading
	file, err := os.Open(location)
	if err != nil {
		return nil, err
	}

	// Close the file when we are done
	defer file.Close()

	return io.ReadAll(file)
}

func fileSave(location string, data []byte) error {
	// Open a file for writing
	file, err := os.Create(location)
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
