package responder

import (
	"encoding/json"
	"net/http"
)

func Respond(rw http.ResponseWriter, code int, data interface{}) {
	rw.WriteHeader(code)

	if data != nil {
		if err := json.NewEncoder(rw).Encode(data); err != nil {
			rw.WriteHeader(http.StatusInternalServerError)
		}
	}
}

func AuthError(w http.ResponseWriter, r *http.Request, message string) {
	type authError struct {
		Message string `json:"message"`
		Code    int    `json:"code"`
	}

	const code = http.StatusInternalServerError

	Respond(w, code, struct {
		Error authError `json:"error"`
	}{
		Error: authError{
			Message: message,
			Code:    code,
		},
	})
}

func Error(w http.ResponseWriter, r *http.Request, code int, err error) {
	if _, ok := err.(json.Marshaler); ok {
		Respond(w, code, struct {
			Errors error `json:"errors"`
		}{
			Errors: err,
		})
		return
	}

	Respond(w, code, struct {
		Error string `json:"error"`
	}{
		Error: err.Error(),
	})
}
