package utils

import (
	"path"
	"path/filepath"
	"runtime"
)

func GetFileRootDir() string {
	_, b, _, _ := runtime.Caller(0)
	d := path.Join(path.Dir(b))
	return filepath.Dir(d)
}
