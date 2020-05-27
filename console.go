package mylog

import (
	"errors"
	"fmt"
	"path"
	"runtime"
	"strings"
	"time"
)

type Clogger interface {
	Debug(formatString string, other ...interface{})
	Info(formatString string, other ...interface{})
	Warn(formatString string, other ...interface{})
	Error(formatString string, other ...interface{})
	Fatal(formatString string, other ...interface{})
}

//声明日志类型，可以选择c或f, 分别对应往终端打印的ConsoleLogger和往文件写入的FileLogger
func NewLogger(typeOfLogStr string, size int64, strs ...string) Clogger {
	switch typeOfLogStr {
	case "c":
		logger := NewConsoleloger(strs[0])
		return logger
	case "f":
		logger := NewFileLogger(strs[0], strs[1], strs[2], size)
		return logger
	default:
		err := errors.New("错误的类型！")
		panic(err)
	}
}

type LogLevel uint16

//定义日志级别
const (
	UNKONWN LogLevel = iota
	DEBUG
	INFO
	WARN
	ERROR
	FATAL
)

type ConsoleLogger struct {
	//LogPath string
	Level LogLevel
}

func levelToLogLevel(s string) (LogLevel, error) {
	sLower := strings.ToLower(s)
	switch sLower {
	case "debug":
		return DEBUG, nil
	case "info":
		return INFO, nil
	case "warn":
		return WARN, nil
	case "error":
		return ERROR, nil
	case "fatal":
		return FATAL, nil
	default:
		err := errors.New("未知的日志类型")
		return UNKONWN, err
	}
}

func getLevelString(lv LogLevel) string {
	switch lv {
	case DEBUG:
		return "DEBUG"
	case INFO:
		return "INFO"
	case WARN:
		return "WARN"
	case ERROR:
		return "ERROR"
	case FATAL:
		return "FATAL"
	default:
		return "DEBUG"
	}
}

func NewConsoleloger(level string) ConsoleLogger {
	logLevel, err := levelToLogLevel(level)
	if err != nil {
		panic(err)
	}
	return ConsoleLogger{
		//LogPath:path,
		Level: logLevel,
	}
}

func getInfo(skip int) (fileName, funcName string, lineNum int) {
	pc, file, line, ok := runtime.Caller(skip)
	if !ok {
		fmt.Printf("runtime.Caller failed!\n")
		return
	}
	funcName = runtime.FuncForPC(pc).Name()
	funcName = strings.Split(funcName, ".")[1]
	fileName = path.Base(file)

	return fileName, funcName, line
}

func (l ConsoleLogger) enable(level LogLevel) bool {
	return l.Level <= level
}

func log(lv LogLevel, formatString string, other ...interface{}) {
	fileName, funcName, lineNum := getInfo(3)
	message := fmt.Sprintf(formatString, other...)
	now := time.Now().Format("2006-01-02 15:04:05")
	level := getLevelString(lv)
	fmt.Printf("[%s][%s][%v  %v  %v行]%s\n", now, level, fileName, funcName, lineNum, message)
}

func (ConsoleLogger ConsoleLogger) Debug(formatString string, other ...interface{}) {
	if ConsoleLogger.enable(DEBUG) {
		log(DEBUG, formatString, other...)
	}
}

func (ConsoleLogger ConsoleLogger) Info(formatString string, other ...interface{}) {
	if ConsoleLogger.enable(INFO) {
		log(INFO, formatString, other...)
	}
}

func (ConsoleLogger ConsoleLogger) Warn(formatString string, other ...interface{}) {
	if ConsoleLogger.enable(WARN) {
		log(ConsoleLogger.Level, formatString, other...)
	}
}

func (ConsoleLogger ConsoleLogger) Error(formatString string, other ...interface{}) {
	if ConsoleLogger.enable(ERROR) {
		log(ConsoleLogger.Level, formatString, other...)
	}
}

func (ConsoleLogger ConsoleLogger) Fatal(formatString string, other ...interface{}) {
	if ConsoleLogger.enable(FATAL) {
		log(ConsoleLogger.Level, formatString, other...)
	}
}
