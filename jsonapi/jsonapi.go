package jsonapi

import (
	"bytes"
	"database/sql"
	"encoding/json"
	"io"
	"log"
	"net/http"

	"github.com/OtaviOuu/mailingList-microservice/mdb"
)

func setJsonHeader(w http.ResponseWriter) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
}

// io.Reader -> qualquer coisa que pode ser lida
func fromJson[T any](body io.Reader, target T) {
	buf := new(bytes.Buffer)
	buf.ReadFrom(body)

	json.Unmarshal(buf.Bytes(), &target)
}

func returnJson[T any](w http.ResponseWriter, withData func() (T, error)) {
	setJsonHeader(w)
	// Retorno da função withData é o conteudo do json
	data, err := withData()
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		ErrJson, err := json.Marshal(&err)
		if err != nil {
			log.Println(err)
			return
		}
		w.Write(ErrJson)
	}

	dataJson, err := json.Marshal(&data)
	if err != nil {
		log.Print(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.Write(dataJson)
}

func returnErr(w http.ResponseWriter, err error, statusCode int) {
	returnJson(w, func() (interface{}, error) {
		errorMessage := struct {
			Err string
		}{
			Err: err.Error(),
		}
		w.WriteHeader(statusCode)
		return errorMessage, nil
	})
}

func CreateEmail(db *sql.DB) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			return
		}

		entry := mdb.EmailEntry{}
		fromJson(r.Body, &entry)

		if err := mdb.CreateEmail(db, entry.Email); err != nil {
			returnErr(w, err, 400)
			return
		}

		returnJson(w, func() (interface{}, error) {
			log.Printf("Json CreateEmail: %v\n", entry.Email)
			return mdb.GetEmail(db, entry.Email)
		})
	})
}
