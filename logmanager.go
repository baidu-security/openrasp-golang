package openrasp

import (
	"io"
	"path/filepath"
	"time"

	"github.com/baidu/openrasp/common"
	"github.com/baidu/openrasp/orlog"
	"github.com/sirupsen/logrus"
)

type LogCode int

const (
	refillDuration = 60 * 1000 * 1000
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

func NewWrapLogger(dirCode common.WorkDirCode) (*WrapLogger, error) {
	var logFilename string
	logDir, err := workSpace.GetDir(dirCode)
	if err != nil {
		return nil, err
	} else {
		logFilename = filepath.Join(logDir, dirCodeToName(dirCode))
		logrusLogger := logrus.New()
		logrusLogger.Formatter = &orlog.OpenRASPFormatter{}
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
	alarmLogger, err := NewWrapLogger(common.LogAlarm)
	if err != nil {
		return nil, err
	}
	policyLogger, err := NewWrapLogger(common.LogPolicy)
	if err != nil {
		return nil, err
	}
	pluginLogger, err := NewWrapLogger(common.LogPlugin)
	if err != nil {
		return nil, err
	}
	raspLogger, err := NewWrapLogger(common.LogRasp)
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

func (lm *LogManager) UpdateFileWriter() {
	maxBackup := GetGeneral().GetInt("log.maxbackup")
	capacity := GetGeneral().GetInt64("log.maxburst")
	lm.alarm.SetOutput(orlog.NewLogger(lm.alarm.filename, maxBackup, orlog.NewTokenBucket(uint64(capacity), time.Duration(refillDuration))))
	lm.policy.SetOutput(orlog.NewLogger(lm.policy.filename, maxBackup, orlog.NewTokenBucket(uint64(capacity), time.Duration(refillDuration))))
	lm.plugin.SetOutput(orlog.NewLogger(lm.plugin.filename, maxBackup, orlog.NewTokenBucket(uint64(capacity), time.Duration(refillDuration))))
	lm.rasp.SetOutput(orlog.NewLogger(lm.rasp.filename, maxBackup, orlog.NewTokenBucket(uint64(capacity), time.Duration(refillDuration))))
}

func (lm *LogManager) OnConfigUpdate() {
	lm.UpdateFileWriter()
}

func (lm *LogManager) PolicyInfo(message string) {
	lm.GetPolicy().Info(message)
}
