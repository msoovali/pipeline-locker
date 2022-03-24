package logger

import (
	"log"
	"os"
)

const (
	infoPrefix  string = "INFO\t"
	errorPrefix string = "ERROR\t"
)

type Logger struct {
	Info  *log.Logger
	Error *log.Logger
}

func New() *Logger {
	infoLog := log.New(os.Stdout, infoPrefix, log.Ldate|log.Ltime)
	errorLog := log.New(os.Stderr, errorPrefix, log.Ldate|log.Ltime|log.Lshortfile)

	return &Logger{
		Info:  infoLog,
		Error: errorLog,
	}
}
