package main

import (
	"encoding/json"
	"errors"
	"log"
	"net/http"
	"strings"

	"github.com/google/uuid"
	"github.com/wbartholomay/chirpy/internal/database"
)

func (cfg *apiConfig) CreateChirpHandler(w http.ResponseWriter, req *http.Request) {
	type parameters struct {
		Body string `json:"body"`
		UserID uuid.UUID `json:"user_id"`
	}


	decoder := json.NewDecoder(req.Body)
	params := parameters{}

	if err := decoder.Decode(&params); err != nil {
		log.Printf("Error decoding parameters: %v\n", err)
		respondWithError(w, 500, "Something went wrong", err)
		return
	}

	cleanedChirp, err := validateAndCleanChirp(w, params.Body)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, err.Error(), err)
	}


	dbChirp, err := cfg.db.CreateChirp(req.Context(), database.CreateChirpParams{
		Body: cleanedChirp,
		UserID: params.UserID,
	})
	if err != nil {
		respondWithDefaultError(w, err)
	}

	chirp := getChirpFromDBChirp(dbChirp)


	respondWithJSON(w, http.StatusCreated, chirp)
	
}

func validateAndCleanChirp(w http.ResponseWriter, body string) (string, error){
	
	if len(body) > 140 {
		return "", errors.New("chirp is too long")
	}

	return getCleanedBody(body), nil
	
}

func getCleanedBody(body string) string {
	profanity := []string{"kerfuffle", "sharbert", "fornax"}

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