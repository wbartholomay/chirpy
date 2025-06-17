package main

import (
	"fmt"
	"net/http"

	"github.com/wbartholomay/chirpy/internal/auth"
)

func (cfg *apiConfig) RefreshTokenHandler (w http.ResponseWriter, req *http.Request) error {
	refTokenString, err := auth.GetBearerToken(req.Header)
	if err != nil {
		return getDefaultApiError(err)
	}

	refToken, err := cfg.db.GetRefreshTokenByID(req.Context(), refTokenString)
	if err != nil {
		return APIError{
			Status: 401,
			ResponseMsg: "Unauthorized",
			ErrorMsg: err.Error(),
		}
	}
	//if token is revoked, respond with unauthorized
	if refToken.RevokedAt.Valid { 
		return APIError{
			Status: 401,
			ResponseMsg: "Unauthorized",
			ErrorMsg: "refresh token is revoked",
		}
	}

	//get user from refresh token, then create jwt for the user
	user, err := cfg.db.GetUserFromRefreshToken(req.Context(), refToken.Token)
	if err != nil {
		return getDefaultApiError(err)
	}

	token, err := auth.MakeJWT(user.ID, cfg.tokenSecret) 
	if err != nil {
		return getDefaultApiError(fmt.Errorf("error creating JWT for user: %w", err))
	}

	type tokenResponse struct {
		Token string `json:"token"`
	}

	resData := tokenResponse {
		Token: token,
	}

	respondWithJSON(w, 200, resData)
	return nil
}