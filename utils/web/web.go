package web

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
)

type Validator interface {
	Validate() error
}

func DecodeHTTPRequest(r *http.Request, decodeTo Validator) error {
	if err := json.NewDecoder(r.Body).Decode(decodeTo); err != nil {
		return err
	}
	return decodeTo.Validate()
}

func HTTPError(w http.ResponseWriter, statusCode int, messages ...string) {
	var errText string
	if messages != nil {
		errText = strings.Join(messages, " ")
	} else {
		errText = fmt.Sprintf("%d %s", statusCode, http.StatusText(statusCode))
	}
	http.Error(w, errText, statusCode)
}
