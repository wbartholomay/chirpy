package main

import (
	"encoding/json"
	"log"
	"net/http"
	"strings"

	"github.com/wbartholomay/chirpy/internal/httphelper"
)



func ValidateChirpHandler(w http.ResponseWriter, req *http.Request) {
	profanity := []string{"kerfuffle", "sharbert", "fornax"}

	type parameters struct {
		Body string `json:"body"`
	}

	type resSuccess struct {
		CleanedBody string `json:"cleaned_body"`
	}

	decoder := json.NewDecoder(req.Body)
	params := parameters{}

	if err := decoder.Decode(&params); err != nil {
		log.Printf("Error decoding parameters: %v\n", err)
		httphelper.RespondWithError(w, 500, "Something went wrong", err)
		return
	}

	if len(params.Body) > 140 {
		httphelper.RespondWithError(w, 400, "Chirp is too long", nil)
		return
	}

	w.WriteHeader(200)
	cleanedBody := getCleanedBody(params.Body, profanity)
	successParams := resSuccess{
		CleanedBody: cleanedBody,
	}

	httphelper.RespondWithJSON(w, 200, successParams)
}

func getCleanedBody(body string, profanity []string) string {
	bodySlice := strings.Split(body, " ")
	for i, word := range bodySlice {
		lowerWord := strings.ToLower(word)
		for _, badWord := range profanity {
			if lowerWord == badWord {
				bodySlice[i] = "****"
			}
		}
	}

	cleanedBody := strings.Join(bodySlice, " ")
	
	return cleanedBody
}