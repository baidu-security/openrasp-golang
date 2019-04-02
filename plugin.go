package openrasp

import (
	"errors"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
	"sync"

	v8 "github.com/baidu-security/openrasp-v8/go"
)

type UpdateListener interface {
	OnPluginUpdate()
}

type PluginManager struct {
	dirPath   string
	plugins   []v8.Plugin
	listeners []UpdateListener
	mu        sync.RWMutex
}

func NewPluginManager(dir string) *PluginManager {
	pm := &PluginManager{
		dirPath: dir,
	}
	return pm
}

func (pm *PluginManager) AttachListener(listener UpdateListener) {
	pm.mu.Lock()
	defer pm.mu.Unlock()
	pm.listeners = append(pm.listeners, listener)
}

func (pm *PluginManager) walkFunc(path string, info os.FileInfo, err error) error {
	if !info.IsDir() && filepath.Ext(path) == ".js" {
		plugin, err := newPlugin(path)
		if err == nil {
			pm.plugins = append(pm.plugins, *plugin)
		}
	}
	return nil
}

func (pm *PluginManager) loadLocalPlugins() {
	pm.mu.Lock()
	defer pm.mu.Unlock()
	pm.plugins = nil
	filepath.Walk(pm.dirPath, pm.walkFunc)
}

func (pm *PluginManager) buildLocalSnapshot() {
	pm.loadLocalPlugins()
	pm.createSnapshot()
}

func (pm *PluginManager) createSnapshot() {
	pm.mu.RLock()
	defer pm.mu.RUnlock()
	if len(pm.plugins) > 0 {
		if v8.CreateSnapshot("", pm.plugins) {
			for _, l := range pm.listeners {
				l.OnPluginUpdate()
			}
		}
	}
}

func (pm *PluginManager) OnUpdate(absPath string) {
	if filepath.Ext(absPath) == ".js" {
		pm.buildLocalSnapshot()
	}
}

func (pm *PluginManager) OnUpdateCloud(source string, filename string) {
	pm.plugins = []v8.Plugin{v8.Plugin{
		Source:   source,
		Filename: filename,
	}}
	pm.createSnapshot()
}

func newPlugin(path string) (*v8.Plugin, error) {
	if filepath.Ext(path) == ".js" {
		content, err := ioutil.ReadFile(path)
		if err != nil {
			return nil, err
		}
		fileWithExt := filepath.Base(path)
		filename := strings.TrimSuffix(fileWithExt, filepath.Ext(fileWithExt))
		p := &v8.Plugin{
			Source:   string(content),
			Filename: filename,
		}
		return p, nil
	} else {
		return nil, errors.New("plugin must be javascript file")
	}
}
