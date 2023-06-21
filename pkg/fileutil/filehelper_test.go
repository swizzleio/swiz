package fileutil

import (
	"github.com/spf13/afero"
	"github.com/stretchr/testify/assert"
	"github.com/swizzleio/swiz/mocks/pkg/fileutil"
	"os"
	"testing"
)

func TestCreateDirIfNotExist(t *testing.T) {
	// Set up test cases
	testCases := []struct {
		name        string
		location    string
		dirLocation string
		shouldError bool
		pathFail    bool
	}{
		{
			name:        "Successful case",
			location:    "path/to/dir",
			dirLocation: "path/to/dir",
			shouldError: false,
		},
		{
			name:        "Error: unable to get path from URL",
			location:    "invalidUrl",
			dirLocation: "",
			shouldError: true,
			pathFail:    true,
		},
		{
			name:        "directory exists",
			location:    "alreadyexists",
			dirLocation: "alreadyexists",
			shouldError: false,
		},
	}

	// Mock the filesystem
	fs := afero.NewMemMapFs()
	assert.NoError(t, fs.MkdirAll("alreadyexists", 0755))

	for _, tt := range testCases {
		t.Run(tt.name, func(t *testing.T) {
			// Mock FileUrlHelper
			mockFuh := mockfileutil.FileUrlHelper{}
			if tt.pathFail {
				mockFuh.On("GetPathFromUrl", tt.location, false).Return(tt.dirLocation, os.ErrInvalid)
			} else {
				mockFuh.On("GetPathFromUrl", tt.location, false).Return(tt.dirLocation, nil)
			}

			f := &FileHelp{
				fh:    &mockFuh,
				appFs: fs,
			}

			// Test CreateDirIfNotExist
			err := f.CreateDirIfNotExist(tt.location)
			if tt.shouldError {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				exists, _ := afero.DirExists(fs, tt.dirLocation)
				assert.True(t, exists)
			}
		})
	}
}
