package goid

import (
	"fmt"
	"runtime"
	"strings"
	"sync"
	"testing"
	_ "unsafe"
)

func TestGoIDAsm(t *testing.T) {
	var wg sync.WaitGroup
	for i := 0; i < 5; i++ {
		wg.Add(1)
		go func() {
			testGoIDAsm(t)
			wg.Done()
		}()
	}
	wg.Wait()
}

func TestGoIDStack(t *testing.T) {
	var wg sync.WaitGroup
	for i := 0; i < 5; i++ {
		wg.Add(1)
		go func() {
			testGoIDStack(t)
			wg.Done()
		}()
	}
	wg.Wait()
}

func TestCompareGoID(t *testing.T) {
	var wg sync.WaitGroup
	for i := 0; i < 5; i++ {
		wg.Add(1)
		go func() {
			idAsm := GoIDAsm()
			idStack := GoIDStack()
			if idAsm != idStack {
				t.Errorf("different goroutine id: GoIDAsm return %d but GoIDStack return %d", idAsm, idStack)
			}
			wg.Done()
		}()
	}
	wg.Wait()
}

func testGoIDAsm(t *testing.T) {
	id := GoIDAsm()
	lines := strings.Split(stackTrace(), "\n")
	for i, line := range lines {
		if !strings.HasPrefix(line, fmt.Sprintf("goroutine %d ", id)) {
			continue
		}
		if i+1 == len(lines) {
			break
		}
		if !strings.Contains(lines[i+1], ".stackTrace") {
			t.Errorf("there are goroutine id %d but it is not me: %s", id, lines[i+1])
		}
		return
	}
	t.Errorf("there are no goroutine %d", id)
}

func testGoIDStack(t *testing.T) {
	id := GoIDStack()
	lines := strings.Split(stackTrace(), "\n")
	for i, line := range lines {
		if !strings.HasPrefix(line, fmt.Sprintf("goroutine %d ", id)) {
			continue
		}
		if i+1 == len(lines) {
			break
		}
		if !strings.Contains(lines[i+1], ".stackTrace") {
			t.Errorf("there are goroutine id %d but it is not me: %s", id, lines[i+1])
		}
		return
	}
	t.Errorf("there are no goroutine %d", id)
}

func stackTrace() string {
	var n int
	for n = 4096; n < 16777216; n *= 2 {
		buf := make([]byte, n)
		ret := runtime.Stack(buf, true)
		if ret != n {
			return string(buf[:ret])
		}
	}
	panic(n)
}
