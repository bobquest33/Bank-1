package log

import (
	"log"
	"os"
)

type LogType struct {
	LogLog			*log.Logger		//用于错误的日志类型。
	LogAssert		*log.Logger		//用于断言(Assert)的日志类型（这些表明Unity自身的一个错误）。
	LogWar			*log.Logger		//用于警告(Warning)的日志类型。
	LogErr			*log.Logger		//用于普通日志消息的日志类型。
	LogException	*log.Logger		//用于异常的日志类型。
}

type LogFileInfo struct {
	logTxt	*os.File
	assTxt	*os.File
	warTxt	*os.File
	errTxt	*os.File
	excTxt	*os.File
}
