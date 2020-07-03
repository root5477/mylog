# mylog
一个提供终端打印输出和文件输出，并且拥有fatal/error/warn/info/debug分级的日志库

#使用方法：
var Mylogger clogger.Clogger
#第一个参数可选"f"或"c",分别是往文件和往控制台输出日志,第二个参数是日志自动分割的文件大小，第三个参数是记录日志的最低级别，第四个参数是日志目录，第五个参数是日志名称
Mylogger = clogger.NewLogger("f", 102400, "debug", `E:outLog\`, "test.log")

Mylogger.Debug("this is a test")

Mylogger.Error("this is a test")
