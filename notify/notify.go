package notify

import (
	"fmt"
	"log"
	"path/filepath"
	"sync"

	"github.com/fsnotify/fsnotify"
)

type EventOp int

const (
	Invalid EventOp = 0
	Create  EventOp = 1 << 0
	Write   EventOp = 1 << 1
	Remove  EventOp = 1 << 2
	Rename  EventOp = 1 << 3
	Chmod   EventOp = 1 << 4
	All             = Create | Write | Remove | Rename | Chmod
)

type Watcher struct {
	watcher    *fsnotify.Watcher
	handlerMap HandlerMap
	mu         sync.RWMutex
}

type HandlerMap map[string]EventHandler

type EventHandler struct {
	EventFilter   EventFilterFunc
	EventDispatch EventDispatchFunc
}

type EventFilterFunc func(name string, op EventOp) bool

type EventDispatchFunc func(name string, op EventOp)

func NewWatcher() (*Watcher, error) {
	fw, err := fsnotify.NewWatcher()
	if err != nil {
		return nil, err
	}
	watcher := &Watcher{
		watcher:    fw,
		handlerMap: make(HandlerMap),
	}
	return watcher, nil
}

func (w *Watcher) Add(name string, handler EventHandler) error {
	w.mu.Lock()
	defer w.mu.Unlock()
	_, ok := w.handlerMap[name]
	if ok {
		return fmt.Errorf("%s has already added.", name)
	}
	w.handlerMap[name] = handler
	return w.watcher.Add(name)
}

func (w *Watcher) Remove(name string) error {
	err := w.watcher.Remove(name)
	w.mu.Lock()
	defer w.mu.Unlock()
	delete(w.handlerMap, name)
	return err
}

func (w *Watcher) Close() error {
	err := w.watcher.Close()
	w.mu.Lock()
	defer w.mu.Unlock()
	for k := range w.handlerMap {
		delete(w.handlerMap, k)
	}
	return err
}

func (w *Watcher) PollAsyn() {
	go func() {
		w.poll()
	}()
}

func (w *Watcher) poll() {
	for {
		var eop EventOp = Invalid
		select {
		case event, ok := <-w.watcher.Events:
			if !ok {
				return
			}
			if event.Op&fsnotify.Write == fsnotify.Write {
				eop |= Write
			} else if event.Op&fsnotify.Create == fsnotify.Create {
				eop |= Create
			} else if event.Op&fsnotify.Remove == fsnotify.Remove {
				eop |= Remove
			} else if event.Op&fsnotify.Rename == fsnotify.Rename {
				eop |= Rename
			} else if event.Op&fsnotify.Chmod == fsnotify.Chmod {
				eop |= Chmod
			}
			w.mu.RLock()
			for k, v := range w.handlerMap {
				if eventMatchWatch(k, event.Name) && v.EventFilter(event.Name, eop) {
					v.EventDispatch(event.Name, eop)
				}
			}
			w.mu.RUnlock()
		case err, ok := <-w.watcher.Errors:
			if !ok {
				return
			}
			log.Println("error : ", err)
		}
	}
}

func eventMatchWatch(watchPath string, eventName string) bool {
	rel, err := filepath.Rel(eventName, watchPath)
	if nil == err {
		return "." == rel || ".." == rel
	}
	return false
}
