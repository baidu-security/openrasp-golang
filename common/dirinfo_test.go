package common

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewWorkDirInfo(t *testing.T) {
	workDir := NewWorkDirInfo("/tmp", "workdir", 0755)
	subDir := NewWorkDirInfo("/tmp/workdir", "1", 0755)
	workDir.appendSubDir(subDir)
	assert.Equal(t, workDir.absPath(), "/tmp/workdir", "they should be equal")
	var err error
	err = workDir.init()
	assert.Nil(t, err)
	err = workDir.init()
	assert.Nil(t, err)
	err = workDir.clear()
	assert.Nil(t, err)
}
