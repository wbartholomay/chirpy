package main

import (
	"database/sql"
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
		RefreshToken string `json:"refresh_token"`
	}

	user := getUserFromDBUser(dbUser)

	token, err := auth.MakeJWT(user.ID, cfg.tokenSecret)
	if err != nil {
		return APIError{
			Status: 401,
			ResponseMsg: "Unauthorized",
			ErrorMsg: err.Error(),
		}
	}

	refreshToken, err := cfg.db.GetRefreshTokenByUser(req.Context(), user.ID)
	if err != nil && err != sql.ErrNoRows{
		return getDefaultApiError(err)
	}

	//if no token is found for the user, or if the token is expired, return a new token
	refreshTokenString := ""
	if err == sql.ErrNoRows || refreshToken.RevokedAt.Valid{
		refreshTokenString = auth.MakeRefreshToken()
	} else {
		refreshTokenString = refreshToken.Token
	}

	resData := userWithToken{
		User: user,
		Token: token,
		RefreshToken: refreshTokenString,
	}


	respondWithJSON(w, 200, resData)
	return nil
}