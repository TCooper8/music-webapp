package main

import (
	"fmt"
	"log"
	"time"
)

const (
	LOG_FATAL = 60
	LOG_ERROR = 50
	LOG_WARN  = 40
	LOG_INFO  = 30
	LOG_DEBUG = 20
	LOG_TRACE = 10
)

type Log struct {
	name  string
	level int
}

func NewLogger(name string, level int) *Log {
	logger := Log{
		name,
		level,
	}

	return &logger
}

func (logger *Log) post(level int, format string, args ...interface{}) {
	if level >= logger.level {
		log.Printf(
			"%s\t | %s\t | %d\t | %s\t",
			time.Now().Format(time.RFC822),
			logger.name,
			level,
			fmt.Sprintf(format, args...),
		)
	}
}

func (logger *Log) Fatal(format string, args ...interface{}) {
	logger.post(LOG_FATAL, format, args...)
}

func (logger *Log) Error(format string, args ...interface{}) {
	logger.post(LOG_ERROR, format, args...)
}

func (logger *Log) Warn(format string, args ...interface{}) {
	logger.post(LOG_WARN, format, args...)
}

func (logger *Log) Info(format string, args ...interface{}) {
	logger.post(LOG_INFO, format, args...)
}

func (logger *Log) Debug(format string, args ...interface{}) {
	logger.post(LOG_DEBUG, format, args...)
}

func (logger *Log) Trace(format string, args ...interface{}) {
	logger.post(LOG_TRACE, format, args...)
}
