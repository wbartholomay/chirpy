package main

import (
	"net/http"
	"sort"

	"github.com/google/uuid"
)

func (cfg *apiConfig) GetAllChirpsHandler(w http.ResponseWriter, req *http.Request) error {
	authorID := req.URL.Query().Get("author_id")
	if authorID != "" {
		return cfg.getChirpsFromAuthor(w, req, authorID)
	}

	dbChirps, err := cfg.db.GetChirps(req.Context())
	if err != nil {
		return getDefaultApiError(err)
	}

	chirps := []Chirp{}

	//convert to json friendly struct
	for _, dbChirp := range dbChirps {
		chirps = append(chirps, getChirpFromDBChirp(dbChirp))
	}

	chirps = sortChirps(req, chirps)

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

func (cfg *apiConfig) getChirpsFromAuthor(w http.ResponseWriter, req *http.Request, authorID string) error {
	//get author uuid
	userID, err := uuid.Parse(authorID)
	if err != nil {
		return getDefaultApiError(err)
	}

	//look for all chirps by that author
	dbChirps, err := cfg.db.GetChirpsByUserID(req.Context(), userID)
	if err != nil {
		return APIError{
			Status: 404,
			ResponseMsg: "no chirps found for user id: " + authorID,
			ErrorMsg: err.Error(),
		}
	}

	chirps := []Chirp{}

	for _, c := range dbChirps {
		chirps = append(chirps, getChirpFromDBChirp(c))
	}

	chirps = sortChirps(req, chirps)

	respondWithJSON(w, 200, chirps)
	return nil
}

func sortChirps(req *http.Request, chirps []Chirp) []Chirp {
	sortParam := req.URL.Query().Get("sort")
	if sortParam == "desc" {
		sort.Slice(chirps, func(i, j int) bool {
			return chirps[i].CreatedAt.After(chirps[j].CreatedAt)
		})
	}

	return chirps
}