package logger

import (
	"fmt"
	rotatelogs "github.com/lestrrat-go/file-rotatelogs"
	"github.com/rifflock/lfshook"
	"github.com/sirupsen/logrus"
	"go.gin.order/src/internal/pojo"
	"log"
	"os"
	"time"
)

var (
	Logger = logrus.New() // 初始化日志对象
)

func LogInit(logConf *pojo.LogCof) {
	//确定日志存放路径
	if _, err := os.Stat(logConf.LogPath); os.IsNotExist(err) {
		if err = os.MkdirAll(logConf.LogPath, 0755); err != nil {
			panic(fmt.Sprintf("日志文件存放地址创建错误：%s", err))
			//fmt.Println("文件创建错误：", err)
			//return
		}
	}
	src, err := os.OpenFile(logConf.LogPath+"/log", os.O_RDWR|os.O_CREATE, 0644) // 初始化日志文件对象
	if err != nil {
		panic(fmt.Sprintf("日志文件创建错误：%s", err))
		//fmt.Println("err: ", err)
		//return
	}
	Logger.Out = src // 把产生的日志内容写进日志文件中
	// 设置日志级别
	level, err := logrus.ParseLevel(logConf.LogLevel)
	if err != nil {
		panic(fmt.Sprintf("日志级别解析错误：%s", err))
	}
	Logger.SetLevel(level)
	// 日志分隔：1. 每天产生的日志写在不同的文件；2. 只保留一定时间的日志（例如：一星期）
	logWriter, _ := rotatelogs.New(
		logConf.LogPath+"/%Y%m%d.log",             // 日志文件名格式
		rotatelogs.WithMaxAge(7*24*time.Hour),     // 最多保留7天之内的日志
		rotatelogs.WithRotationTime(24*time.Hour), // 一天保存一个日志文件
		//rotatelogs.WithLinkName(logConf.LinkName), // 为最新日志建立软连接
	)
	infoFormatter := &logrus.TextFormatter{
		TimestampFormat: "2006-01-02 15:04:05",
		FullTimestamp:   true,
	}
	errorFormatter := &logrus.JSONFormatter{
		TimestampFormat: "2006-01-02 15:04:05",
	}
	writeMap := lfshook.WriterMap{
		logrus.InfoLevel:  logWriter, // info级别使用logWriter写日志
		logrus.WarnLevel:  logWriter,
		logrus.FatalLevel: logWriter,
		logrus.DebugLevel: logWriter,
		logrus.ErrorLevel: logWriter,
		logrus.PanicLevel: logWriter,
	}
	infoHook := lfshook.NewHook(writeMap, logrus.Formatter(infoFormatter))
	warnHook := lfshook.NewHook(writeMap, logrus.Formatter(errorFormatter))
	Logger.AddHook(infoHook) //info级别数据
	Logger.AddHook(warnHook) //warn级别日志
	log.Printf("日志初始化成功\n")
}
