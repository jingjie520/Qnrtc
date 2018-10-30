package logutil

import (
	"core/utils/config"
	"fmt"
	"net/http"
	"os"
	"path"
	"time"

	rotatelogs "github.com/lestrrat-go/file-rotatelogs"
	"github.com/pkg/errors"
	"github.com/rifflock/lfshook"
	logrus "github.com/sirupsen/logrus"
)

// Log Logger日志对象
var Log *logrus.Logger

/*
NewLogger Create Logger
*/
func NewLogger() *logrus.Logger {
	if Log != nil {
		return Log
	}
	Log = logrus.New()

	Conf := config.Conf
	if Conf == nil {
		fmt.Println("加载配置文件conf.ini失败。请检查当前目录下是否存在该文件。")
		os.Exit(-1)
	}

	logPath := Conf.LogPath
	logName := Conf.LogName

	ConfigLocalFilesystemLogger(logPath, logName, time.Hour*24*365, time.Hour*24)
	return Log
}

func init() {
	Log = NewLogger()
}

/*
Error 记录错误日志
*/
func Error(str string, err error) {
	Log.Errorln(str, err)
}

//RequestInfo 记录访问信息
func RequestInfo(request *http.Request) {
	Log.WithFields(logrus.Fields{
		"IP":  request.RemoteAddr,
		"URL": request.URL,
	}).Info("Request")
}

/*
ConfigLocalFilesystemLogger 本地日志配置
*/
func ConfigLocalFilesystemLogger(logPath string, logFileName string, maxAge time.Duration, rotationTime time.Duration) {
	baseLogPaht := path.Join(logPath, logFileName)
	writer, err := rotatelogs.New(
		baseLogPaht+".%Y%m%d%H%M.log",
		rotatelogs.WithLinkName(baseLogPaht),      // 生成软链，指向最新日志文件
		rotatelogs.WithMaxAge(maxAge),             // 文件最大保存时间
		rotatelogs.WithRotationTime(rotationTime), // 日志切割时间间隔
	)
	if err != nil {
		Log.Errorf("config local file system logger error. %+v", errors.WithStack(err))
	}

	lfHook := lfshook.NewHook(lfshook.WriterMap{
		logrus.DebugLevel: writer, // 为不同级别设置不同的输出目的
		logrus.InfoLevel:  writer,
		logrus.WarnLevel:  writer,
		logrus.ErrorLevel: writer,
		logrus.FatalLevel: writer,
		logrus.PanicLevel: writer,
	}, &logrus.JSONFormatter{})
	Log.AddHook(lfHook)
}
