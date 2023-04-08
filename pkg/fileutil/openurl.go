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

	return nil, fmt.Errorf("Unsupported protocol: %s", u.Scheme)
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
