package clock

import "time"

// Clock 结构体
type Clock struct {
	date time.Time
}

// New 实例化
func New() *Clock {
	return &Clock{date: time.Now()}
}

// Parse 解析时间字符串
func Parse(s string) (time.Time, error) {
	return time.Parse(FullDateFomate, s)
}

// Format 格式化返回时间字符串
func (c *Clock) Format(layout string) string {
	if layout == "" {
		layout = DateFormat
	}
	return c.date.Format(layout)
}

// Value 获取时间戳
func (c *Clock) Value() int64 {
	return c.date.Unix()
}
