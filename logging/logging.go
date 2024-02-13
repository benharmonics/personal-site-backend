package logging

import (
	"fmt"
	"time"
)

const (
	LogLevelDebug = "debug"
	LogLevelInfo  = "info"
	LogLevelWarn  = "warn"
	LogLevelError = "error"
)

func Debug(a ...any) {
	if logLevel != LogLevelDebug {
		return
	}
	fmt.Print(debugString, newTimestamp(), fmt.Sprintln(a...))
}

func Debugf(format string, a ...any) {
	if logLevel != LogLevelDebug {
		return
	}
	fmt.Print(debugString, newTimestamp(), fmt.Sprintf(format, a...))
}

func Info(a ...any) {
	if logLevel == LogLevelWarn || logLevel == LogLevelError {
		return
	}
	fmt.Print(infoString, newTimestamp(), fmt.Sprintln(a...))
}

func Infof(format string, a ...any) {
	if logLevel == LogLevelWarn || logLevel == LogLevelError {
		return
	}
	fmt.Print(infoString, newTimestamp(), fmt.Sprintf(format, a...))
}

func Warn(a ...any) {
	if logLevel == LogLevelError {
		return
	}
	fmt.Print(warnString, newTimestamp(), fmt.Sprintln(a...))
}

func Warnf(format string, a ...any) {
	if logLevel == LogLevelError {
		return
	}
	fmt.Print(warnString, newTimestamp(), fmt.Sprintf(format, a...))
}

func Error(a ...any) {
	fmt.Print(errorString, newTimestamp(), fmt.Sprintln(a...))
}

func Errorf(format string, a ...any) {
	fmt.Print(errorString, newTimestamp(), fmt.Sprintf(format, a...))
}

func newTimestamp() string {
	if !timestamp {
		return ""
	}
	return fmt.Sprintf("(%s) ", time.Now().Format(timeformat))
}
