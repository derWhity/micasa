// Package log provides a simple wrapper to ease the use of GoKit's logger
package log

import (
	"github.com/go-kit/kit/log"
	"github.com/go-kit/kit/log/levels"
)

const (
	// LvlDebug is the debug level
	LvlDebug = iota
	// LvlInfo is the informational level
	LvlInfo
	// LvlWarn is the warning level
	LvlWarn
	// LvlError is the error level
	LvlError
	// LvlCrit is the critical level
	LvlCrit

	// FldMessage is the name of the log field that contains the log message
	FldMessage = "msg"
	// FldError is the name of the log field that contains an error
	FldError = "err"
	// FldTimestamp is the name of the log field that contains the timestamp of the log entry
	FldTimestamp = "ts"
	// FldFile is the name of the log field for storing file name information
	FldFile = "file"
	// FldPath is the name of the log field for storing path name information
	FldPath = "path"
	// FldSession is the name of the log field for storing the session ID
	FldSession = "session"
	// FldUser is the name of the log field for storing the ID of the currently active user
	FldUser = "user"
	// FldVersion is the version number of the application
	FldVersion = "ver"
)

// Logger is a helper that uses a levels logger to perform logging operations
// It is used to simplify the typical logging operations a little bit
type Logger struct {
	logger   levels.Levels
	minLevel int
}

// New creates a new Logger instance with the given GoKit logger at its base
func New(kitLogger log.Logger, minLevel int) Logger {
	return Logger{levels.New(kitLogger), minLevel}
}

// Performs logging with the given logger
func doLog(logger log.Logger, msg string, keyvals ...interface{}) {
	keyvals = append(keyvals, FldMessage, msg)
	logger.Log(keyvals...)
}

// RawLogger returns the levels logger used as base for the logger
func (l *Logger) RawLogger() levels.Levels {
	return l.logger
}

// Debug logs a Debug message
func (l *Logger) Debug(msg string, keyvals ...interface{}) {
	if l.minLevel <= LvlDebug {
		doLog(l.logger.Debug(), msg, keyvals...)
	}
}

// Info logs an Informational message
func (l *Logger) Info(msg string, keyvals ...interface{}) {
	if l.minLevel <= LvlInfo {
		doLog(l.logger.Info(), msg, keyvals...)
	}
}

// Warn logs a Warning message
func (l *Logger) Warn(msg string, keyvals ...interface{}) {
	if l.minLevel <= LvlWarn {
		doLog(l.logger.Warn(), msg, keyvals...)
	}
}

// Error logs an Error message
func (l *Logger) Error(msg string, err error, keyvals ...interface{}) {
	if l.minLevel <= LvlError {
		keyvals = append(keyvals, FldError, err)
		doLog(l.logger.Error(), msg, keyvals...)
	}
}

// Crit logs a Critical message
func (l *Logger) Crit(msg string, keyvals ...interface{}) {
	if l.minLevel <= LvlCrit {
		doLog(l.logger.Crit(), msg, keyvals...)
	}
}

// Levels returns the internal levels logger for use inside of GoKit
func (l *Logger) Levels() levels.Levels {
	return l.logger
}

// With returns a new Logger instance with the given key-value pairs set in its context
func (l *Logger) With(keyvals ...interface{}) Logger {
	return Logger{l.logger.With(keyvals...), l.minLevel}
}
