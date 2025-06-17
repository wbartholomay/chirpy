package main

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/wbartholomay/chirpy/internal/auth"
)

func (cfg *apiConfig) LoginUserHandler (w http.ResponseWriter, req *http.Request) error{
	type params struct {
		Email string `json:"email"`
		Password string `json:"password"`
		ExpiresInSeconds int `json:"expires_in_seconds"`
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
			ResponseMsg: "incorrect email or password",
			ErrorMsg: err.Error(),
		}
	}

	err = auth.CheckPasswordHash(dbUser.HashedPassword, reqParams.Password)
	if err != nil {
		return APIError{
			Status: 401,
			ResponseMsg: "incorrect email or password",
			ErrorMsg: err.Error(),
		}
	}

	type userWithToken struct{
		User
		Token string `json:"token"`
	}

	user := getUserFromDBUser(dbUser)

	expiresIn := time.Hour
	if 0 < reqParams.ExpiresInSeconds && reqParams.ExpiresInSeconds < 3600{
		expiresIn = time.Duration(reqParams.ExpiresInSeconds) * time.Second
	}

	token, err := auth.MakeJWT(user.ID, cfg.tokenSecret, expiresIn)
	if err != nil {
		return APIError{
			Status: 401,
			ResponseMsg: "Unauthorized",
			ErrorMsg: err.Error(),
		}
	}

	resData := userWithToken{
		User: user,
		Token: token,
	}


	respondWithJSON(w, 200, resData)
	return nil
}