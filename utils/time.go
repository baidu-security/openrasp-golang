package utils

import "time"

const (
	ISO8601TimestampFormat = "2006-01-02T15:04:05-0700"
)

func CurrentISO8601Time() string {
	return time.Now().Format(ISO8601TimestampFormat)
}
