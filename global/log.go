package global

import (
	"fmt"
	"os"
	"path"
	"regexp"
	"time"

	"github.com/mocheer/charon/cts"
	"github.com/mocheer/pluto/fs"
	"github.com/mocheer/pluto/js"
	"github.com/mocheer/pluto/ts/clock"
	"github.com/sirupsen/logrus"
)

type Logger struct {
	*logrus.Logger
	// 日志文件保存时长，超过时长后自动清除
	KeepDay float64
}

// Log 日志操作对象
var Log = &Logger{Logger: logrus.New(), KeepDay: 10}

// Init 初始化日志
func (log *Logger) Init() {
	log.ClearOld()
	logFileName := path.Join(cts.LogDir, fmt.Sprintf("/logrus.%s.log", clock.Now().Fmt(clock.FmtCompactDate)))
	file, err := fs.OpenOrCreate(logFileName, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	if err == nil {
		log.Out = file
	} else {
		Log.Error("打开日志文件失败", err)
	}
	// 每隔一天重新切换日志文件输出
	js.SetTimeout(log.Init, time.Duration(clock.GetDayLastMillisecond()))
}

// ClearOld 历史的日志文件需要删除
func (log *Logger) ClearOld() {
	fs.EachFilesRemove(cts.LogDir, func(filename string, fi os.FileInfo) bool {
		timeStrArray := regexp.MustCompile(`logrus\.(.+)\.log`).FindStringSubmatch(filename)
		if len(timeStrArray) == 0 {
			return true
		}
		return clock.MustParse(timeStrArray[1], clock.FmtCompactDate).SinceDays() > log.KeepDay
	})
}
