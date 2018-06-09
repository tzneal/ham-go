package db

import (
	"strings"
	"time"
)

const timestampFormat = "20060102 15:04"

func TimeToUTCString(t time.Time) string {
	t = t.In(time.UTC)
	return t.Format(timestampFormat)
}

func NormalizeCall(call string) string {
	call = strings.TrimSpace(call)
	return strings.ToUpper(call)
}

func UTCStringToTime(s string) (time.Time, error) {
	return time.Parse(timestampFormat, s)
}
