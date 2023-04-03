package logger

import (
	"log"
	"os"
)

type Logger struct {
	LogInfo  *log.Logger
	LogError *log.Logger
}

func NewLogger() *Logger {
	infoLog := log.New(os.Stdout, "\033[34m[INFO]\033[0m\t", log.Ldate|log.Ltime)
	errorLog := log.New(os.Stderr, "\033[31m[ERROR]\033[0m\t", log.Ldate|log.Ltime|log.Lshortfile)

	return &Logger{
		LogInfo:  infoLog,
		LogError: errorLog,
	}
}
