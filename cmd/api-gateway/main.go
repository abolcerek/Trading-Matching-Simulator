package main

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/abolcerek/Trading-Matching-Simulator/internal/database"
	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
)

type apiConfig struct {
	database *database.Queries
	JWT_secret string
	platform string
}

const balance = 1000


func main() {
	godotenv.Load()
	dbURL := os.Getenv("DB_URL")
	platform := os.Getenv("PLATFORM")
	jwtSecret := os.Getenv("JWT_SECRET")
	port := os.Getenv("PORT")
	db, err := sql.Open("postgres", dbURL)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	err = db.Ping()
	if err != nil {
		log.Fatalf("Error connecting to the database: %v", err)
	}
	mux := http.NewServeMux()
	server := &http.Server{
		Addr: ":" + port,
		Handler: mux,
	}
	ApiCfg := apiConfig{}
	ApiCfg.database = database.New(db)
	ApiCfg.JWT_secret = jwtSecret
	ApiCfg.platform = platform
	mux.HandleFunc("POST /api/users", ApiCfg.HandlerCreateUser)
	mux.HandleFunc("PUT /api/users", ApiCfg.HandlerUpdateUser)
	mux.HandleFunc("POST /api/login", ApiCfg.HandlerLogin)
	mux.HandleFunc("POST /api/orders", ApiCfg.HandlerCreateOrder)
	mux.HandleFunc("DELETE /api/orders/{orderID}", ApiCfg.HandlerCancelOrder)
	log.Fatal(server.ListenAndServe())
}