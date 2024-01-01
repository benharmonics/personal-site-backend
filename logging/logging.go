package logging

import (
	"fmt"
	"net/http"
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

var color, debug bool
var debugString, infoString, warnString, errorString string

func init() {
	SetColor(false)
}

func SetDebug(active bool) {
	debug = active
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

func Debug(messages ...any) {
	if !debug {
		return
	}
	fmt.Print(debugString, fmt.Sprintln(messages...))
}

func Info(messages ...any) {
	fmt.Print(infoString, fmt.Sprintln(messages...))
}

func Warning(messages ...any) {
	fmt.Print(warnString, fmt.Sprintln(messages...))
}

func Error(messages ...any) {
	fmt.Print(errorString, fmt.Sprintln(messages...))
}

func HTTPOk(r *http.Request) {
	if !color {
		fmt.Printf("%s - \"%s %s %s\" 200 OK\n", r.RemoteAddr, r.Method, r.Proto, r.URL)
		return
	}
	fmt.Printf("%s - \"%s%s %s %s%s\" %s200 OK%s\n", r.RemoteAddr, boldGreen, r.Method, r.Proto, r.URL, reset, green, reset)
}

func HTTPError(r *http.Request, statusCode int) {
	errText := fmt.Sprintf("%d %s", statusCode, http.StatusText(statusCode))
	if !color {
		fmt.Printf("%s - \"%s %s %s\" %s\n", r.RemoteAddr, r.Method, r.Proto, r.URL, errText)
		return
	}
	fmt.Printf("%s - \"%s%s %s %s%s\" %s%s%s\n", r.RemoteAddr, boldGreen, r.Method, r.Proto, r.URL, reset, red, errText, reset)
}
