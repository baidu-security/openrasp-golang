package orlog

import (
	"bytes"

	"github.com/baidu/openrasp/utils"
	"github.com/sirupsen/logrus"
)

type OpenRASPFormatter struct {
	TimestampFormat      string
	WithTimestamp        bool
	WithoutLineSeparator bool
}

func (f *OpenRASPFormatter) Format(entry *logrus.Entry) ([]byte, error) {

	var b *bytes.Buffer
	if entry.Buffer != nil {
		b = entry.Buffer
	} else {
		b = &bytes.Buffer{}
	}

	timestampFormat := f.TimestampFormat
	if timestampFormat == "" {
		timestampFormat = utils.ISO8601TimestampFormat
	}

	if f.WithTimestamp {
		timeString := entry.Time.Format(timestampFormat)
		b.WriteString(timeString)
		b.WriteByte(' ')
	}

	b.WriteString(entry.Message)
	if !f.WithoutLineSeparator {
		b.WriteByte('\n')
	}

	return b.Bytes(), nil
}
