package httphelper

import (
	"encoding/json"
	"log"
	"net/http"
)

func RespondWithError(w http.ResponseWriter, statusCode int, errMessage string) {
	type resError struct {
		Error string `json:"error"`
	}

	w.WriteHeader(statusCode)

	resData, err := json.Marshal(resError{
			Error: errMessage,
	})
	if err != nil {
		log.Fatal(err)
	}

	w.Write(resData)
}