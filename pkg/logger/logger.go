package logger

import (
	"log"
	"os"
)

type Logger struct {
	debug bool
}

func NewLogger(debug bool) *Logger {
	return &Logger{debug: debug}
}

func (l *Logger) Info(v ...interface{}) {
	log.Println("[INFO]", v)
}

func (l *Logger) Error(v ...interface{}) {
	log.Println("[ERROR]", v)
}

func (l *Logger) Fatal(v ...interface{}) {
	log.Println("[FATAL]", v)
	os.Exit(1) 
}
