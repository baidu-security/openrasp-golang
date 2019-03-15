package notify

import (
	"os"
	"path/filepath"
	"sync/atomic"
	"testing"
	"time"

	"github.com/baidu/openrasp/test"
)

func TestEventMatchWatch(t *testing.T) {
	var res bool = false
	res = eventMatchWatch("/a/b", "/a/b/c")
	if !res {
		t.Fatalf("eventMatchWatch not work.")
	}

	res = eventMatchWatch("/a/b", "/a/b")
	if !res {
		t.Fatalf("eventMatchWatch not work.")
	}

	res = eventMatchWatch("/a/b", "/a/notb")
	if res {
		t.Fatalf("eventMatchWatch not work.")
	}

	res = eventMatchWatch("/a/b", "/a")
	if res {
		t.Fatalf("eventMatchWatch not work.")
	}

	res = eventMatchWatch("/a/b", "/else")
	if res {
		t.Fatalf("eventMatchWatch not work.")
	}

	res = eventMatchWatch("./a/b", "./a/b/c")
	if !res {
		t.Fatalf("eventMatchWatch not work.")
	}

	res = eventMatchWatch("./a/b", "./a/b")
	if !res {
		t.Fatalf("eventMatchWatch not work.")
	}
}

func TestNotifyCloseWithoutWatch(t *testing.T) {
	w, err := NewWatcher()
	if err != nil {
		t.Fatalf("Failed to create watcher")
	}

	<-time.After(50 * time.Millisecond)
	w.Close()

	// Wait for the close to complete.
	<-time.After(50 * time.Millisecond)
	isWatcherReallyClosed(t, w)
}

func TestNotifyCloseWithWatch(t *testing.T) {
	testDir := test.TempMkdir(t)
	defer os.RemoveAll(testDir)

	w, err := NewWatcher()
	if err != nil {
		t.Fatalf("Failed to create watcher")
	}
	handler := EventHandler{
		EventFilter:   succcessFilterFunc,
		EventDispatch: func(name string, op EventOp) {},
	}
	err = w.Add(testDir, handler)
	if err != nil {
		t.Fatalf("Failed to add.")
	}

	// Wait until readEvents has reached unix.Read, and Close.
	<-time.After(50 * time.Millisecond)
	w.Close()

	// Wait for the close to complete.
	<-time.After(50 * time.Millisecond)
	isWatcherReallyClosed(t, w)
}

func TestNotifyRepeatedAdd(t *testing.T) {
	testDir := test.TempMkdir(t)
	defer os.RemoveAll(testDir)

	w, err := NewWatcher()
	if err != nil {
		t.Fatalf("Failed to create watcher")
	}
	handler := EventHandler{
		EventFilter:   succcessFilterFunc,
		EventDispatch: func(name string, op EventOp) {},
	}
	err = w.Add(testDir, handler)
	if err != nil {
		t.Fatalf("Failed to add.")
	}

	err = w.Add(testDir, handler)
	if err == nil {
		t.Fatalf("Repeated add should fail.")
	}

	// Wait until readEvents has reached unix.Read, and Close.
	<-time.After(50 * time.Millisecond)
	w.Close()

	// Wait for the close to complete.
	<-time.After(50 * time.Millisecond)
	isWatcherReallyClosed(t, w)
}

func TestNotifyCloseAfterRead(t *testing.T) {
	testDir := test.TempMkdir(t)
	defer os.RemoveAll(testDir)

	w, err := NewWatcher()
	if err != nil {
		t.Fatalf("Failed to create watcher")
	}

	var count, expected uint64
	ops := make([]EventOp, 5)
	names := make([]string, 5)
	handler := EventHandler{
		EventFilter: succcessFilterFunc,
		EventDispatch: func(name string, op EventOp) {
			atomic.AddUint64(&count, 1)
			names = append(names, name)
			ops = append(ops, op)
		},
	}

	err = w.Add(testDir, handler)
	if err != nil {
		t.Fatalf("Failed to add.")
	}

	w.PollAsyn()

	// Generate an event.
	os.Create(filepath.Join(testDir, "file"))
	expected += 1
	<-time.After(50 * time.Millisecond)
	checkEventNumber(t, count, expected)

	testSubDir := filepath.Join(testDir, "dir")
	os.Mkdir(testSubDir, 0755)
	expected += 1
	<-time.After(50 * time.Millisecond)
	checkEventNumber(t, count, expected)

	os.Create(filepath.Join(testSubDir, "subfile1"))
	<-time.After(50 * time.Millisecond)
	checkEventNumber(t, count, expected)

	subHandler := EventHandler{
		EventFilter: succcessFilterFunc,
		EventDispatch: func(name string, op EventOp) {
			atomic.AddUint64(&count, 1)
			names = append(names, name)
			ops = append(ops, op)
		},
	}

	err = w.Add(testSubDir, subHandler)
	if err != nil {
		t.Fatalf("Failed to add.")
	}

	//create
	f2, err := os.Create(filepath.Join(testSubDir, "subfile2"))
	expected += 1
	if err != nil {
		t.Fatalf("Failed to create subfile2.")
	}
	<-time.After(50 * time.Millisecond)
	checkEventNumber(t, count, expected)
	//write
	f2.WriteString("This is subfile2!\n")
	expected += 1
	<-time.After(50 * time.Millisecond)
	checkEventNumber(t, count, expected)

	//rename and create
	os.Rename(filepath.Join(testSubDir, "subfile2"), filepath.Join(testSubDir, "subfile2.bak"))
	expected += 2
	<-time.After(50 * time.Millisecond)
	checkEventNumber(t, count, expected)
	//chmod
	os.Chmod(filepath.Join(testSubDir, "subfile2.bak"), 0777)
	expected += 1
	<-time.After(50 * time.Millisecond)
	checkEventNumber(t, count, expected)
	//remove
	os.RemoveAll(filepath.Join(testSubDir, "subfile2.bak"))
	expected += 1
	<-time.After(50 * time.Millisecond)
	checkEventNumber(t, count, expected)

	err = w.Remove(testSubDir)
	if err != nil {
		t.Fatalf("Failed to remove.")
	}

	os.Create(filepath.Join(testSubDir, "subfile3"))
	<-time.After(50 * time.Millisecond)
	checkEventNumber(t, count, expected)

	// Wait for readEvents to read the event, then close the watcher.
	<-time.After(50 * time.Millisecond)
	w.Close()

	// Wait for the close to complete.
	<-time.After(50 * time.Millisecond)
	isWatcherReallyClosed(t, w)
}

func checkEventNumber(t *testing.T, count uint64, expected uint64) {
	if expected != count {
		t.Fatalf("event count is %d, which should be %d", count, expected)
	}
}

func succcessFilterFunc(name string, op EventOp) bool {
	return true
}

func failFilterFunc(name string, op EventOp) bool {
	return false
}

func isWatcherReallyClosed(t *testing.T, w *Watcher) {
	select {
	case err, ok := <-w.watcher.Errors:
		if ok {
			t.Fatalf("w.watcher.Errors is not closed; readEvents is still alive after closing (error: %v)", err)
		}
	default:
		t.Fatalf("w.watcher.Errors would have blocked; readEvents is still alive!")
	}

	select {
	case _, ok := <-w.watcher.Events:
		if ok {
			t.Fatalf("w.watcher.Events is not closed; readEvents is still alive after closing")
		}
	default:
		t.Fatalf("w.watcher.Events would have blocked; readEvents is still alive!")
	}
}
