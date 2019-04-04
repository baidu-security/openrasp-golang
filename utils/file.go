package utils

import (
	"os"
	"time"
)

func IsExist(filename string) bool {
	_, err := os.Stat(filename)
	return !os.IsNotExist(err)
}

func GetLastModified(filename string) (time.Time, bool) {
	if IsExist(filename) {
		file, err := os.Stat(filename)
		if err == nil {
			return file.ModTime(), true
		}
	}
	return time.Now(), false
}
