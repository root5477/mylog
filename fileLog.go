package mylog

import (
	"fmt"
	"os"
	"path"
	"time"
)

type FileLogger struct {
	Level        LogLevel
	FilePath     string
	FileName     string
	MaxFileSize  int64
	FileObj      *os.File
	ErrorFileObj *os.File
}

func NewFileLogger(level, filePath, fileName string, size int64) *FileLogger {
	logLevel, err := levelToLogLevel(level)
	if err != nil {
		panic(err)
	}
	fl := &FileLogger{
		Level:       logLevel,
		FilePath:    filePath,
		FileName:    fileName,
		MaxFileSize: size,
	}
	err = fl.initFile()
	if err != nil {
		panic(err)
	}
	return fl
}

//初始化好打开的日志文件
func (f1 *FileLogger) initFile() error {
	FullFileName := path.Join(f1.FilePath, f1.FileName)
	fileObj, err := os.OpenFile(FullFileName, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0644)
	if err != nil {
		fmt.Printf(" os.OpenFile %s失败, the err is:%v", FullFileName, err)
		return err
	}
	errFileObj, err2 := os.OpenFile(FullFileName+".err", os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0644)
	if err2 != nil {
		fmt.Printf(" os.OpenFile %s失败, the err is:%v", FullFileName+".err", err2)
		return err2
	}
	//到这里说明日志文件都已经打开
	f1.FileObj = fileObj
	f1.ErrorFileObj = errFileObj
	return nil
}

//关闭日志文件
func (f *FileLogger) Close() {
	f.FileObj.Close()
	f.ErrorFileObj.Close()
}

//每次写入日志前，调用此方法判断日志文件大小是否需要切割,true表示需要切割，false表示不需要切割
func (logFile *FileLogger) checkLogFileSize(file *os.File) bool {
	fileInfo, err := file.Stat()
	if err != nil {
		fmt.Println("获取日志文件大小失败, err is:", err)
		return false
	}
	if fileInfo.Size() < logFile.MaxFileSize {
		return false
	} else {
		return true
	}
}

func (f *FileLogger) cutLogFileOnSize(file *os.File) (*os.File, error) {

	nowStr := time.Now().Format("20060102_1504_05000")
	fileInfo, err := file.Stat()
	if err != nil {
		fmt.Printf("获取当前日志文件名失败, err is:%v", err)
		return nil, err
	}
	fileName := fileInfo.Name()
	logNeedBakPath := path.Join(f.FilePath, fileName)
	logBakPath := fmt.Sprintf("%s.bak&s", logNeedBakPath, nowStr)

	// 1--关闭当前日志文件
	file.Close()
	// 2--rename原来的日志文件

	errBak := os.Rename(logNeedBakPath, logBakPath)
	if errBak != nil {
		fmt.Printf("备份文件失败了，err is:%v", errBak)
		return nil, errBak
	}
	// 3--创建新的日志文件并打开
	newFile, errNewFile := os.OpenFile(logNeedBakPath, os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0644)
	if errNewFile != nil {
		fmt.Printf("创建新日志文件%s失败了，err is:%v", logNeedBakPath, errNewFile)
		return nil, errNewFile
	}
	// 4--将新打开的日志文件赋值给 file.FileObj
	return newFile, nil
}

func (file *FileLogger) logFile(lv LogLevel, formatString string, other ...interface{}) {
	fileName, funcName, lineNum := getInfo(3)
	message := fmt.Sprintf(formatString, other...)
	now := time.Now().Format("2006-01-02 15:04:05")
	level := getLevelString(lv)
	//检查日志文件大小
	cutTag := file.checkLogFileSize(file.FileObj)
	if cutTag == true {
		//需要切割
		newFile, errCut := file.cutLogFileOnSize(file.FileObj)
		if errCut != nil {
			fmt.Printf("切割日志文件出错，err is:%v", errCut)
			return
		}
		file.FileObj = newFile
	}
	fmt.Fprintf(file.FileObj, "[%s][%s][%v  %v  %v行]%s\n", now, level, fileName, funcName, lineNum, message)
	//如果日志级别比Error还大，应再记录一份error日志
	if lv >= ERROR {
		cutTag := file.checkLogFileSize(file.ErrorFileObj)
		if cutTag == true {
			//需要切割
			newFile, errCut := file.cutLogFileOnSize(file.ErrorFileObj)
			if errCut != nil {
				fmt.Printf("切割日志文件出错，err is:%v", errCut)
				return
			}
			file.ErrorFileObj = newFile
		}
		fmt.Fprintf(file.ErrorFileObj, "[%s][%s][%v  %v  %v行]%s\n", now, level, fileName, funcName, lineNum, message)
	}
}

func (fileLogger *FileLogger) Debug(formatStr string, other ...interface{}) {
	//timeNow := time.Now().Format("2006/1/2")
	//logName := fileLogger.FilePath + "/debug" + strings.Split(timeNow, "/")[1] + "_"+ strings.Split(timeNow, "/")[2] + ".log"
	fileLogger.logFile(DEBUG, formatStr, other...)
}

func (fileLogger *FileLogger) Info(formatStr string, other ...interface{}) {
	fileLogger.logFile(INFO, formatStr, other...)
}

func (fileLogger *FileLogger) Warn(formatStr string, other ...interface{}) {
	fileLogger.logFile(WARN, formatStr, other...)
}

func (fileLogger *FileLogger) Error(formatStr string, other ...interface{}) {
	fileLogger.logFile(ERROR, formatStr, other...)
}

func (fileLogger *FileLogger) Fatal(formatStr string, other ...interface{}) {
	fileLogger.logFile(FATAL, formatStr, other...)
}
