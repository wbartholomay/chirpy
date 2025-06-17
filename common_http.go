package main

import (
	"encoding/json"
	"log"
	"log/slog"
	"net/http"
)

type apiFunc func(w http.ResponseWriter, r *http.Request) error

type APIError struct {
	Status int
	ResponseMsg string
	ErrorMsg string
}

func (err APIError) Error() string {
	return err.ErrorMsg
}

func getDefaultApiError(err error) APIError{
	return APIError{
		Status: 500,
		ResponseMsg: "Internal server error",
		ErrorMsg: err.Error(),
	}
}


func makeHandler(h apiFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if err := h(w, r); err != nil {
			e, ok := err.(APIError)
			if !ok {
				slog.Error("API failed with non API error.", "err", err)
				respondWithError(w, 500, "Something went wrong.")
			} else {
				slog.Error("API error", "err", e, "status", e.Status)
				respondWithError(w, e.Status, e.ResponseMsg)
			}
		}
	}
}

func respondWithError(w http.ResponseWriter, code int, msg string) {
	type errorResponse struct {
		Error string `json:"error"`
	}
	respondWithJSON(w, code, errorResponse{
		Error: msg,
	})
}

func respondWithJSON(w http.ResponseWriter, code int, payload interface{}) {
	w.Header().Set("Content-Type", "application/json")
	dat, err := json.Marshal(payload)
	if err != nil {
		log.Printf("Error marshalling JSON: %s", err)
		w.WriteHeader(500)
		return
	}
	w.WriteHeader(code)
	w.Write(dat)
}
