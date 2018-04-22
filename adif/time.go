package adif

import "time"

func NowUTCDate() string {
	t := time.Now().In(time.UTC)
	return t.Format("20060102")
}

func NowUTCTime() string {
	t := time.Now().In(time.UTC)
	return t.Format("1504")
}

func NowUTCTimestamp() string {
	t := time.Now().In(time.UTC)
	return t.Format("20060102 15:04")
}
