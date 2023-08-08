package fileutil

import (
	"errors"
	"github.com/spf13/afero"
	"github.com/stretchr/testify/assert"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"testing"
)

func TestGetPathFromUrl(t *testing.T) {
	tests := []struct {
		name             string
		location         string
		preserveFilename bool
		expectedPath     string
		expectedErr      error
	}{
		{
			name:             "valid file path",
			location:         "file:///tmp/testfile",
			preserveFilename: true,
			expectedPath:     "/tmp/testfile",
			expectedErr:      nil,
		},
		{
			name:             "invalid scheme",
			location:         "https://example.com",
			preserveFilename: true,
			expectedPath:     "",
			expectedErr:      errors.New("unsupported protocol: https"),
		},
	}

	appFs := afero.NewMemMapFs()

	fileUrlHelper := &FileUrlHelp{appFs: appFs}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			result, err := fileUrlHelper.GetPathFromUrl(test.location, test.preserveFilename)
			assert.Equal(t, test.expectedPath, result)
			if test.expectedErr != nil {
				assert.EqualError(t, err, test.expectedErr.Error())
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func TestOpenUrlWithBaseDir(t *testing.T) {
	testServer := httptest.NewServer(http.HandlerFunc(func(res http.ResponseWriter, req *http.Request) {
		res.Write([]byte("fake response"))
	}))
	defer testServer.Close()

	tests := []struct {
		name        string
		baseDir     string
		location    string
		expectedRes []byte
		expectedErr error
	}{
		{
			name:        "valid http file path",
			baseDir:     "",
			location:    testServer.URL,
			expectedRes: []byte("fake response"),
			expectedErr: nil,
		},
		{
			name:        "invalid http file path",
			baseDir:     "",
			location:    "http://example.com/foobar",
			expectedRes: nil,
			expectedErr: errors.New("unsupported protocol: http"),
		},
		{
			name:        "invalid file path",
			baseDir:     "",
			location:    "file:///invalid/file/path",
			expectedRes: nil,
			expectedErr: errors.New("open /invalid/file/path: file does not exist"),
		},
	}

	appFs := afero.NewMemMapFs()
	fileUrlHelper := &FileUrlHelp{appFs: appFs}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			result, err := fileUrlHelper.OpenUrlWithBaseDir(test.baseDir, test.location)
			assert.Equal(t, test.expectedRes, result)
			if test.expectedErr != nil {
				assert.EqualError(t, err, test.expectedErr.Error())
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func TestUrlWithBaseDir(t *testing.T) {
	tests := []struct {
		name        string
		baseDir     string
		location    string
		expectedURL string
		expectedErr error
	}{
		{
			name:        "valid file URL",
			baseDir:     "/base",
			location:    "file://host/path",
			expectedURL: "file:///base/host/path",
			expectedErr: nil,
		},
		{
			name:        "valid https URL",
			baseDir:     "/base",
			location:    "https://example.com",
			expectedURL: "https://example.com",
			expectedErr: nil,
		},
	}

	appFs := afero.NewMemMapFs()
	fileUrlHelper := &FileUrlHelp{
		appFs: appFs,
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			result, err := fileUrlHelper.UrlWithBaseDir(test.baseDir, test.location)
			assert.Equal(t, test.expectedURL, result)
			if test.expectedErr != nil {
				assert.EqualError(t, err, test.expectedErr.Error())
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func TestWriteUrlWithBaseDir(t *testing.T) {
	tests := []struct {
		name        string
		baseDir     string
		location    string
		data        []byte
		expectedErr error
	}{
		{
			name:        "valid file URL",
			baseDir:     "",
			location:    "file:///tmp/testfile",
			data:        []byte("test data"),
			expectedErr: nil,
		},
	}

	appFs := afero.NewMemMapFs()
	fileUrlHelper := &FileUrlHelp{
		appFs: appFs,
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			err := fileUrlHelper.WriteUrlWithBaseDir(test.baseDir, test.location, test.data)
			if test.expectedErr != nil {
				assert.EqualError(t, err, test.expectedErr.Error())
			} else {
				assert.NoError(t, err)
				// Now let's read the file back to ensure it was written correctly
				fullPath, _ := fileUrlHelper.GetPathFromUrl(test.location, true)
				contents, _ := afero.ReadFile(appFs, filepath.Clean(fullPath))
				assert.Equal(t, test.data, contents)
			}
		})
	}
}

func TestGetScheme(t *testing.T) {
	tests := []struct {
		name        string
		location    string
		expectedRes string
		expectedErr error
	}{
		{
			name:        "http scheme",
			location:    "http://example.com",
			expectedRes: "http",
			expectedErr: nil,
		},
		{
			name:        "https scheme",
			location:    "https://example.com",
			expectedRes: "https",
			expectedErr: nil,
		},
		{
			name:        "file scheme",
			location:    "file:///path/to/file",
			expectedRes: "file",
			expectedErr: nil,
		},
		{
			name:        "invalid URL",
			location:    "://invalid",
			expectedRes: "",
			expectedErr: errors.New("parse \"://invalid\": missing protocol scheme"),
		},
	}

	fileUrlHelper := &FileUrlHelp{}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			result, err := fileUrlHelper.GetScheme(test.location)
			assert.Equal(t, test.expectedRes, result)
			if test.expectedErr != nil {
				assert.EqualError(t, err, test.expectedErr.Error())
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func TestExpandUserPath(t *testing.T) {
	tests := []struct {
		name        string
		filePath    string
		expectedRes string
		expectedErr error
	}{
		{
			name:        "valid user path",
			filePath:    "~/testfile",
			expectedRes: filepath.Join(os.Getenv("HOME"), "testfile"),
			expectedErr: nil,
		},
	}

	appFs := afero.NewMemMapFs()
	fileUrlHelper := &FileUrlHelp{
		appFs: appFs,
	}

	// This test is private method dependent, it might change according to implementation
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			// Use reflection to call the unexported method expandUserPath
			result, err := fileUrlHelper.expandUserPath(test.filePath)
			assert.Equal(t, test.expectedRes, result)
			if test.expectedErr != nil {
				assert.EqualError(t, err, test.expectedErr.Error())
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func TestHttpGet(t *testing.T) {
	testServer := httptest.NewServer(http.HandlerFunc(func(res http.ResponseWriter, req *http.Request) {
		res.Write([]byte("fake response"))
	}))
	defer testServer.Close()

	tests := []struct {
		name        string
		location    string
		expectedRes []byte
		expectedErr error
	}{
		{
			name:        "valid https location",
			location:    testServer.URL,
			expectedRes: []byte("fake response"),
			expectedErr: nil,
		},
		{
			name:        "invalid location",
			location:    "https://invalid-url",
			expectedRes: nil,
			expectedErr: errors.New("no such host"),
		},
	}

	fileUrlHelper := &FileUrlHelp{}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			// Use reflection to call the unexported method httpGet
			result, err := fileUrlHelper.httpGet(test.location)
			assert.Equal(t, test.expectedRes, result)
			if test.expectedErr != nil {
				assert.Contains(t, err.Error(), test.expectedErr.Error())
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func TestFileGet(t *testing.T) {
	appFs := afero.NewMemMapFs()
	fileUrlHelper := &FileUrlHelp{
		appFs: appFs,
	}

	// Create a test file
	afero.WriteFile(appFs, "/tmp/testfile", []byte("test data"), 0644)

	tests := []struct {
		name        string
		location    string
		expectedRes []byte
		expectedErr error
	}{
		{
			name:        "valid file location",
			location:    "file:///tmp/testfile",
			expectedRes: []byte("test data"),
			expectedErr: nil,
		},
		{
			name:        "invalid file location",
			location:    "file:///invalid/testfile",
			expectedRes: nil,
			expectedErr: errors.New("open /invalid/testfile: file does not exist"),
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			// Use reflection to call the unexported method fileGet
			result, err := fileUrlHelper.fileGet(test.location)
			assert.Equal(t, test.expectedRes, result)
			if test.expectedErr != nil {
				assert.Contains(t, err.Error(), test.expectedErr.Error())
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func TestFileSave(t *testing.T) {
	appFs := afero.NewMemMapFs()
	fileUrlHelper := &FileUrlHelp{
		appFs: appFs,
	}

	tests := []struct {
		name        string
		location    string
		data        []byte
		expectedErr error
	}{
		{
			name:        "save to valid file location",
			location:    "file:///tmp/testfile",
			data:        []byte("test data"),
			expectedErr: nil,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			// Use reflection to call the unexported method fileSave
			err := fileUrlHelper.fileSave(test.location, test.data)
			if test.expectedErr != nil {
				assert.EqualError(t, err, test.expectedErr.Error())
			} else {
				assert.NoError(t, err)
				// Now let's read the file back to ensure it was written correctly
				fullPath, _ := fileUrlHelper.GetPathFromUrl(test.location, true)
				contents, _ := afero.ReadFile(appFs, filepath.Clean(fullPath))
				assert.Equal(t, test.data, contents)
			}
		})
	}
}
