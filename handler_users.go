package main

import (
	"encoding/json"
	"net/http"
)

func (cfg *apiConfig) CreateUserHandler (w http.ResponseWriter, req *http.Request) error{
	type params struct {
		Email string `json:"email"`
	}

	decoder := json.NewDecoder(req.Body)
	reqParams := params{}
	err := decoder.Decode(&reqParams)
	if err != nil {
		return getDefaultApiError(err)
	}

	dbUser, err := cfg.db.CreateUser(req.Context(), reqParams.Email)
	if err != nil {
		return getDefaultApiError(err)
	}

	user := getUserFromDBUser(dbUser)

	respondWithJSON(w, 201, user)
	return nil
}