package main

import (
	"encoding/json"
	"errors"
	"net/http"
	"strings"

	"github.com/wbartholomay/chirpy/internal/database"
)

func (cfg *apiConfig) CreateChirpHandler(w http.ResponseWriter, req *http.Request) error {
	type parameters struct {
		Body string `json:"body"`
	}


	user, err := cfg.authenticateAndGetUser(req)
	if err != nil {
		return APIError{
			Status: 401,
			ResponseMsg: "Unauthorized",
			ErrorMsg: err.Error(),
		}
	}

	decoder := json.NewDecoder(req.Body)
	params := parameters{}

	if err := decoder.Decode(&params); err != nil {
		return getDefaultApiError(err)
	}

	cleanedChirp, err := validateAndCleanChirp(params.Body)
	if err != nil {
		return APIError{
			Status: http.StatusBadRequest,
			ResponseMsg: err.Error(),
			ErrorMsg: err.Error(),
		}
	}

	dbChirp, err := cfg.db.CreateChirp(req.Context(), database.CreateChirpParams{
		Body: cleanedChirp,
		UserID: user.ID,
	})
	if err != nil {
		return getDefaultApiError(err)
	}

	chirp := getChirpFromDBChirp(dbChirp)


	respondWithJSON(w, http.StatusCreated, chirp)
	return nil
}

func validateAndCleanChirp(body string) (string, error){
	
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