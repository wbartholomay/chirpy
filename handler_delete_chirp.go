package main

import (
	"fmt"
	"net/http"

	"github.com/google/uuid"
)


func (cfg *apiConfig) DeleteChirpHandler(w http.ResponseWriter, req *http.Request) error {
	user, err := cfg.authenticateAndGetUser(req)
	if err != nil {
		return APIError{
			Status: 401,
			ResponseMsg: "Unauthorized",
			ErrorMsg: err.Error(),
		}
	}

	chirp_id, err := uuid.Parse(req.PathValue("chirp_id"))
	if err != nil {
		return getDefaultApiError(err)
	}

	chirp, err := cfg.db.GetChirp(req.Context(), chirp_id)
	if err != nil {
		return APIError{
			Status: 404,
			ResponseMsg: fmt.Sprintf("no chirp found with chirp id: %v", chirp_id),
			ErrorMsg: err.Error(),
		}
	}

	if chirp.UserID != user.ID {
		return APIError{
			Status: 403,
			ResponseMsg: "Forbidden",
			ErrorMsg: "user not authorized to delete this resource",
		}
	}

	if err = cfg.db.DeleteChirp(req.Context(), chirp.ID); err != nil {
		return getDefaultApiError(err)
	}

	respondWithJSON(w, 204, nil)
	return nil
}