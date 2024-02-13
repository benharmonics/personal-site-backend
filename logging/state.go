package logging

import (
	"fmt"
	"time"
)

const (
	debugText = "DEBUG"
	infoText  = "INFO"
	warnText  = "WARNING"
	errorText = "ERROR"

	red             = "\033[31m"
	boldRed         = "\033[1;31m"
	green           = "\033[32m"
	boldGreen       = "\033[1;32m"
	brightBoldGreen = "\033[1;92m"
	boldYellow      = "\033[1;33m"
	brightBoldBlue  = "\033[1;94m"
	reset           = "\033[0m"
)

var (
	color, timestamp bool
	timeformat       string = time.RFC3339
	logLevel         string = LogLevelInfo

	debugString = fmt.Sprintf("%s: ", debugText)
	infoString  = fmt.Sprintf("%s: ", infoText)
	warnString  = fmt.Sprintf("%s: ", warnText)
	errorString = fmt.Sprintf("%s: ", errorText)
)

func SetTime(active bool)         { timestamp = active }
func SetTimeFormat(format string) { timeformat = format }

func SetLogLevel(level string) error {
	switch level {
	case LogLevelDebug, LogLevelInfo, LogLevelWarn, LogLevelError:
		logLevel = level
		return nil
	default:
		return fmt.Errorf("unknown log level %s", level)
	}
}

func SetColor(active bool) {
	if color = active; color {
		debugString = fmt.Sprintf("%s%s%s: ", brightBoldBlue, debugText, reset)
		infoString = fmt.Sprintf("%s%s%s: ", brightBoldGreen, infoText, reset)
		warnString = fmt.Sprintf("%s%s%s: ", boldYellow, warnText, reset)
		errorString = fmt.Sprintf("%s%s%s: ", boldRed, errorText, reset)
	} else {
		debugString = fmt.Sprintf("%s: ", debugText)
		infoString = fmt.Sprintf("%s: ", infoText)
		warnString = fmt.Sprintf("%s: ", warnText)
		errorString = fmt.Sprintf("%s: ", errorText)
	}
}
