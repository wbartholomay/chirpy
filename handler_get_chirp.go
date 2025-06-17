package main

import "net/http"

func (cfg *apiConfig) GetChirpsHandler(w http.ResponseWriter, req *http.Request) error {
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