package main

import (
	"encoding/json"
	"net/http"

	"github.com/wbartholomay/chirpy/internal/auth"
)

func (cfg *apiConfig) LoginUserHandler (w http.ResponseWriter, req *http.Request) error{
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

	dbUser, err := cfg.db.GetUserByEmail(req.Context(), reqParams.Email)
	if err != nil {
		return APIError{
			Status: 401,
			Msg: "incorrect email or password",
		}
	}

	err = auth.CheckPasswordHash(dbUser.HashedPassword, reqParams.Password)
	if err != nil {
		return APIError{
			Status: 401,
			Msg: "incorrect email or password",
		}
	}

	user := getUserFromDBUser(dbUser)
	respondWithJSON(w, 200, user)
	return nil
}