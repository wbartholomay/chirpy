package main

import (
	"fmt"
	"net/http"

	"github.com/wbartholomay/chirpy/internal/auth"
)

func (cfg *apiConfig) RevokeTokenHandler (w http.ResponseWriter, req *http.Request) error {
	token, err := auth.GetBearerToken(req.Header)
	if err != nil {
		return APIError{
			Status: 401,
			ResponseMsg: "Unauthorized",
			ErrorMsg: err.Error(),
		}
	}

	err = cfg.db.RevokeRefreshToken(req.Context(), token)
	if err != nil {
		return getDefaultApiError(fmt.Errorf("error revoking token: %w", err))
	}

	respondWithJSON(w, 204, nil)
	return nil
}