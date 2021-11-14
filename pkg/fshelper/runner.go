package fshelper

import (
	"github.com/spf13/afero"
	"log"
	"os"
	"path"
)

const TempDirName = "swiztmpdir"

// RunInTempDir runs the specified function in a temporary directory that gets cleaned up
func (f FsHelp) RunInTempDir(useFullPath bool, exec TempDirExecFunc) error {
	log.Printf("[DEBUG] Creating temp dir %v", TempDirName)

	// Create dir if it doesn't exist. If the dir exists, it's probably because
	// of an earlier crash
	exists, err := afero.Exists(f.Fs, TempDirName)
	if err != nil {
		return err
	}

	if !exists {
		err = f.Fs.Mkdir(TempDirName, 0700)
		if err != nil {
			return err
		}
	}

	dir := TempDirName
	if useFullPath {
		wdDir, wdErr := os.Getwd()
		if wdErr != nil {
			return wdErr
		}
		dir = path.Join(wdDir, TempDirName)
	}

	// Execute
	err = exec(dir)
	if err != nil {
		return err
	}

	// Clean up
	log.Printf("[DEBUG] Cleaning up temp directory")
	err = os.RemoveAll(TempDirName)
	if err != nil {
		log.Printf("[WARN] Error removing temporary directory %v", TempDirName)
	}

	return nil
}
