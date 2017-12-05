package log

import (
	"fmt"
	"log"
	"os"
	"runtime"
	_ "time"
	"path/filepath"
	"strings"
	"github.com/fatih/color"
)

type defaultLogger struct {
	*log.Logger
}

func GetCallerInfo(calldepth int) (fn string, file string, line int, ok bool) {
	var pc uintptr
	pc, file, line, ok = runtime.Caller(calldepth)
	if ok {
		if pfn := runtime.FuncForPC(pc); pfn != nil {
			fn = pfn.Name()
		}else {
			fn = "?()"
		}
	}else {
		fn = "?()"
		file = "?"
		line = 0
	}
	return
}

func GetCallerFuncname(calldepth int, shortpath bool) (fn string) {
	fn, _, _, _ = GetCallerInfo(calldepth)
	if shortpath {
		dotName := filepath.Ext(fn)
		fn = strings.TrimLeft(dotName, ".") + "()"
	}
	return fn
}

func (l *defaultLogger)Trace(v ...interface{}) {
	l.Output(calldepth, header("TRACE", fmt.Sprint(v...)))
}

func (l *defaultLogger)Tracef(format string, v ...interface{}) {
	l.Output(calldepth, header("TRACE", fmt.Sprintf(format, v...)))
}

func (l *defaultLogger) Debug(v ...interface{}) {
	l.Output(calldepth, header(color.CyanString("DEBUG"), fmt.Sprint(v...)))
}

func (l *defaultLogger) Debugf(format string, v ...interface{}) {
	l.Output(calldepth, header(color.CyanString("DEBUG"), fmt.Sprintf(format, v...)))
}

func (l *defaultLogger) Info(v ...interface{}) {
	l.Output(calldepth, header(color.GreenString("INFO "), fmt.Sprint(v...)))
}

func (l *defaultLogger) Infof(format string, v ...interface{}) {
	l.Output(calldepth, header(color.GreenString("INFO "), fmt.Sprintf(format, v...)))
}

func (l *defaultLogger) Notice(v ...interface{}) {
	l.Output(calldepth, header(color.HiBlueString("NOTICE "), fmt.Sprint(v...)))
}

func (l *defaultLogger) Noticef(format string, v ...interface{}) {
	l.Output(calldepth, header(color.HiBlueString("NOTICE "), fmt.Sprintf(format, v...)))
}


func (l *defaultLogger) Warn(v ...interface{}) {
	l.Output(calldepth, header(color.YellowString("WARN "), fmt.Sprint(v...)))
}

func (l *defaultLogger) Warnf(format string, v ...interface{}) {
	l.Output(calldepth, header(color.YellowString("WARN "), fmt.Sprintf(format, v...)))
}

func (l *defaultLogger) Error(v ...interface{}) {
	l.Output(calldepth, header(color.RedString("ERROR"), fmt.Sprint(v...)))
}

func (l *defaultLogger) Errorf(format string, v ...interface{}) {
	l.Output(calldepth, header(color.RedString("ERROR"), fmt.Sprintf(format, v...)))
}

func (l *defaultLogger) Fatal(v ...interface{}) {
	l.Output(calldepth, header(color.MagentaString("FATAL"), fmt.Sprint(v...)))
	os.Exit(1)
}

func (l *defaultLogger) Fatalf(format string, v ...interface{}) {
	l.Output(calldepth, header(color.MagentaString("FATAL"), fmt.Sprintf(format, v...)))
	os.Exit(1)
}

func (l *defaultLogger) Panic(v ...interface{}) {
	l.Logger.Panic(v)
}

func (l *defaultLogger) Panicf(format string, v ...interface{}) {
	l.Logger.Panicf(format, v...)
}

func header(lvl, msg string) string {
	return fmt.Sprintf("%s %s: %s", lvl, GetCallerFuncname(calldepth + 2, true), msg)
}

func headerOrig(lvl, msg string) string {
	return fmt.Sprintf("%s: %s", lvl, msg)
}
