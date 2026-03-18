package httpapi

import (
	"encoding/json"
	"net/http"
)

type APIError struct {
	Code string `json:"code"`
	Msg  string `json:"msg"`
}

type APIErr struct {
	Status int
	Code   string
}

var (
	ErrBadRequest   = APIErr{Status: http.StatusBadRequest, Code: "BAD_REQUEST"}
	ErrUnauthorized = APIErr{Status: http.StatusUnauthorized, Code: "UNAUTHORIZED"}
	ErrInternal     = APIErr{Status: http.StatusInternalServerError, Code: "INTERNAL_ERROR"}
)

func (e APIErr) Respond(w http.ResponseWriter, msg string) {
	writeJSON(w, e.Status, APIError{Code: e.Code, Msg: msg})
}

func RespondOK(w http.ResponseWriter, data any) {
	writeJSON(w, http.StatusOK, data)
}

func writeJSON(w http.ResponseWriter, status int, data any) {
	b, err := json.Marshal(data)
	if err != nil {
		w.Header().Set("Content-Type", "application/json; charset=utf-8")
		w.WriteHeader(http.StatusInternalServerError)
		_, _ = w.Write([]byte(`{"code":"INTERNAL_ERROR","msg":"response encoding failed"}`))
		return
	}
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(status)
	_, _ = w.Write(b)
}
