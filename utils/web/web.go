package web

import (
	"encoding/json"
	"fmt"
	"net/http"
)

type Validator interface {
	Validate() error
}

func DecodeHTTPRequest(r *http.Request, target Validator) error {
	if err := json.NewDecoder(r.Body).Decode(target); err != nil {
		return err
	}
	return target.Validate()
}

func EncodeHTTPResponse(w http.ResponseWriter, r *http.Request, target interface{}) error {
	enc := json.NewEncoder(w)
	enc.SetIndent("", "  ")
	return enc.Encode(target)
}

func HTTPError(w http.ResponseWriter, statusCode int, messages ...any) {
	var errText string
	if messages != nil {
		for _, message := range messages {
			errText += fmt.Sprintf("%s, ", message)
		}
		errText = errText[:len(errText)-2]
	} else {
		errText = fmt.Sprintf("%d %s", statusCode, http.StatusText(statusCode))
	}
	http.Error(w, errText, statusCode)
}
