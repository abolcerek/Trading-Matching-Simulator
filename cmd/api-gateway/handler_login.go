package main

import (
	"net/http"
	"encoding/json"
	"time"
	"github.com/abolcerek/Trading-Matching-Simulator/internal/auth"
	"github.com/google/uuid"
)

type UserLogin struct {
	ID uuid.UUID `json:"id"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	Email string `json:"email"`
	Token string `json:"token"`
	Balance int64 `json:"balance"`
}

func (cfg *apiConfig) HandlerLogin(w http.ResponseWriter, r *http.Request) {
	err_params := error_parameters{}
	user_parms := User{}
	decoder := json.NewDecoder(r.Body)
	err := decoder.Decode(&user_parms)
	if err != nil {
		err_params.Error = "Something went wrong"
		handleErrors(w, &err_params, 400)
		return
	}
	database_user, err := cfg.database.GetUserByEmail(r.Context(), user_parms.Email)
	if err != nil {
		err_params.Error = "Invalid email or password"
		handleErrors(w, &err_params, 401)
		return
	}
	valid, err := auth.CheckPasswordHash(user_parms.Password, database_user.HashedPassword)
	if err != nil {
		err_params.Error = "Invalid email or password"
		handleErrors(w, &err_params, 401)
		return
	}
	if !valid {
		err_params.Error = "Invalid email or password"
		handleErrors(w, &err_params, 401)
		return
	}
	jwt, err := auth.MakeJWT(database_user.ID, cfg.JWT_secret)
	if err != nil {
		err_params.Error = "Error generating JSON web token"
		handleErrors(w, &err_params, 401)
		return
	}
	login_resp := UserLogin{
		ID: database_user.ID,
		CreatedAt: database_user.CreatedAt,
		UpdatedAt: database_user.UpdatedAt,
		Email: database_user.Email,
		Token: jwt,
		Balance: database_user.Balance,
	}
	data, err := json.Marshal(&login_resp)
	if err != nil {
		err_params.Error = "Something went wrong"
		handleErrors(w, &err_params, 400)
		return
	}
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(200)
	w.Write(data)
}