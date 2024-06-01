package util

import (
	"github.com/spf13/afero"
	"io"
	"os"
	"path/filepath"
)

// CopyFile copies a single file from src to dst using the provided Fs (file system).
func CopyFile(fs afero.Fs, src string, dst string) error {
	srcFile, err := os.Open(src)
	if err != nil {
		return err
	}
	defer srcFile.Close()

	// Ensure the destination directory exists
	if err := fs.MkdirAll(filepath.Dir(dst), 0755); err != nil {
		return err
	}

	dstFile, err := fs.Create(dst)
	if err != nil {
		return err
	}
	defer dstFile.Close()

	_, err = io.Copy(dstFile, srcFile)
	return err
}

// CopyDir recursively copies a directory from the OS filesystem to an afero Fs.
func CopyDir(fs afero.Fs, srcDir string, dstDir string) error {
	// Walk the source directory tree
	return filepath.Walk(srcDir, func(srcPath string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		// Create a relative path from the srcDir for the destination directory
		relPath, err := filepath.Rel(srcDir, srcPath)
		if err != nil {
			return err
		}
		dstPath := filepath.Join(dstDir, relPath)

		// Check if it's a directory, if so, create it
		if info.IsDir() {
			return fs.MkdirAll(dstPath, info.Mode())
		}

		// It's a file, copy it
		return CopyFile(fs, srcPath, dstPath)
	})
}
