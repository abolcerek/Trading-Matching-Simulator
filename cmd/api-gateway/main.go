package main

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"os"
	"sync"

	"github.com/abolcerek/Trading-Matching-Simulator/internal/database"
	"github.com/abolcerek/Trading-Matching-Simulator/internal/engine"
	"github.com/abolcerek/Trading-Matching-Simulator/internal/engine/types"
	"github.com/coder/websocket"
	"github.com/google/uuid"
	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
)

type apiConfig struct {
	database *database.Queries
	db *sql.DB
	JWT_secret string
	platform string
	orderbook *engine.OrderBook
	orderChannel chan types.Envelope
	writerChannel chan []engine.Event
	wsChannel chan []engine.Event
	requestChannel chan engine.SnapshotRequest
	ConnectionRegistry *Registry
}

const balance = 1000
const place = "place"
const cancel = "cancel"

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
	ApiCfg.db = db
	ApiCfg.JWT_secret = jwtSecret
	ApiCfg.platform = platform
	orderbook := engine.NewOrderBook()
	ApiCfg.orderbook = orderbook
	ApiCfg.orderChannel = make(chan types.Envelope, 100)
	ApiCfg.writerChannel = make(chan []engine.Event, 100)
	ApiCfg.wsChannel = make(chan []engine.Event, 100)
	ApiCfg.requestChannel = make(chan engine.SnapshotRequest)
	ApiCfg.ConnectionRegistry = &Registry{
		Connections: map[*websocket.Conn]uuid.UUID{},
		mu: sync.Mutex{},
	}
	err = ApiCfg.Replay()
	if err != nil {
		log.Fatalf("Error replaying orders from the database: %v", err)
	}
	go engine.RunEngine(ApiCfg.orderbook, ApiCfg.orderChannel, ApiCfg.writerChannel, ApiCfg.wsChannel, ApiCfg.requestChannel)
	go ApiCfg.Consumer(ApiCfg.writerChannel)
	go ApiCfg.Ws(ApiCfg.wsChannel)
	mux.HandleFunc("GET /api/ws", ApiCfg.HandlerWebSocket)
	mux.HandleFunc("POST /api/users", ApiCfg.HandlerCreateUser)
	mux.HandleFunc("PUT /api/users", ApiCfg.HandlerUpdateUser)
	mux.HandleFunc("POST /api/login", ApiCfg.HandlerLogin)
	mux.HandleFunc("POST /api/orders", ApiCfg.HandlerCreateOrder)
	mux.HandleFunc("GET /api/orders/{orderID}", ApiCfg.HandlerGetOrder)
	mux.HandleFunc("DELETE /api/orders/{orderID}", ApiCfg.HandlerCancelOrder)
	mux.HandleFunc("GET /api/book", ApiCfg.HandlerGetBook)
	log.Fatal(server.ListenAndServe())
}
