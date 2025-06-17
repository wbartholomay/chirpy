package main

import (
	"encoding/json"
	"net/http"

	"github.com/wbartholomay/chirpy/internal/auth"
	"github.com/wbartholomay/chirpy/internal/database"
)

func (cfg *apiConfig) CreateUserHandler (w http.ResponseWriter, req *http.Request) error{
	type params struct {
		Email string `json:"email"`
		Password string `json:"password"`
	}

	decoder := json.NewDecoder(req.Body)
	reqParams := params{}
	err := decoder.Decode(&reqParams)
	if err != nil {
		return getDefaultApiError(err)
	}

	hashedPassword, err := auth.HashPassword(reqParams.Password)
	if err != nil {
		return getDefaultApiError(err)
	}

	dbParams := database.CreateUserParams{
		Email: reqParams.Email,
		HashedPassword: hashedPassword,
	}
	
	dbUser, err := cfg.db.CreateUser(req.Context(), dbParams)
	if err != nil{
		return getDefaultApiError(err)
	}

	user := getUserFromDBUser(dbUser)

	respondWithJSON(w, 201, user)
	return nil
}