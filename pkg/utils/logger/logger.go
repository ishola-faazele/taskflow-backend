package logger

import (
	"log"
	"os"
)

type StdLogger struct {
	info  *log.Logger
	error *log.Logger
	warn  *log.Logger
}

func NewStdLogger() *StdLogger {
	return &StdLogger{
		info:  log.New(os.Stdout, "INFO: ", log.LstdFlags),
		error: log.New(os.Stderr, "ERROR: ", log.LstdFlags),
		warn:  log.New(os.Stdout, "WARN: ", log.LstdFlags),
	}
}

func (l *StdLogger) Info(msg string) {
	l.info.Println(msg)
}

func (l *StdLogger) Error(msg string) {
	l.error.Println(msg)
}
func (l *StdLogger) Warn(msg string) {
	l.warn.Println(msg)
}
