package common

import (
	"errors"
	"os"
	"path/filepath"
	"sync"
	"syscall"

	"github.com/baidu/openrasp/notify"
)

type WorkDirCode int

const (
	Root WorkDirCode = iota
	Assets
	Conf
	Locale
	Logs
	Plugins
	LogAlarm
	LogPlugin
	LogPolicy
	LogRasp
)

type NotifyListener interface {
	OnUpdate(absPath string)
}

type WorkSpace struct {
	workDirs WorkDirInfo
	dirMap   map[WorkDirCode]*codeInfo
	mu       sync.RWMutex
	watcher  *notify.Watcher
	pollOnce sync.Once
	active   bool
}

type WorkDirInfo struct {
	name string
	mode os.FileMode
	code WorkDirCode
	subs []WorkDirInfo
}

type codeInfo struct {
	path      string
	listeners []NotifyListener
}

func NewWorkSpace() *WorkSpace {
	newWatcher, _ := notify.NewWatcher()
	return &WorkSpace{
		workDirs: WorkDirInfo{
			"rasp",
			0755,
			Root,
			[]WorkDirInfo{
				WorkDirInfo{
					"assets",
					0755,
					Assets,
					[]WorkDirInfo{},
				},
				WorkDirInfo{
					"conf",
					0755,
					Conf,
					[]WorkDirInfo{},
				},
				WorkDirInfo{
					"locale",
					0755,
					Locale,
					[]WorkDirInfo{},
				},
				WorkDirInfo{
					"logs",
					0755,
					Logs,
					[]WorkDirInfo{
						WorkDirInfo{
							"alarm",
							0755,
							LogAlarm,
							[]WorkDirInfo{},
						},
						WorkDirInfo{
							"plugin",
							0755,
							LogPlugin,
							[]WorkDirInfo{},
						},
						WorkDirInfo{
							"policy",
							0755,
							LogPolicy,
							[]WorkDirInfo{},
						},
						WorkDirInfo{
							"rasp",
							0755,
							LogRasp,
							[]WorkDirInfo{},
						},
					},
				},
				WorkDirInfo{
					"plugins",
					0755,
					Plugins,
					[]WorkDirInfo{},
				},
			},
		},
		dirMap:  make(map[WorkDirCode]*codeInfo),
		watcher: newWatcher,
		active:  false,
	}
}

func newCodeInfo(absPath string) *codeInfo {
	return &codeInfo{
		absPath,
		make([]NotifyListener, 0)}
}

func (ci *codeInfo) getPath() string {
	return ci.path
}

func (ci *codeInfo) attachListener(listener NotifyListener) {
	ci.listeners = append(ci.listeners, listener)
}

func (ci *codeInfo) notify(absPath string) {
	for _, o := range ci.listeners {
		if o != nil {
			o.OnUpdate(absPath)
		}
	}
}

func (wdi *WorkDirInfo) init(base string, registerFunc func(code WorkDirCode, absPath string)) error {
	currentDir := filepath.Join(base, wdi.name)
	if _, err := os.Stat(currentDir); err == nil {
		err := syscall.Access(currentDir, syscall.O_RDWR)
		if err != nil {
			return err
		} else {
			registerFunc(wdi.code, currentDir)
		}
	} else if os.IsNotExist(err) {
		err := os.Mkdir(currentDir, wdi.mode)
		if err != nil {
			return err
		} else {
			registerFunc(wdi.code, currentDir)
		}
	} else {
		return err
	}
	for _, v := range wdi.subs {
		err := v.init(currentDir, registerFunc)
		if err != nil {
			return err
		}
	}
	return nil
}

func (wdi *WorkDirInfo) clear(base string) error {
	currentDir := filepath.Join(base, wdi.name)
	return os.RemoveAll(currentDir)
}

func (ws *WorkSpace) Init() error {
	executableDir, err := getExecutableDir()
	if err != nil {
		return err
	} else {
		err := ws.workDirs.init(executableDir, ws.register)
		if err != nil {
			return err
		} else {
			ws.active = true
			return nil
		}
	}
}

func (ws *WorkSpace) clear() error {
	executableDir, err := getExecutableDir()
	if err != nil {
		return err
	} else {
		err := ws.workDirs.clear(executableDir)
		if err != nil {
			return err
		} else {
			ws.active = false
			return nil
		}
	}
}

func (ws *WorkSpace) register(code WorkDirCode, absPath string) {
	ws.mu.Lock()
	defer ws.mu.Unlock()
	ws.dirMap[code] = newCodeInfo(absPath)
}

func (ws *WorkSpace) GetDir(code WorkDirCode) (string, error) {
	if !ws.Active() {
		return "", errors.New("work space is not active.")
	}
	ws.mu.RLock()
	defer ws.mu.RUnlock()
	v, ok := ws.dirMap[code]
	if ok {
		return v.getPath(), nil
	} else {
		return "", errors.New("invalid code parameter.")
	}
}

func (ws *WorkSpace) Active() bool {
	return ws.active
}

func (ws *WorkSpace) StartWatch(codes ...WorkDirCode) bool {
	if ws.watcher == nil {
		return false
	}
	ws.pollOnce.Do(func() {
		ws.watcher.PollAsyn()
	})
	wtachCodes := []WorkDirCode(codes)
	for _, code := range wtachCodes {
		handler := notify.EventHandler{
			EventFilter:   filterNone,
			EventDispatch: buildDispatchFunc(ws, code),
		}
		path, err := ws.GetDir(code)
		if err != nil {
			return false
		}
		err = ws.watcher.Add(path, handler)
		if err != nil {
			return false
		}
	}
	return true
}

func (ws *WorkSpace) RegisterListener(code WorkDirCode, listener NotifyListener) {
	ws.mu.Lock()
	defer ws.mu.Unlock()
	ws.dirMap[code].attachListener(listener)
}

func filterNone(name string, op notify.EventOp) bool {
	return true
}

func buildDispatchFunc(ws *WorkSpace, code WorkDirCode) notify.EventDispatchFunc {
	return func(name string, op notify.EventOp) {
		ws.mu.Lock()
		defer ws.mu.Unlock()
		ws.dirMap[code].notify(name)
	}
}

func getExecutableDir() (string, error) {
	ex, err := os.Executable()
	if err != nil {
		return "", err
	}
	return filepath.Dir(ex), nil
}
