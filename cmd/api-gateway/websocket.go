package main

import (
	"context"
	"encoding/json"
	"fmt"
	"maps"
	"net/http"
	"sync"
	"time"

	"github.com/abolcerek/Trading-Matching-Simulator/internal/engine"
	"github.com/coder/websocket"
	"github.com/google/uuid"
)

type Registry struct {
	Connections map[*websocket.Conn]uuid.UUID
	mu sync.Mutex
}
func (cfg *apiConfig) HandlerWebSocket(w http.ResponseWriter, r *http.Request) {
	conn, err := websocket.Accept(w, r, &websocket.AcceptOptions{
		OriginPatterns: []string{"localhost:5173"},
	})
	if err != nil {
		fmt.Printf("Error connecting to websocket: %v", err)
		return
	}
	defer conn.CloseNow()
	reply := make(chan engine.Snapshot)
	cfg.requestChannel <- engine.SnapshotRequest{Reply: reply}
	snapshot := <- reply
	data, err := json.Marshal(&snapshot)
	if err != nil {
		fmt.Printf("Error marshalling snapshot: %v", err)
	}
	err = conn.Write(r.Context(), websocket.MessageText, data)
	if err != nil {
		fmt.Printf("Error writing to websocket: %v", err)
		return
	}
	cfg.connectionRegistry.add(conn)
	defer cfg.connectionRegistry.remove(conn)
	for {
		_, _, err := conn.Read(r.Context())
		if err != nil {
			break
		}
	}
}

func (reg *Registry) add(conn *websocket.Conn) {
	reg.mu.Lock()
	id := uuid.New()
	reg.Connections[conn] = id
	reg.mu.Unlock()
}

func (reg *Registry) remove(conn *websocket.Conn) {
	reg.mu.Lock()
	delete(reg.Connections, conn)
	reg.mu.Unlock()
}

func (reg *Registry) broadcast(snapshot []byte) {
	reg.mu.Lock()
	connections := maps.Clone(reg.Connections)
	reg.mu.Unlock()
	for conn := range connections {
		ctx, cancel := context.WithTimeout(context.Background(), time.Second * 5)
		err := conn.Write(ctx, websocket.MessageText, snapshot)
		if err != nil {
			fmt.Printf("Error writing to connection: %v", err)
		}
		cancel()
	}
}

func (cfg *apiConfig) Ws(ch <-chan []engine.Event) {
	for range ch {
		reply := make(chan engine.Snapshot)
		cfg.requestChannel <- engine.SnapshotRequest{Reply: reply}
		snapshot := <- reply
		data, err := json.Marshal(&snapshot)
		if err != nil {
			fmt.Printf("Error marshalling snapshot: %v", err)
		}
		cfg.connectionRegistry.broadcast(data)
	}
}