package common

import (
	"errors"
	"sync"

	"github.com/baidu-security/openrasp-golang/notify"
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
	workDirs *WorkDirInfo
	dirMap   map[WorkDirCode]*codeInfo
	mu       sync.RWMutex
	watcher  *notify.Watcher
	pollOnce sync.Once
	active   bool
}

type codeInfo struct {
	workDir   *WorkDirInfo
	listeners []NotifyListener
}

func NewWorkSpace(baseDir string) *WorkSpace {
	newWatcher, _ := notify.NewWatcher()
	workDirMap := make(map[WorkDirCode]*codeInfo)

	raspDir := NewWorkDirInfo(baseDir, "rasp", 0755)
	workDirMap[Root] = newCodeInfo(raspDir)

	assetsDir := NewWorkDirInfo(raspDir.absPath(), "assets", 0755)
	workDirMap[Assets] = newCodeInfo(assetsDir)
	raspDir.appendSubDir(assetsDir)

	confDir := NewWorkDirInfo(raspDir.absPath(), "conf", 0777)
	workDirMap[Conf] = newCodeInfo(confDir)
	raspDir.appendSubDir(confDir)

	localeDir := NewWorkDirInfo(raspDir.absPath(), "locale", 0755)
	workDirMap[Locale] = newCodeInfo(localeDir)
	raspDir.appendSubDir(localeDir)

	logsDir := NewWorkDirInfo(raspDir.absPath(), "logs", 0755)
	workDirMap[Logs] = newCodeInfo(logsDir)
	raspDir.appendSubDir(logsDir)

	logAlarmDir := NewWorkDirInfo(logsDir.absPath(), "alarm", 0777)
	workDirMap[LogAlarm] = newCodeInfo(logAlarmDir)
	logsDir.appendSubDir(logAlarmDir)

	logPluginDir := NewWorkDirInfo(logsDir.absPath(), "plugin", 0777)
	workDirMap[LogPlugin] = newCodeInfo(logPluginDir)
	logsDir.appendSubDir(logPluginDir)

	logPolicyDir := NewWorkDirInfo(logsDir.absPath(), "policy", 0777)
	workDirMap[LogPolicy] = newCodeInfo(logPolicyDir)
	logsDir.appendSubDir(logPolicyDir)

	logRaspDir := NewWorkDirInfo(logsDir.absPath(), "rasp", 0777)
	workDirMap[LogRasp] = newCodeInfo(logRaspDir)
	logsDir.appendSubDir(logRaspDir)

	pluginsDir := NewWorkDirInfo(raspDir.absPath(), "plugins", 0777)
	workDirMap[Plugins] = newCodeInfo(pluginsDir)
	raspDir.appendSubDir(pluginsDir)

	return &WorkSpace{
		workDirs: raspDir,
		dirMap:   workDirMap,
		watcher:  newWatcher,
		active:   false,
	}
}

func newCodeInfo(workDir *WorkDirInfo) *codeInfo {
	return &codeInfo{
		workDir,
		make([]NotifyListener, 0)}
}

func (ci *codeInfo) getPath() string {
	return ci.workDir.absPath()
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

func (ws *WorkSpace) Init() error {
	err := ws.workDirs.init()
	if err != nil {
		return err
	} else {
		ws.active = true
		return nil
	}
}

func (ws *WorkSpace) clear() error {
	err := ws.workDirs.clear()
	if err != nil {
		return err
	} else {
		ws.active = false
		return nil
	}
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
