package openrasp

import (
	"io"
	"path/filepath"
	"time"

	"github.com/baidu-security/openrasp-golang/common"
	"github.com/baidu-security/openrasp-golang/orlog"
	"github.com/baidu-security/openrasp-golang/utils"
	"github.com/sirupsen/logrus"
)

type LogCode int

const (
	refillDuration = 60 * 1000 * 1000
	duration       = time.Duration(refillDuration)
)

type LogManager struct {
	alarm  *WrapLogger
	policy *WrapLogger
	plugin *WrapLogger
	rasp   *WrapLogger
}

type WrapLogger struct {
	logger   *logrus.Logger
	filename string
	dirCode  common.WorkDirCode
}

func NewWrapLogger(dirCode common.WorkDirCode, f *orlog.OpenRASPFormatter) (*WrapLogger, error) {
	var logFilename string
	logDir, err := workSpace.GetDir(dirCode)
	if err != nil {
		return nil, err
	} else {
		logFilename = filepath.Join(logDir, dirCodeToName(dirCode))
		logrusLogger := logrus.New()
		logrusLogger.Formatter = f
		wl := &WrapLogger{
			logger:   logrusLogger,
			filename: logFilename,
			dirCode:  dirCode,
		}
		return wl, nil
	}
}

func (wl *WrapLogger) Info(message string) {
	wl.logger.Info(message)
}

func (wl *WrapLogger) SetOutput(output io.Writer) {
	wl.logger.SetOutput(output)
}

func (wl *WrapLogger) SetLevel(l orlog.Level) {
	wl.logger.SetLevel(orlog.LevelTransform(l))
}

func (wl *WrapLogger) ClearHooks() {
	wl.logger.ReplaceHooks(make(logrus.LevelHooks))
}

func (wl *WrapLogger) AddHook(hook orlog.Hook) {
	wl.logger.AddHook(hook)
}

func dirCodeToName(dirCode common.WorkDirCode) string {
	switch dirCode {
	case common.LogAlarm:
		return "alarm.log"
	case common.LogPolicy:
		return "policy.log"
	case common.LogPlugin:
		return "plugin.log"
	case common.LogRasp:
		return "rasp.log"
	default:
		return ""
	}
}

func InitLogManager() (*LogManager, error) {
	alarmLogger, err := NewWrapLogger(common.LogAlarm, &orlog.OpenRASPFormatter{})
	if err != nil {
		return nil, err
	}
	policyLogger, err := NewWrapLogger(common.LogPolicy, &orlog.OpenRASPFormatter{})
	if err != nil {
		return nil, err
	}
	pluginLogger, err := NewWrapLogger(common.LogPlugin, &orlog.OpenRASPFormatter{
		TimestampFormat:      utils.ISO8601TimestampFormat,
		WithTimestamp:        true,
		WithoutLineSeparator: true,
	})
	if err != nil {
		return nil, err
	}
	raspLogger, err := NewWrapLogger(common.LogRasp, &orlog.OpenRASPFormatter{})
	if err != nil {
		return nil, err
	}
	lm := &LogManager{
		alarm:  alarmLogger,
		policy: policyLogger,
		plugin: pluginLogger,
		rasp:   raspLogger,
	}
	return lm, nil
}

func (lm *LogManager) GetPolicy() *WrapLogger {
	return lm.policy
}

func (lm *LogManager) GetPlugin() *WrapLogger {
	return lm.plugin
}

func (lm *LogManager) GetAlarm() *WrapLogger {
	return lm.alarm
}

func (lm *LogManager) UpdateFileWriter() {
	maxBackup := GetGeneral().GetInt("log.maxbackup")
	capacity := GetGeneral().GetInt64("log.maxburst")
	lm.alarm.SetOutput(orlog.NewFileWriter(lm.alarm.filename, maxBackup, orlog.NewTokenBucket(uint64(capacity), duration)))
	lm.policy.SetOutput(orlog.NewFileWriter(lm.policy.filename, maxBackup, orlog.NewTokenBucket(uint64(capacity), duration)))
	lm.plugin.SetOutput(orlog.NewFileWriter(lm.plugin.filename, maxBackup, orlog.NewTokenBucket(uint64(capacity), duration)))
	lm.rasp.SetOutput(orlog.NewFileWriter(lm.rasp.filename, maxBackup, orlog.NewTokenBucket(uint64(capacity), duration)))
	debugLevel := GetGeneral().GetInt("debug.level")
	if debugLevel > 0 {
		lm.rasp.SetLevel(orlog.DebugLevel)
	}
}

func (lm *LogManager) UpdateHttpHook() {
	backendUrl := GetBasic().GetString("cloud.backend_url")
	appId := GetBasic().GetString("cloud.app_id")
	appSecret := GetBasic().GetString("cloud.app_secret")
	capacity := GetGeneral().GetInt64("log.maxburst")
	lm.alarm.ClearHooks()
	lm.alarm.AddHook(orlog.NewHttpHook(backendUrl, appId, appSecret, orlog.InfoLevel, orlog.NewTokenBucket(uint64(capacity), duration)))
	lm.policy.ClearHooks()
	lm.policy.AddHook(orlog.NewHttpHook(backendUrl, appId, appSecret, orlog.InfoLevel, orlog.NewTokenBucket(uint64(capacity), duration)))
	lm.rasp.ClearHooks()
	lm.rasp.AddHook(orlog.NewHttpHook(backendUrl, appId, appSecret, orlog.WarnLevel, orlog.NewTokenBucket(uint64(capacity), duration)))
}

func (lm *LogManager) OnConfigUpdate() {
	lm.UpdateFileWriter()
	clouldEnable := GetBasic().GetBool("cloud.enable")
	if clouldEnable {
		lm.UpdateHttpHook()
	}
}

func (lm *LogManager) PolicyInfo(message string) {
	lm.GetPolicy().Info(message)
}

func (lm *LogManager) PluginInfo(message string) {
	lm.GetPlugin().Info(message)
}

func (lm *LogManager) AlarmInfo(message string) {
	lm.GetAlarm().Info(message)
}
