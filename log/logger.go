package log

import (
	"log"
	"os"
	//"flag"
)

const (
	calldepth = 3
)

type LogLevel int

//
// 0  	TRACE  最详细一级的日志，打印系统内部，库内部的调试日志信息
// 1  	DEBUG  普通调试日志，较详细， 用于业务调试，测试流程等
//
// 2 	INFO   正常日志， 一般用于打印表示程序运行状态、逻辑的提醒性日志
// 3	NOTICE 统计日志， 用于业务，指标数据输出，统计、分析的数据日志
// 4	WARN   非正常日志, 业务警号，但不会影响正常运行的错误或者警告
//
// 5 	ERROR  非正常日志，系统出现严重错误，但是不影响其他业务处理等。
// 6 	FATAL  系统致命错误， 可能导致行为可能不正常, 程序主动退出
// 7 	Panic  系统异常，必须终止程序运行
//
const (
	LOG_LEVEL_TRACE  LogLevel = iota
	LOG_LEVEL_DEBUG
	LOG_LEVEL_INFO
	LOG_LEVEL_NOTICE
	LOG_LEVEL_WARN
	LOG_LEVEL_ERROR
	LOG_LEVEL_FATAL
	LOG_LEVEL_PANIC
)

var loglevel LogLevel = LOG_LEVEL_INFO
var l Logger = &defaultLogger{log.New(os.Stderr, "", log.LstdFlags|log.Lshortfile)}

func init() {

	/*loglevel = LogLevel(*flag.Int("loglevel", int(LOG_LEVEL_INFO),
		"set loglevel: 0-trace, 1-debug, 2-info, 3-notice, 4-wran, 5-error, 6-fatal, 7-panic")	)
	*/
}

type Logger interface {
	// more detailed info output for tracing
	Trace(v ...interface{})
	Tracef(format string, v ...interface{})

	// detailed info debug usage
	Debug(v ...interface{})
	Debugf(format string, v ...interface{})

	// normal output
	Info(v ...interface{})
	Infof(format string, v ...interface{})

	Notice(v ...interface{})
	Noticef(format string, v ...interface{})

	Warn(v ...interface{})
	Warnf(format string, v ...interface{})

	Error(v ...interface{})
	Errorf(format string, v ...interface{})

	Fatal(v ...interface{})
	Fatalf(format string, v ...interface{})

	Panic(v ...interface{})
	Panicf(format string, v ...interface{})
}

func SetLogger(logger Logger) {
	l = logger
}

func SetLogLevel(level  LogLevel) {
	loglevel = level
}

func GetLogLevel() LogLevel {
	return  loglevel
}

func Trace(v ...interface{}) {
	if loglevel > LOG_LEVEL_TRACE {
		return
	}
	l.Trace(v)
}

func Tracef(format string, v ...interface{}) {
	if loglevel > LOG_LEVEL_TRACE {
		return
	}
	l.Tracef(format, v...)
}

func Debug(v ...interface{}) {
	if loglevel > LOG_LEVEL_DEBUG {
		return
	}
	l.Debug(v)
}
func Debugf(format string, v ...interface{}) {
	if loglevel > LOG_LEVEL_DEBUG {
		return
	}
	l.Debugf(format, v...)
}

func Info(v ...interface{}) {
	if loglevel > LOG_LEVEL_INFO {
		return
	}
	l.Info(v)
}
func Infof(format string, v ...interface{}) {
	if loglevel > LOG_LEVEL_INFO {
		return
	}
	l.Infof(format, v...)
}

func Notice(v ...interface{}) {
	if loglevel < LOG_LEVEL_NOTICE {
		return
	}
	l.Notice(v)
}
func Noticef(format string, v ...interface{}) {
	if loglevel < LOG_LEVEL_NOTICE {
		return
	}
	l.Noticef(format, v...)
}

func Warn(v ...interface{}) {
	if loglevel < LOG_LEVEL_WARN {
		return
	}
	l.Warn(v)
}
func Warnf(format string, v ...interface{}) {
	if loglevel < LOG_LEVEL_WARN {
		return
	}
	l.Warnf(format, v...)
}

func Error(v ...interface{}) {
	if loglevel < LOG_LEVEL_ERROR {
		return
	}
	l.Error(v)
}
func Errorf(format string, v ...interface{}) {
	if loglevel < LOG_LEVEL_ERROR {
		return
	}
	l.Errorf(format, v...)
}

func Fatal(v ...interface{}) {
	if loglevel < LOG_LEVEL_FATAL {
		return
	}
	l.Fatal(v)
}
func Fatalf(format string, v ...interface{}) {
	if loglevel < LOG_LEVEL_FATAL {
		return
	}
	l.Fatalf(format, v...)
}

func Panic(v ...interface{}) {
	if loglevel < LOG_LEVEL_PANIC {
		return
	}
	l.Panic(v)
}
func Panicf(format string, v ...interface{}) {
	if loglevel < LOG_LEVEL_PANIC {
		return
	}
	l.Panicf(format, v...)
}
