package utils

import (
	"runtime"
	"path/filepath"
	"path"
)

func Directory() (dir string) {
	_, dirname, _, _ := runtime.Caller(2)
	dir = filepath.Dir(dirname)
	return
}

func Resource(resource string) (file string) {
	return path.Join(Directory(), resource)
}
