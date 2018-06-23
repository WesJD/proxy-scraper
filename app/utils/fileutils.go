package utils

import (
	"path/filepath"
	"os"
	"github.com/craftimize/sneaker/app/utils"
	"path"
)

var execPath string

func init() {
	ex, err := os.Executable()
	utils.CheckError(err)
	execPath = filepath.Dir(ex)
}

func Resource(resource string) (file string) {
	return path.Join(execPath, resource)
}

func Exists(resource string) bool {
	_, err := os.Stat(resource)
	return err == nil
}
