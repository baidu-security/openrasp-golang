package test

import (
	"io/ioutil"
	"testing"
)

// TempMkdir makes a temporary directory
func TempMkdir(t *testing.T) string {
	dir, err := ioutil.TempDir("", "openrasp-notify-test")
	if err != nil {
		t.Fatalf("failed to create test directory: %s", err)
	}
	return dir
}
