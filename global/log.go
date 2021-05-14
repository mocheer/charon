package global

import (
	"fmt"
	"os"
	"path"
	"regexp"
	"time"

	"github.com/mocheer/charon/cts"
	"github.com/mocheer/pluto/clock"
	"github.com/mocheer/pluto/fs"
	"github.com/mocheer/pluto/js/window"
	"github.com/sirupsen/logrus"
)

type Logger struct {
	*logrus.Logger
}

// Log 日志操作对象
var Log = &Logger{logrus.New()}

// Init 初始化日志
func (log *Logger) Init() {
	logFileName := path.Join(cts.LogDir, fmt.Sprintf("/logrus.%s.log", clock.Now().Fmt(clock.FmtCompactDate)))
	file, err := fs.OpenOrCreate(logFileName, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	if err == nil {
		log.Out = file
	} else {
		Log.Error("打开日志文件失败", err)
	}
	// 每隔一天重新切换日志文件输出
	window.SetTimeout(log.Init, time.Duration(clock.GetDayLastMillisecond()))
}

// CheckClear 历史的日志文件需要删除
func (log *Logger) CheckClear(num float64) {
	fs.EachDirToRemove(cts.LogDir, func(filename string) bool {
		timeStrArray := regexp.MustCompile(`logrus\.(.+)\.log`).FindStringSubmatch(filename)
		if len(timeStrArray) == 0 {
			return true
		}
		return clock.MustParse(timeStrArray[1], clock.FmtCompactDate).SinceDays() > num
	})
}
