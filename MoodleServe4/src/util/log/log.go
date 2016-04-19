package log

import (
	"fmt"
	"os"
	"log"
)

var fileStr LogFileInfo
var logType	LogType

func Init() {
	fileStr = LogFileInfo{}
	logType = LogType{}

	err := defineLogInfo()
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println("Log Init")
}

func createLogFile(filename string) (*os.File, error) {
	if _, err := os.Stat(filename); err != nil {
		file, err := os.Create(filename)
		if err != nil {
			return nil, err
		}
		return file, nil
	} else {
		file, err := os.OpenFile(filename, os.O_RDWR|os.O_APPEND, os.ModePerm)
		if err != nil {
			return nil, err
		}
		return file, nil
	}
}

func defineLogType() error {
	path := "/Users/tangs/IdeaProjects/MoodleServe4/src/util/log/"
	filename := path + "logTxt.txt"
	fPtr, err :=  createLogFile(filename)
	if err != nil {
		return err
		//失败了怎么办,应该不允许失败这里
	}
	fileStr.logTxt = fPtr

	filename = path + "assTxt.txt"
	fPtr, err =  createLogFile(filename)
	if err != nil {
		return err
	}
	fileStr.assTxt = fPtr

	filename = path + "warTxt.txt"
	fPtr, err =  createLogFile(filename)
	if err != nil {
		return err
	}
	fileStr.warTxt = fPtr

	filename = path + "errTxt.txt"
	fPtr, err =  createLogFile(filename)
	if err != nil {
		return err
	}
	fileStr.errTxt = fPtr

	filename = path + "excTxt.txt"
	fPtr, err =  createLogFile(filename)
	if err != nil {
		return err
	}
	fileStr.excTxt = fPtr
	return nil
}

func defineLogInfo() (error) {
	err := defineLogType()
	if err != nil {
		return err
		//how to solve the question if create or open file failed.
	}
	logType.LogLog			= log.New(fileStr.logTxt,		"[Log] : ",			log.Ldate|log.Ltime|log.Llongfile)
	logType.LogAssert		= log.New(fileStr.assTxt,		"[Assert] : ",		log.Ldate|log.Ltime|log.Llongfile)
	logType.LogWar			= log.New(fileStr.warTxt,		"[Waring] : ",		log.Ldate|log.Ltime|log.Llongfile)
	logType.LogErr			= log.New(fileStr.errTxt,		"[Error] : ",		log.Ldate|log.Ltime|log.Llongfile)
	logType.LogException	= log.New(fileStr.excTxt,		"[Exception] : ",	log.Ldate|log.Ltime|log.Llongfile)
	return nil
}


func AddLog (v ...interface{})  {
	logType.LogLog.Output(2, fmt.Sprintf("%s", v...))
}

func AddAssert (v ...interface{}) {
	logType.LogAssert.Output(2, fmt.Sprintf("%s", v...))
}

func AddWarning (v ...interface{}) {
	logType.LogWar.Output(2, fmt.Sprintf("%s", v...))
}

func AddError(v ...interface{}) {
	logType.LogErr.Output(2, fmt.Sprintf("%s", v...))
}

func AddException(v ...interface{}) {
	logType.LogException.Output(2, fmt.Sprintf("%s", v...))
}

func AddTrance(v ...interface{})  {
	// do something
}