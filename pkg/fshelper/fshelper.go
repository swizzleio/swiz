package fshelper

// TempDirExecFunc is the function that runs in the temp path
type TempDirExecFunc func(tmpPath string) error

type FsHelper interface {
	CreatePath(filename string) error
	ReadJson(filename string, obj interface{}) error
	WriteJson(filename string, obj interface{}) error
	WriteString(filename string, data string) error
	RunInTempDir(useFullPath bool, exec TempDirExecFunc) error
}
