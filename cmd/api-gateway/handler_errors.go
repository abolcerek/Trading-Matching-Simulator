package main

import (
	"encoding/json"
	"net/http"
	"log"
)

type error_parameters struct {
	Error string `json:"error"`
}

func handleErrors(w http.ResponseWriter, err_params *error_parameters, status_code int) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(status_code)
	data, err := json.Marshal(err_params)
	if err != nil {
		log.Printf("Error marshalling JSON")
		return
	}
	w.Write(data) 
}