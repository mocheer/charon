package logger

import (
	"fmt"
	"os"
	"regexp"
	"time"

	"github.com/mocheer/pluto/clock"
	"github.com/mocheer/pluto/fs"
	"github.com/mocheer/pluto/tm"
	"github.com/sirupsen/logrus"
)

type Logger struct {
	*logrus.Logger
}

var log *Logger

// Init 初始化日志
func Init() {
	if log == nil {
		log = &Logger{logrus.New()}
	}
	// 历史的日志文件需要删除，这里只删除一个月前的文件
	fs.EachDirToRemove("log", func(filename string) bool {
		timeStrArray := regexp.MustCompile(`logrus\.(.+)\.log`).FindStringSubmatch(filename)
		if len(timeStrArray) == 0 {
			return true
		}
		return clock.MustParse(timeStrArray[1], clock.FmtCompactDate).SinceDays() > 10
	})

	logFileName := fmt.Sprintf("log/logrus.%s.log", clock.Now().Fmt(clock.FmtCompactDate))
	file, err := fs.OpenOrCreate(logFileName, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	if err == nil {
		log.Out = file
	} else {
		log.Info("打开日志文件失败", err)
	}
	// 每隔一天重新定位日志文件
	tm.SetTimeout(Init, time.Duration(clock.GetDayLastMillisecond()))

}

// Trace 输出调试信息
func Trace(args ...interface{}) {
	log.Trace(args...)
}

// Debug 输出调试信息
func Debug(args ...interface{}) {
	log.Debug(args...)
}

// Info 输出调试信息
func Info(args ...interface{}) {
	log.Info(args...)
}

// Warn 输出调试信息
func Warn(args ...interface{}) {
	log.Warn(args...)
}

// Error 输出调试信息
func Error(args ...interface{}) {
	log.Error(args...)
}
