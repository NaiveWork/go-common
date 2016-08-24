package logger

import (
	"fmt"
	"os"
	"path"
	"runtime"
	"time"
)

const timeFormat = "2006-01-02 15:04:05"

// LogLevel is one of DEBUG/NOTICE/WARNING/FATAL
type LogLevel int

const (
	logDebug LogLevel = iota
	logNotice
	logWarning
	logFatal
)

// Logger is the struct of logger
type Logger struct {
	Level  LogLevel
	Outter *os.File
}

func newLogger() *Logger {
	return &Logger{
		Level:  logDebug,
		Outter: os.Stdout,
	}
}

var (
	log    = newLogger()
	logStr = map[LogLevel]string{
		logDebug:   "DEBUG",
		logNotice:  "NOTICE",
		logWarning: "WARNING",
		logFatal:   "FATAL",
	}
)

// Debug to print debug level log
func Debug(format string, args ...interface{}) {
	log.debug(format, args...)
}

// Notice to print notice level log
func Notice(format string, args ...interface{}) {
	log.notice(format, args...)
}

// Warning to print warning level log
func Warning(format string, args ...interface{}) {
	log.warning(format, args...)
}

// Fatal to print fatal level log
func Fatal(format string, args ...interface{}) {
	log.fatal(format, args...)
}

// SetLogLevel to set the lowest log level
func SetLogLevel(lvl LogLevel) {
	log.setLogLevel(lvl)
}

// OpenLog to open the log file
func OpenLog(file string) (err error) {
	writer, err := os.OpenFile(file, os.O_WRONLY|os.O_APPEND|os.O_CREATE, os.FileMode(0666))
	if err != nil {
		return
	}
	log.setWriter(writer)
	return
}

// Close to close the opened log file
func Close() {
	log.close()
}

func (logger *Logger) setLogLevel(level LogLevel) {
	logger.Level = level
}

func (logger *Logger) setWriter(writer *os.File) {
	logger.Outter = writer
}

func (logger *Logger) close() {
	logger.Outter.Close()
}

func (logger *Logger) log(level LogLevel, format string, args ...interface{}) {
	if level < logger.Level {
		return
	}
	_, file, lineno, ok := runtime.Caller(3)
	if !ok {
		file = "unkown"
		lineno = 0
	}
	_, file = path.Split(file)
	ft := fmt.Sprintf("%s: %s: %d: %s:%d:", logStr[level],
		time.Now().Format(timeFormat),
		time.Now().Unix(), file, lineno)
	st := fmt.Sprintf(format, args...)
	fmt.Fprintf(logger.Outter, "%s %s\n", ft, st)
}

func (logger *Logger) debug(format string, args ...interface{}) {
	logger.log(logDebug, format, args...)
}

func (logger *Logger) notice(format string, args ...interface{}) {
	logger.log(logNotice, format, args...)
}

func (logger *Logger) warning(format string, args ...interface{}) {
	logger.log(logWarning, format, args...)
}

func (logger *Logger) fatal(format string, args ...interface{}) {
	logger.log(logFatal, format, args...)
}

// Printf 增加这个方法为了实现hbase/zk的logger
func (logger *Logger) Printf(format string, args ...interface{}) {
	logger.log(logDebug, format, args...)
}
