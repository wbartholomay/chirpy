package main

import (
	"net/http"

	"github.com/wbartholomay/chirpy/internal/auth"
	"github.com/wbartholomay/chirpy/internal/database"
)

func (cfg *apiConfig) authenticateAndGetUser(req *http.Request) (database.User, error) {
	token, err := auth.GetBearerToken(req.Header)
	if err != nil {
		return database.User{}, err
	}

	userID, err := auth.ValidateJWT(token, cfg.tokenSecret)
	if err != nil {
		return database.User{}, err
	}

	return cfg.db.GetUserByID(req.Context(), userID)
}