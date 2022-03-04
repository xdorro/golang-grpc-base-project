package utils

import (
	"os"
	"path/filepath"
)

// Exists reports whether the named file or directory exists.
func Exists(name string) bool {
	if _, err := os.Stat(name); err != nil {
		if os.IsNotExist(err) {
			return false
		}
	}

	return true
}

// MakeDir creates a directory if it does not exist.
func MakeDir(dir string) error {
	dirConvert := filepath.Dir(dir)
	if !Exists(dirConvert) {
		err := os.Mkdir(dirConvert, 0775)
		if err != nil {
			return err
		}
	}

	return nil
}
