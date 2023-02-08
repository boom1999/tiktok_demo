package util

import "time"

func GetCurrentTimeMillisecond() int64 {
	return time.Now().UnixNano() / 1e6
}
