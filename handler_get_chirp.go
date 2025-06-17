package main

import (
	"net/http"

	"github.com/google/uuid"
)

func (cfg *apiConfig) GetAllChirpsHandler(w http.ResponseWriter, req *http.Request) error {
	dbChirps, err := cfg.db.GetChirps(req.Context())
	if err != nil {
		return getDefaultApiError(err)
	}

	chirps := []Chirp{}

	//convert to json friendly struct
	for _, dbChirp := range dbChirps {
		chirps = append(chirps, getChirpFromDBChirp(dbChirp))
	}

	respondWithJSON(w, 200, chirps)
	return nil
}

func (cfg *apiConfig) GetChirpByIdHandler(w http.ResponseWriter, req *http.Request) error {
	chirp_id, err := uuid.Parse(req.PathValue("chirp_id"))
	if err != nil {
		return getDefaultApiError(err)
	}

	dbChirp, err := cfg.db.GetChirp(req.Context(), chirp_id)
	if err != nil {
		return APIError{
			Status: http.StatusNotFound,
			ResponseMsg: "Chirp not found",
			ErrorMsg: err.Error(),
		}
	}

	chirp := getChirpFromDBChirp(dbChirp)
	respondWithJSON(w, 200, chirp)

	return nil
}