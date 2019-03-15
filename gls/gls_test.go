package gls

import (
	"sync"
	"testing"
)

func TestGls(t *testing.T) {
	wg := sync.WaitGroup{}
	wg.Add(1)
	go func() {
		Initialize()
		defer func() {
			if "gls" != Get("name") {
				t.Errorf("the value of key 'name' should be 'gls'")
			}
			Clear()
			if Activated() {
				t.Errorf("gls did not be removed successfully")
			}
		}()

		if !Activated() {
			t.Errorf("gls should have been activated")
		}
		if nil != Get("name") {
			t.Errorf("key 'name' should not be found")
		}
		Set("name", "gls")
		if "gls" != Get("name") {
			t.Errorf("the value of key 'name' should be 'gls'")
		}

		wg.Done()
	}()
	wg.Wait()
	if Activated() {
		t.Errorf("parent gls should be not activated")
	}
}
