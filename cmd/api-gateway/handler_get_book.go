package main

import (
	"encoding/json"
	"net/http"

	"github.com/abolcerek/Trading-Matching-Simulator/internal/auth"
	"github.com/abolcerek/Trading-Matching-Simulator/internal/engine"
)
func (cfg *apiConfig) HandlerGetBook(w http.ResponseWriter, r *http.Request) {
	err_params := error_parameters{}
	token, err := auth.GetBearerToken(r.Header)
	if err != nil {
		err_params.Error = "Invalid access token"
		handleErrors(w, &err_params, 401)
		return
	}
	_, err = auth.ValidateJWT(token, cfg.JWT_secret)
	if err != nil {
		err_params.Error = "Invalid access token"
		handleErrors(w, &err_params, 401)
		return
	}
	reply := make(chan engine.Snapshot)
	cfg.requestChannel <- engine.SnapshotRequest{Reply: reply}
	snapshot := <- reply
	data, err := json.Marshal(&snapshot)
	if err != nil {		
		err_params.Error = "Error getting snapshot"
		handleErrors(w, &err_params, 400)
		return
	}
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(200)
	w.Write(data)
}