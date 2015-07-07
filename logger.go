package yorm

import "fmt"

type loggerInterface interface {
	Debug(string, ...interface{})
	Warn(string, ...interface{})
	Error(string, ...interface{})
}

var yogger loggerInterface = &yormLogger{}
var loggerLevel = ""

var (
	CloseLevel = ""
	DebugLevel = "Debug"
	WarnLevel  = "Warn"
	ErrorLevel = "Error"
)

var lmap = map[string]int{CloseLevel: 100, DebugLevel: 1, WarnLevel: 7, ErrorLevel: 14}

type yormLogger struct {
}

func InitLogger(l loggerInterface) {
	yogger = l
}

func SetLoggerLevel(s string) {
	switch s {
	case CloseLevel, DebugLevel, WarnLevel, ErrorLevel:
		loggerLevel = s
	default:
	}
}

func writeMsg(l string, s string, f ...interface{}) {
	if lmap[loggerLevel] > lmap[l] {
		return
	}
	msg := s
	if len(f) > 0 {
		msg = fmt.Sprintf(s, f...)
	}
	fmt.Printf("[%s] %s\n", l, msg)
}

func (yl *yormLogger) Debug(s string, f ...interface{}) {
	writeMsg(DebugLevel, s, f...)
}

func (yl *yormLogger) Warn(s string, f ...interface{}) {
	writeMsg(WarnLevel, s, f...)
}

func (yl *yormLogger) Error(s string, f ...interface{}) {
	writeMsg(ErrorLevel, s, f...)
}
