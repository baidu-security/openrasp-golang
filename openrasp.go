package openrasp

import (
	"log"
	"path/filepath"

	"github.com/baidu/openrasp/common"
	"github.com/baidu/openrasp/config"
)

var workSpace *common.WorkSpace
var commonGlobals *common.Globals
var basic *config.BasicConfig
var general *config.GeneralConfig
var logManager *LogManager

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

	confDir, err := workSpace.GetDir(common.Conf)
	if err != nil {
		log.Printf("%v", err)
	} else {
		path := filepath.Join(confDir, "openrasp.yml")
		err := basic.LoadProperties(path)
		if err != nil {
			log.Printf("%v", err)
		}
	}
	if !basic.GetBool("cloud.enable") {
		workSpace.StartWatch(common.Conf)
		workSpace.RegisterListener(common.Conf, general)
	}
}

func IsAvailable() bool {
	return commonGlobals != nil && basic != nil && general != nil
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
