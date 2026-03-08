package log

import (
	"fmt"
	"io"
	"os"
	"runtime"
	"season-studio/go-utils/ioex"
	"sync"
	"time"
)

type LogLevel byte

const (
	LOG_ERROR LogLevel = 1 << iota
	LOG_INFO
	LOG_WARN
	LOG_DEBUG
)

var (
	logLevel              = LOG_INFO
	showFilInfo           = false
	chMsg                 = make(chan *[]byte, 1)
	loggerErr   io.Writer = os.Stderr
	loggerInfo  io.Writer = os.Stdout
	loggerWarn  io.Writer = os.Stderr
	loggerDebug io.Writer = os.Stdout
)

var bufPool = sync.Pool{
	New: func() any {
		ret := make([]byte, 1)
		return &ret
	},
}

func loggerWorkerProc() {
	for {
		buf := <-chMsg
		if buf == nil {
			ioex.Flush(loggerErr)
			ioex.Flush(loggerInfo)
			ioex.Flush(loggerDebug)
			ioex.Flush(loggerWarn)
			continue
		}
		var writer io.Writer
		switch LogLevel((*buf)[0]) {
		case LOG_ERROR:
			writer = loggerErr
		case LOG_INFO:
			writer = loggerInfo
		case LOG_WARN:
			writer = loggerWarn
		case LOG_DEBUG:
			writer = loggerDebug
		}
		writer.Write((*buf)[1:])
		bufPool.Put(buf)
	}
}

func formatFileInfo(buf []byte) []byte {
	_, file, line, ok := runtime.Caller(4)
	if !ok {
		return buf
	}
	return fmt.Appendf(buf, "%s:%d ", file, line)
}

func formatlogPrefix(logType LogLevel, buf []byte) []byte {
	buf[0] = byte(logType)
	now := time.Now()
	buf = fmt.Appendf(buf, "%04d/%02d/%02d %02d:%02d:%02d.%03d ", now.Year(), now.Month(), now.Day(), now.Hour(), now.Minute(), now.Second(), now.Nanosecond()/1e6)
	switch logType {
	case LOG_ERROR:
		buf = fmt.Appendf(buf, "\x1b[31m[ERROR] \x1b[0m")
	case LOG_INFO:
		buf = fmt.Appendf(buf, "\x1b[32m[INFO] \x1b[0m")
	case LOG_WARN:
		buf = fmt.Appendf(buf, "\x1b[35m[WARN] \x1b[0m")
	case LOG_DEBUG:
		buf = fmt.Appendf(buf, "\x1b[33m[DEBUG] \x1b[0m")
	}
	if showFilInfo {
		buf = formatFileInfo(buf)
	}
	return buf
}

func sendLogMsg(logType LogLevel, v []any) {
	buf := bufPool.Get().(*[]byte)
	*buf = (*buf)[:1]
	*buf = formatlogPrefix(logType, *buf)
	for i, vi := range v {
		if i > 0 {
			*buf = fmt.Append(*buf, " ")
		}
		*buf = fmt.Append(*buf, vi)
	}
	*buf = fmt.Appendln(*buf)
	chMsg <- buf
}

func sendLogMsgf(logType LogLevel, format string, v []any) {
	buf := bufPool.Get().(*[]byte)
	*buf = (*buf)[:1]
	*buf = formatlogPrefix(logType, *buf)
	*buf = fmt.Appendf(*buf, format, v...)
	*buf = fmt.Appendln(*buf)
	chMsg <- buf
}

type LogMsgHandler func(func(v ...any))

func logexOutput(handler LogMsgHandler) []any {
	buf := []any{}
	handler(func(v ...any) {
		buf = append(buf, v...)
	})
	return buf
}

func Error(msg ...any) {
	sendLogMsg(LOG_ERROR, msg)
}

func Errorf(format string, v ...any) {
	sendLogMsgf(LOG_ERROR, format, v)
}

func ErrorEx(handler LogMsgHandler) {
	sendLogMsg(LOG_ERROR, logexOutput(handler))
}

func Info(msg ...any) {
	if logLevel < LOG_INFO {
		return
	}
	sendLogMsg(LOG_INFO, msg)
}

func Infof(format string, v ...any) {
	if logLevel < LOG_INFO {
		return
	}
	sendLogMsgf(LOG_INFO, format, v)
}

func InfoEx(handler LogMsgHandler) {
	if logLevel < LOG_INFO {
		return
	}
	sendLogMsg(LOG_INFO, logexOutput(handler))
}

func Warn(msg ...any) {
	if logLevel < LOG_WARN {
		return
	}
	sendLogMsg(LOG_WARN, msg)
}

func Warnf(format string, v ...any) {
	if logLevel < LOG_WARN {
		return
	}
	sendLogMsgf(LOG_WARN, format, v)
}

func WarnEx(handler LogMsgHandler) {
	if logLevel < LOG_WARN {
		return
	}
	sendLogMsg(LOG_WARN, logexOutput(handler))
}

func Debug(msg ...any) {
	if logLevel < LOG_DEBUG {
		return
	}
	sendLogMsg(LOG_DEBUG, msg)
}

func Debugf(format string, v ...any) {
	if logLevel < LOG_DEBUG {
		return
	}
	sendLogMsgf(LOG_DEBUG, format, v)
}

func DebugEx(handler LogMsgHandler) {
	if logLevel < LOG_DEBUG {
		return
	}
	sendLogMsg(LOG_DEBUG, logexOutput(handler))
}

func SetOutput(logType LogLevel, writer io.Writer) error {
	if writer == nil {
		return fmt.Errorf("writer cannot be nil")
	}

	var target *io.Writer

	switch logType {
	case LOG_ERROR:
		target = &loggerErr
	case LOG_WARN:
		target = &loggerWarn
	case LOG_DEBUG:
		target = &loggerDebug
	default:
		target = &loggerInfo
	}

	*target = writer

	return nil
}

func SetOutputFile(logType LogLevel, filePath string) error {
	f, err := os.OpenFile(filePath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return err
	}
	SetOutput(logType, f)
	return nil
}

func ShowFileInfo(enable bool) {
	showFilInfo = enable
}

func SetLevel(logType LogLevel) {
	logLevel = logType
}

func GetLevel() LogLevel {
	return logLevel
}

func Flush() {
	chMsg <- nil
	chMsg <- nil
}

func init() {
	go loggerWorkerProc()
}
