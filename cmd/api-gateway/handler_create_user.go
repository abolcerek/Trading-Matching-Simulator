package main

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/abolcerek/Trading-Matching-Simulator/internal/auth"
	"github.com/abolcerek/Trading-Matching-Simulator/internal/database"
	"github.com/google/uuid"
)

type User struct {
	Password string `json:"password"`
	Email string `json:"email"`
}

type User_response struct {
	ID uuid.UUID `json:"id"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	Email string `json:"email"`
	Balance int64 `json:"balance"`
}



func (cfg *apiConfig) HandlerCreateUser(w http.ResponseWriter, r *http.Request) {
	req_params := User{}
	err_params := error_parameters{}
	decoder := json.NewDecoder(r.Body)
	err := decoder.Decode(&req_params)
	if err != nil {
		err_params.Error = "Something went wrong"
		handleErrors(w, &err_params, 400)
		return
	}
	hashed_password, err := auth.MakeHashPassword(req_params.Password)
	if err != nil {
		err_params.Error = "Something went wrong"
		handleErrors(w, &err_params, 400)
		return
	}
	create_user_params := database.CreateUserParams{
		ID: uuid.New(),
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
		Email: req_params.Email,
		HashedPassword: hashed_password,
		Balance: balance,
	}
	database_user, err := cfg.database.CreateUser(r.Context(), create_user_params)
	if err != nil {
		err_params.Error = "Something went wrong"
		handleErrors(w, &err_params, 400)
		return
	}
	user_resp := User_response{
		ID: database_user.ID,
		CreatedAt: database_user.CreatedAt,
		UpdatedAt: database_user.UpdatedAt,
		Email: database_user.Email,
		Balance: database_user.Balance,
	}
	data, err := json.Marshal(user_resp)
	if err != nil {
		err_params.Error = "Something went wrong"
		handleErrors(w, &err_params, 400)
		return
	}
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(201)
	w.Write(data)
}