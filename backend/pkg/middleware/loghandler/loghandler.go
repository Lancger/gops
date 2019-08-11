package handler

import (
	"fmt"
	"os"
	"path"
	"time"

	"gops/backend/glo/comfunc"

	"github.com/gin-gonic/gin"
	rotatelogs "github.com/lestrrat-go/file-rotatelogs"
	"github.com/rifflock/lfshook"
	logutils "github.com/sirupsen/logrus"
)

// Logger 日志记录句柄
func Logger() gin.HandlerFunc {
	logClient := logutils.New()
	var logPath = "./logs" // 日志打印到指定的目录
	// 目录不存在则创建
	var logName = "access.log"
	if isExists, _ := comfunc.PathExists(logPath); isExists == false {
		os.MkdirAll(logPath, os.ModePerm)
	}
	fileName := path.Join(logPath, logName)
	//禁止logrus的输出
	src, err := os.OpenFile(fileName, os.O_APPEND|os.O_WRONLY, os.ModeAppend)
	if err != nil {
		fmt.Println("err", err)
	}

	//apiLogPath := "gin-api.log"
	logWriter, err := rotatelogs.New(
		fileName+".%Y-%m-%d.log",
		rotatelogs.WithLinkName(logName),          // 生成软链，指向最新日志文件
		rotatelogs.WithMaxAge(7*24*time.Hour),     // 文件最大保存时间
		rotatelogs.WithRotationTime(24*time.Hour), // 日志切割时间间隔
	)
	// 设置日志输出的路径
	logClient.Out = src
	logClient.SetLevel(logutils.DebugLevel)
	writeMap := lfshook.WriterMap{
		logutils.InfoLevel:  logWriter,
		logutils.FatalLevel: logWriter,
		logutils.DebugLevel: logWriter, // 为不同级别设置不同的输出目的
		logutils.WarnLevel:  logWriter,
		logutils.ErrorLevel: logWriter,
		logutils.PanicLevel: logWriter,
	}
	lfHook := lfshook.NewHook(writeMap, &logutils.JSONFormatter{})
	logClient.AddHook(lfHook)

	return func(c *gin.Context) {
		// type jsonStruct struct {
		// 	StatusCode int
		// 	Lantency   time.Duration
		// 	ClientIP   string
		// 	Method     string
		// 	URI        string
		// }
		// var (
		// 	message string
		// )
		// 开始时间
		start := time.Now()
		// 处理请求
		c.Next()
		// 结束时间
		end := time.Now()
		//执行时间
		// logMessage := jsonStruct{
		// 	Lantency:   end.Sub(start),
		// 	URI:        c.Request.URL.Path,
		// 	ClientIP:   c.ClientIP(),
		// 	Method:     c.Request.Method,
		// 	StatusCode: c.Writer.Status(),
		// }
		// jsonBytes, err := json.Marshal(&logMessage)
		// if err != nil {
		// 	message = fmt.Sprintf("err, %s", err.Error())
		// } else {
		// 	message = string(jsonBytes)
		// }
		// 这里是指定日志打印出来的格式。分别是状态码，执行时间,请求ip,请求方法,请求路由
		logClient.WithFields(logutils.Fields{
			"uri":    c.Request.URL.Path,
			"ip":     c.ClientIP(),
			"method": c.Request.Method,
			"code":   c.Writer.Status(),
			"ms":     end.Sub(start) / 1e6,
		}).Info("OK")
	}
}

func LogRusutil(v ...interface{}) {
	logClient2 := logutils.New()
	var logPath = "./logs" // 日志打印到指定的目录
	// 目录不存在则创建
	var logName = "debug.log"
	if isExists, _ := comfunc.PathExists(logPath); isExists == false {
		os.MkdirAll(logPath, os.ModePerm)
	}
	fileName := path.Join(logPath, logName)
	//禁止logrus的输出
	src, err := os.OpenFile(fileName, os.O_APPEND|os.O_WRONLY, os.ModeAppend)
	if err != nil {
		fmt.Println("err", err)
	}
	// 设置日志输出的路径
	logClient2.Out = src
	logClient2.SetLevel(logutils.DebugLevel)
	//apiLogPath := "gin-api.log"
	logWriter2, err := rotatelogs.New(
		fileName+".%Y-%m-%d.log",
		rotatelogs.WithLinkName(logName),          // 生成软链，指向最新日志文件
		rotatelogs.WithMaxAge(7*24*time.Hour),     // 文件最大保存时间
		rotatelogs.WithRotationTime(24*time.Hour), // 日志切割时间间隔
	)
	writeMap := lfshook.WriterMap{
		logutils.InfoLevel:  logWriter2,
		logutils.FatalLevel: logWriter2,
		logutils.DebugLevel: logWriter2, // 为不同级别设置不同的输出目的
		logutils.WarnLevel:  logWriter2,
		logutils.ErrorLevel: logWriter2,
		logutils.PanicLevel: logWriter2,
	}
	lfHook := lfshook.NewHook(writeMap, &logutils.JSONFormatter{})
	logClient2.AddHook(lfHook)
	// logClient2.WithFields(logutils.Fields{
	// 	"uri":    c.Request.URL.Path,
	// 	"ip":     c.ClientIP(),
	// 	"method": c.Request.Method,
	// 	"code":   c.Writer.Status(),
	// 	"ms":     end.Sub(start) / 1e6,
	// }).Info("OK")
	logClient2.Info(v...)
}
