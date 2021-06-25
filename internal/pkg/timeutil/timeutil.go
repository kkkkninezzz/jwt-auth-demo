package timeutil

import "time"

// 获取当前时间的秒数  UNIX时间戳
func CurrentTimeSeconds() int64 {
	return time.Now().Unix()
}

// 获取当前时间的纳秒数
func CurrentTimeNanos() int64 {
	return time.Now().UnixNano()
}

// 获取当前时间的豪秒数
func CurrentTimeMillis() int64 {
	return time.Now().UnixNano() / 1e9
}

// 获取当前时间的下一个时间，增量为d
func NextTimeSeconds(d time.Duration) int64 {
	return time.Now().Add(d).Unix()
}
