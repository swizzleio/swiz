package fshelper

// TempDirExecFunc is the function that runs in the temp path
type TempDirExecFunc func(path string) error

type FsHelper interface {
	CreatePath(filename string) error
	ReadJson(filename string, obj interface{}) error
	WriteJson(filename string, obj interface{}) error
	RunInTempDir(useFullPath bool, exec TempDirExecFunc) error
}
