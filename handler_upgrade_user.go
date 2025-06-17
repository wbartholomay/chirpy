package main

import (
	"encoding/json"
	"net/http"

	"github.com/google/uuid"
	"github.com/wbartholomay/chirpy/internal/auth"
)

func (cfg *apiConfig) UpgradeUserHandler(w http.ResponseWriter, req *http.Request) error {
	apiKey, err := auth.GetAPIKey(req.Header)
	if err != nil || apiKey != cfg.polkaKey{
		return APIError{
			Status: 401,
			ResponseMsg: "Unauthorized",
			ErrorMsg: err.Error(),
		}
	}


	type reqParams struct {
		Event string `json:"event"`
		Data struct{
			UserID string `json:"user_id"`
		} `json:"data"`
	}

	decoder := json.NewDecoder(req.Body)
	
	reqData := reqParams{}
	err = decoder.Decode(&reqData)
	if err != nil {
		return getDefaultApiError(err)
	}

	if reqData.Event != "user.upgraded" {
		respondWithJSON(w, 204, nil)
		return nil
	}

	userID, err := uuid.ParseBytes(([]byte)(reqData.Data.UserID))
	if err != nil {
		return getDefaultApiError(err)
	}

	_, err = cfg.db.UpgradeToChirpyRed(req.Context(), userID)
	if err != nil {
		return APIError{
			Status: 404,
			ResponseMsg: "no user with id: " + reqData.Data.UserID,
			ErrorMsg: err.Error(),
		}
	}

	respondWithJSON(w, 204, nil)
	return nil
}