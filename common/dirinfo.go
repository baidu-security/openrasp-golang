package common

import (
	"os"
	"path/filepath"
	"syscall"
)

type WorkDirInfo struct {
	path string
	name string
	mode os.FileMode
	subs []*WorkDirInfo
}

func NewWorkDirInfo(path, name string, mode os.FileMode) *WorkDirInfo {
	subdirs := make([]*WorkDirInfo, 0)
	return &WorkDirInfo{
		path: path,
		name: name,
		mode: mode,
		subs: subdirs,
	}
}

func (wdi *WorkDirInfo) appendSubDir(subdir *WorkDirInfo) {
	wdi.subs = append(wdi.subs, subdir)
}

func (wdi *WorkDirInfo) absPath() string {
	return filepath.Join(wdi.path, wdi.name)
}

func (wdi *WorkDirInfo) init() error {
	currentDir := wdi.absPath()
	if _, err := os.Stat(currentDir); err == nil {
		err := syscall.Access(currentDir, syscall.O_RDWR)
		if err != nil {
			return err
		}
	} else if os.IsNotExist(err) {
		err := os.Mkdir(currentDir, wdi.mode)
		if err != nil {
			return err
		}
	} else {
		return err
	}
	for _, v := range wdi.subs {
		err := v.init()
		if err != nil {
			return err
		}
	}
	return nil
}

func (wdi *WorkDirInfo) clear() error {
	currentDir := wdi.absPath()
	return os.RemoveAll(currentDir)
}
