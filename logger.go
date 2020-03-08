package yorm

import (
	"fmt"
	"log"
	"os"
)

var (
	Close = ""
	Debug = "Debug"
	Warn  = "Warn"
	Error = "Error"

	level = Close
)

var lmap = map[string]int{
	Close: 100,
	Debug: 1,
	Warn:  7,
	Error: 14}

type yormLogger struct {
	innerLogger *log.Logger
}

var logger = &yormLogger{innerLogger: log.New(os.Stderr, "", log.Lshortfile|log.LstdFlags)}

func InitLogger(l *log.Logger) {
	logger.innerLogger = l
}

func SetLoggerLevel(s string) {
	switch s {
	case Close, Debug, Warn, Error:
		level = s

	default:
	}
}

func writeMsg(lv string, s string, f ...interface{}) {
	if lmap[level] > lmap[lv] {
		return
	}
	msg := s
	if len(f) > 0 {
		msg = fmt.Sprintf(s, f...)
	}
	logger.innerLogger.Printf("[%s] %s\n", lv, msg)
}

func (y *yormLogger) Debug(s string, f ...interface{}) {
	writeMsg(Debug, s, f...)
}

func (y *yormLogger) Warn(s string, f ...interface{}) {
	writeMsg(Warn, s, f...)
}

func (y *yormLogger) Error(s string, f ...interface{}) {
	writeMsg(Error, s, f...)
}
