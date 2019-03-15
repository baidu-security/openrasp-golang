package common

import (
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/baidu/openrasp/test"
)

func TestWorkDir(t *testing.T) {
	wd, _ := getExecutableDir()
	binDir, err := filepath.Abs(filepath.Dir(os.Args[0]))
	if err != nil {
		t.Errorf("Fail to get binDir.")
	}
	if wd != binDir {
		t.Errorf("Work directory may be wrong.")
	}
}

func TestWorkDirInitAccessable(t *testing.T) {
	var workSpace *WorkSpace = NewWorkSpace()
	testDir := test.TempMkdir(t)
	defer os.RemoveAll(testDir)
	subDir := filepath.Join(testDir, "accessable")
	os.Mkdir(subDir, 0755)
	err := workSpace.workDirs.init(subDir, workSpace.register)
	if err != nil {
		t.Errorf("Fail to init.")
	}
}

func TestWorkDirInitAlreadyInitialized(t *testing.T) {
	var workSpace *WorkSpace = NewWorkSpace()
	testDir := test.TempMkdir(t)
	defer os.RemoveAll(testDir)
	subDir := filepath.Join(testDir, "accessable")
	os.Mkdir(subDir, 0755)
	err := workSpace.workDirs.init(subDir, workSpace.register)
	if err != nil {
		t.Errorf("Fail to init.")
	}
	err = workSpace.workDirs.init(subDir, workSpace.register)
	if err != nil {
		t.Errorf("Fail to init.")
	}
}

func TestWorkDirInitNotAccessable(t *testing.T) {
	var workSpace *WorkSpace = NewWorkSpace()
	testDir := test.TempMkdir(t)
	defer os.RemoveAll(testDir)
	subDir := filepath.Join(testDir, "notaccessable")
	os.Mkdir(subDir, 0444)
	err := workSpace.workDirs.init(subDir, workSpace.register)
	if err == nil {
		t.Errorf("Init should return error.")
	}
}

func TestWorkSpaceInitClear(t *testing.T) {
	var workSpace *WorkSpace = NewWorkSpace()
	workSpace.Init()
	if !workSpace.Active() {
		t.Errorf("Actice() should return true.")
	}
	_, errInit := workSpace.GetDir(Root)
	if errInit != nil {
		t.Errorf("Fail to get root dir path.")
	}
	workSpace.clear()
	if workSpace.Active() {
		t.Errorf("Actice() should return false after work space clear.")
	}
	_, errClear := workSpace.GetDir(Root)
	if errClear == nil {
		t.Errorf("GetDir should return error.")
	}
}

func TestWorkSpaceWithoutInit(t *testing.T) {
	var workSpace *WorkSpace = NewWorkSpace()
	if workSpace.Active() {
		t.Errorf("Actice() should return false after work space clear.")
	}
	_, err := workSpace.GetDir(Root)
	if err == nil {
		t.Errorf("GetDir should return error.")
	}
}

type testListener struct {
	latest string
}

func (tl *testListener) OnUpdate(absPath string) {
	tl.latest = absPath
}

func TestWorkSpaceInitWatch(t *testing.T) {
	var workSpace *WorkSpace = NewWorkSpace()
	workSpace.Init()
	if !workSpace.Active() {
		t.Errorf("Actice() should return true.")
	}
	_, errInit := workSpace.GetDir(Root)
	if errInit != nil {
		t.Errorf("Fail to get root dir path.")
	}
	if !workSpace.StartWatch(Conf, Plugins) {
		t.Errorf("Fail to watch conf dir.")
	}
	if workSpace.StartWatch(Plugins) {
		t.Errorf("Cannot watch plugins dir twice.")
	}
	tl := &testListener{
		latest: "",
	}
	confDir, errConf := workSpace.GetDir(Conf)
	if errConf != nil {
		t.Errorf("Fail to get conf dir path.")
	}
	workSpace.RegisterListener(Conf, tl)
	testConfFile := filepath.Join(confDir, "openrasp.yml")
	os.Mkdir(testConfFile, 0755)
	defer os.RemoveAll(testConfFile)
	<-time.After(50 * time.Millisecond)
	if tl.latest != testConfFile {
		t.Errorf("Fail to watch conf dir.")
	}
	workSpace.clear()
	if workSpace.Active() {
		t.Errorf("Actice() should return false after work space clear.")
	}
	_, errClear := workSpace.GetDir(Root)
	if errClear == nil {
		t.Errorf("GetDir should return error.")
	}
}
