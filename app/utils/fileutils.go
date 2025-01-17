package utils

import (
	"path/filepath"
	"os"
	"path"
	"strings"
)

var execPath string

func init() {
	ex, err := os.Executable()
	if strings.HasPrefix(ex, "/tmp") {
		ex = ""
	}
	CheckError(err)
	execPath = filepath.Dir(ex)
}

func Resource(resource string) (file string) {
	return path.Join(execPath, resource)
}

func Exists(resource string) bool {
	_, err := os.Stat(resource)
	return err == nil
}
