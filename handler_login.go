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

	type userWithToken struct{
		User
		Token string `json:"token"`
	}

	user := getUserFromDBUser(dbUser)
	token, err := auth.GetBearerToken(req.Header)
	if err != nil {
		return APIError{
			Status: 401,
			Msg: "Unauthorized",
		}
	}

	resData := userWithToken{
		User: user,
		Token: token,
	}


	respondWithJSON(w, 200, resData)
	return nil
}