package logging

import (
	"fmt"
	"net/http"
)

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
