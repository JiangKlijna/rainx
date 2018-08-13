package main

import (
	"os"
	"fmt"
	"time"
	"bytes"
)

const (
	L_DEBUG = uint(iota)
	L_INFO
	L_WARNING
	L_ERROR
)

var (
	TAG_DEBUG   = []byte(" [D] ")
	TAG_INFO    = []byte(" [I] ")
	TAG_WARNING = []byte(" [W] ")
	TAG_ERROR   = []byte(" [E] ")
)

// Simple Logger Interface
type Logger interface {
	Debug(v ...interface{})
	Info(v ...interface{})
	Warning(v ...interface{})
	Error(v ...interface{})
	Close()
}

// Logger Implement
type loggerImpl struct {
	file  *os.File
	line  *os.File
	level uint
}

// New creates a new Logger
func NewLogger(filename string, level uint) (Logger, error) {
	f, err := os.OpenFile(filename, os.O_CREATE|os.O_WRONLY|os.O_APPEND, os.ModePerm)
	if err != nil {
		return nil, err
	}
	return &loggerImpl{
		file:  f,
		line:  os.Stdout,
		level: level,
	}, nil
}

// print Debug Log
func (l *loggerImpl) Debug(v ...interface{}) {
	l.print(L_DEBUG, TAG_DEBUG, v...)
}

// print Info Log
func (l *loggerImpl) Info(v ...interface{}) {
	l.print(L_INFO, TAG_INFO, v...)
}

// print Warning Log
func (l *loggerImpl) Warning(v ...interface{}) {
	l.print(L_WARNING, TAG_WARNING, v...)
}

// print Error Log
func (l *loggerImpl) Error(v ...interface{}) {
	l.print(L_ERROR, TAG_ERROR, v...)
}

// print Log
func (l *loggerImpl) print(level uint, tag []byte, v ...interface{}) {
	data := l.logData(tag, v...)
	l.line.Write(data)
	if l.level < level {
		l.file.Write(data)
	}
}

// customize print Log
func (l *loggerImpl) Print(v ...interface{}) {
	data := []byte(fmt.Sprintln(v...))
	l.line.Write(data)
	l.file.Write(data)
}

// return the log data
func (l *loggerImpl) logData(tag []byte, v ...interface{}) []byte {
	buf := bytes.Buffer{}
	buf.WriteString(time.Now().Format("2006-01-02 15:04:05"))
	buf.Write(tag)
	buf.WriteString(fmt.Sprintln(v...))
	return buf.Bytes()
}

// Close the File
func (l *loggerImpl) Close() {
	l.file.Close()
}
