package orlog

import (
	"fmt"
	"os"

	"github.com/baidu-security/openrasp-golang/cloud"

	"github.com/sirupsen/logrus"
)

type Hook interface {
	Levels() []logrus.Level
	Fire(*logrus.Entry) error
}

type HttpHook struct {
	hookLevel Level
	Writer    *HttpWriter
}

func NewHttpHook(t string, cm *cloud.Client, level Level, tokenBucket *TokenBucket) *HttpHook {
	hw := NewHttpWriter(t, cm, tokenBucket)
	hh := &HttpHook{
		hookLevel: level,
		Writer:    hw,
	}
	return hh
}

func (hook *HttpHook) Fire(entry *logrus.Entry) error {
	line, err := entry.String()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Unable to read entry, %v", err)
		return err
	}
	_, err = hook.Writer.Write([]byte("[\n" + line + "]"))
	return err
}

func (hook *HttpHook) Levels() []logrus.Level {
	switch hook.hookLevel {
	case WarnLevel:
		return []logrus.Level{logrus.WarnLevel}
	default:
		return []logrus.Level{logrus.InfoLevel}
	}
}
