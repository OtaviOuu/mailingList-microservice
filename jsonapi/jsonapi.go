package jsonapi

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"
)

func setJsonHeader(w http.ResponseWriter) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
}

// io.Reader -> qualquer coisa que pode ser lida
func fromJson[T any](body io.Reader, target T) error {
	buf := new(bytes.Buffer)
	buf.ReadFrom(body)

	json.Unmarshal(buf.Bytes(), &target)
}
