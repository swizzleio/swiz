package fshelper

import (
	"log"
	"os"
	"path"
)

const TempDirName = "swiztmpdir"

// RunInTempDir runs the specified function in a temporary directory that gets cleaned up
func (f FsHelp) RunInTempDir(useFullPath bool, exec TempDirExecFunc) error {
	log.Printf("[DEBUG] Creating temp dir %v", TempDirName)

	// Create dir
	err := f.Fs.Mkdir(TempDirName, 0700)
	if err != nil {
		return err
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
