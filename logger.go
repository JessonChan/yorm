package yorm

import (
	"fmt"
)

type Logger interface {
	Debug(string, ...interface{})
	Warn(string, ...interface{})
	Error(string, ...interface{})
}

var log Logger = &yormLogger{}

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
}

func InitLogger(l Logger) {
	log = l
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
	fmt.Printf("[%s] %s\n", lv, msg)
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
