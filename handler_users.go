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

func (cfg *apiConfig) UpdateUserCredentialsHandler (w http.ResponseWriter, req *http.Request) error {
	user, err := cfg.authenticateAndGetUser(req)
	if err != nil {
		return APIError{
			Status: 401,
			ResponseMsg: "Unauthorized",
			ErrorMsg: err.Error(),
		}
	}

	type reqParams struct {
		Email string `json:"email"`
		Password string `json:"password"`
	}

	p := reqParams{}

	decoder := json.NewDecoder(req.Body)
	if 	err = decoder.Decode(&p); err != nil {
		return getDefaultApiError(err)
	}

	//hash password
	hashedPassword, err := auth.HashPassword(p.Password)
	if err != nil {
		return getDefaultApiError(err)
	}

	params := database.UpdateEmailAndPasswordParams{
		ID: user.ID,
		Email: p.Email,
		HashedPassword: hashedPassword ,
	}

	updatedUser, err := cfg.db.UpdateEmailAndPassword(req.Context(), params); 
	if err != nil {
		return getDefaultApiError(err)
	}

	resData := getUserFromDBUser(updatedUser)

	respondWithJSON(w, 200, resData)
	return nil
}