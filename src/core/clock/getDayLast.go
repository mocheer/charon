package clock

import "time"

// GetNextDay 获取明天对应的时间对象
func GetNextDay() time.Time {
	return time.Now().AddDate(0, 0, 1)
}

// GetDayLastSecond 获取一天内剩余的时间，秒数
func GetDayLastSecond() int64 {
	return GetNextDay().Unix() - time.Now().Unix()
}

// GetDayLastMillisecond 获取一天内剩余的时间，毫秒数
func GetDayLastMillisecond() int64 {
	return (GetNextDay().UnixNano() - time.Now().UnixNano()) / 1e6
}
