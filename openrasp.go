package openrasp

import (
	"log"
	"path/filepath"

	"github.com/baidu-security/openrasp-golang/common"
	"github.com/baidu-security/openrasp-golang/config"
	"github.com/baidu-security/openrasp-golang/orlog"
	v8 "github.com/baidu-security/openrasp-v8/go"
)

var workSpace *common.WorkSpace
var commonGlobals *common.Globals
var basic *config.BasicConfig
var general *config.GeneralConfig
var logManager *LogManager
var pluginManager *PluginManager
var whiteList *WhiteList
var buildinAction *BuildinAction
var complete bool

func init() {
	workSpace = common.NewWorkSpace()
	workSpace.Init()
	if !workSpace.Active() {
		log.Printf("Fail to init workspace.")
		return
	}

	rootDir, err := workSpace.GetDir(common.Root)
	if err != nil {
		log.Printf("Unable to get root dir, cuz of %v", err)
		return
	}

	commonGlobals = common.NewGlobals(rootDir)
	basic = config.NewBasicConfig()
	general = config.NewGeneralConfig()

	logManager, err = InitLogManager()
	if err != nil {
		log.Printf("Unable to init log manager, cuz of %v", err)
		return
	}
	logManager.UpdateFileWriter()
	GetGeneral().AttachListener(logManager)

	whiteList = NewWhiteList()
	GetGeneral().AttachListener(whiteList)

	if !v8.Initialize(logManager.PluginInfo) {
		GetLog().RaspWarn("Unable to initialize v8.", orlog.Plugin)
		return
	} else {
		GetLog().RaspDebug("Initialize v8 successfully.", orlog.Plugin)
	}

	confDir, err := workSpace.GetDir(common.Conf)
	if err != nil {
		GetLog().RaspWarn(err.Error(), orlog.Config)
		return
	}

	pluginDir, err := workSpace.GetDir(common.Plugins)
	if err != nil {
		GetLog().RaspWarn(err.Error(), orlog.Config)
		return
	}
	pluginManager = NewPluginManager(pluginDir)

	buildinAction = NewBuildinAction()
	pluginManager.AttachListener(buildinAction)

	yamlPath := filepath.Join(confDir, "openrasp.yml")
	err = basic.LoadYaml(yamlPath)
	if err != nil {
		GetLog().RaspWarn(err.Error(), orlog.Log)
	}

	if !basic.GetBool("cloud.enable") {
		general.LoadYaml(yamlPath)
		pluginManager.buildLocalSnapshot()
		workSpace.StartWatch(common.Conf)
		workSpace.RegisterListener(common.Conf, general)
		workSpace.StartWatch(common.Plugins)
		workSpace.RegisterListener(common.Plugins, pluginManager)
	}
	InitContextGetters()

	complete = true
	GetLog().RaspInfo("Initialize OpenRASP successfully.", orlog.Runtime)
}

func IsComplete() bool {
	return complete
}

func GetWorkSpace() *common.WorkSpace {
	return workSpace
}

func GetGlobals() *common.Globals {
	return commonGlobals
}

func GetBasic() *config.BasicConfig {
	return basic
}

func GetGeneral() *config.GeneralConfig {
	return general
}

func GetLog() *LogManager {
	return logManager
}

func GetWhite() *WhiteList {
	return whiteList
}

func GetAction() *BuildinAction {
	return buildinAction
}
