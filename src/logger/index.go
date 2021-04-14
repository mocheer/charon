package logger

import (
	"fmt"
	"os"
	"time"

	"github.com/mocheer/pluto/clock"
	"github.com/mocheer/pluto/fs"
	"github.com/mocheer/pluto/tm"
	"github.com/sirupsen/logrus"
)

var log = logrus.New()

// Init 初始化日志
func Init() {
	logFileName := fmt.Sprintf("log/logrus.%s.log", clock.New().Format(clock.SimpleDateFormat))
	file, err := fs.OpenOrCreate(logFileName, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	if err == nil {
		log.Out = file
	} else {
		log.Info("输出日志到文件失败", err)
	}
	//
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
