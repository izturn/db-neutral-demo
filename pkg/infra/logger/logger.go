package logger

import (
	"log"
)

type Logger struct{}

func (logger Logger) Log(args ...interface{}) {
	log.Println(args...)
}

func New() *Logger {
	return &Logger{}
}
