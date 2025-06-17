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
	Msg string
}

func (err APIError) Error() string {
	return err.Msg
}

func getDefaultApiError(err error) APIError{
	return APIError{
		Status: 500,
		Msg: err.Error(),
	}
}

func makeHandler(h apiFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if err := h(w, r); err != nil {
			e, ok := err.(APIError)
			if !ok {
				slog.Error("API failed with non API error.", "err", err)
				respondWithError(w, 500, "Something went wrong.", err)
			} else {
				slog.Error("API error", "err", e, "status", e.Status)
				switch e.Status{
					case 500:
						respondWithError(w, 500, "Something went wrong.", err)
					default:
						respondWithError(w, e.Status, err.Error(), err)
				}
			}
		}
	}
}

func respondWithError(w http.ResponseWriter, code int, msg string, err error) {
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
