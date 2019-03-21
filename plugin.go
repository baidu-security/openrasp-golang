package openrasp

import (
	"errors"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"

	v8 "github.com/baidu-security/openrasp-v8/go"
)

type PluginManager struct {
	dirPath string
	plugins []string
}

func NewPluginManager(dir string) *PluginManager {
	pm := &PluginManager{
		dirPath: dir,
		plugins: []string{},
	}
	return pm
}

func (pm *PluginManager) walkFunc(path string, info os.FileInfo, err error) error {
	if !info.IsDir() && filepath.Ext(path) == ".js" {
		pm.plugins = append(pm.plugins, path)
	}
	return nil
}

func (pm *PluginManager) clearPlugins() {
	pm.plugins = nil
}

func (pm *PluginManager) searchPlugins() {
	filepath.Walk(pm.dirPath, pm.walkFunc)
}

func (pm *PluginManager) LoadLocalPlugins() {
	pm.clearPlugins()
	pm.searchPlugins()
	pm.buildSnapshot()
}

func (pm *PluginManager) buildSnapshot() {
	var plugins []v8.Plugin
	for _, pluginFile := range pm.plugins {
		plugin, err := NewPlugin(pluginFile)
		if err == nil {
			plugins = append(plugins, *plugin)
		}
	}
	if len(plugins) > 0 {
		v8.CreateSnapshot("", plugins)
	}
}

func (pm *PluginManager) OnUpdate(absPath string) {
	if filepath.Ext(absPath) == ".js" {
		pm.LoadLocalPlugins()
	}
}

func NewPlugin(path string) (*v8.Plugin, error) {
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
