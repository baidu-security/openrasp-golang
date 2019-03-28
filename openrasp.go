package openrasp

import (
	"log"
	"path/filepath"

	"github.com/baidu-security/openrasp-golang/common"
	"github.com/baidu-security/openrasp-golang/config"
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
		log.Printf("Unable to init v8.")
		return
	}

	confDir, err := workSpace.GetDir(common.Conf)
	if err != nil {
		log.Printf("%v", err)
		return
	}

	pluginDir, err := workSpace.GetDir(common.Plugins)
	if err != nil {
		log.Printf("%v", err)
		return
	}
	pluginManager = NewPluginManager(pluginDir)

	buildinAction = NewBuildinAction()
	pluginManager.AttachListener(buildinAction)

	complete = true

	yamlPath := filepath.Join(confDir, "openrasp.yml")
	err = basic.LoadYaml(yamlPath)
	if err != nil {
		log.Printf("%v", err)
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
