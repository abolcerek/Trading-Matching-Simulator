package main

import (
	"encoding/json"
	"net/http"

	"github.com/abolcerek/Trading-Matching-Simulator/internal/auth"
	"github.com/abolcerek/Trading-Matching-Simulator/internal/database"
)

type UserUpdate struct {
	Current_Password string `json:"current_password"`
	New_Password string `json:"new_password"`
	Email string `json:"email"`
}

func (cfg *apiConfig) HandlerUpdateUser(w http.ResponseWriter, r *http.Request) {
	err_params := error_parameters{}
	token, err := auth.GetBearerToken(r.Header)
	if err != nil {
		err_params.Error = "Invalid access token"
		handleErrors(w, &err_params, 400)
		return
	}
	id, err := auth.ValidateJWT(token, cfg.JWT_secret)
	if err != nil {
		err_params.Error = "Invalid access token"
		handleErrors(w, &err_params, 400)
		return
	}
	user := UserUpdate{}
	decoder := json.NewDecoder(r.Body)
	err = decoder.Decode(&user)
	if err != nil {
		err_params.Error = "Something went wrong"
		handleErrors(w, &err_params, 400)
		return
	}
	if user.New_Password == "" {
    	err_params.Error = "New password cannot be empty"
    	handleErrors(w, &err_params, 400)
    	return
	}
	db_user, err := cfg.database.GetUserByID(r.Context(), id)
	if err != nil {
		err_params.Error = "Something went wrong"
		handleErrors(w, &err_params, 500)
		return
	}
	valid, err := auth.CheckPasswordHash(user.Current_Password, db_user.HashedPassword)
	if err != nil {
		err_params.Error = "Current password is incorrect"
		handleErrors(w, &err_params, 401)
		return
	}
	if !valid {
		err_params.Error = "Current password is incorrect"
		handleErrors(w, &err_params, 401)
		return
	}
	hashed_password, err := auth.MakeHashPassword(user.New_Password)
	if err != nil {
		err_params.Error = "Error creating password"
		handleErrors(w, &err_params, 401)
		return
	}
	update_user_params := database.UpdateUserParams{
		Email: db_user.Email,
		HashedPassword: hashed_password,
		ID: id,
	}
	err = cfg.database.UpdateUser(r.Context(), update_user_params)
	if err != nil {
		err_params.Error = "Error creating password"
		handleErrors(w, &err_params, 400)
		return
	}
	database_user, err := cfg.database.GetUserByID(r.Context(), id)
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
	data, err := json.Marshal(&user_resp)
	if err != nil {
		err_params.Error = "Something went wrong"
		handleErrors(w, &err_params, 400)
		return
	}
	w.WriteHeader(200)
	w.Write(data)

}