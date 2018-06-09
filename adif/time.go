package adif

import "time"

// UTCDate returns the specified date in ADIF UTC format
func UTCDate(t time.Time) string {
	return t.In(time.UTC).Format("20060102")
}

// NowUTCDate returns the current date in ADIF UTC format
func NowUTCDate() string {
	return UTCDate(time.Now())
}

// UTCTime returns the specified time in ADIF UTC format
func UTCTime(t time.Time) string {
	return t.In(time.UTC).Format("1504")
}

// NowUTCTime returns the current time in ADIF UTC format
func NowUTCTime() string {
	return UTCTime(time.Now())
}

// NowUTCTimestamp returns the current date/time in ADIF UTC format
func NowUTCTimestamp() string {
	return UTCTimestamp(time.Now())
}

func UTCTimestamp(t time.Time) string {
	return t.In(time.UTC).Format("20060102 15:04")
}
