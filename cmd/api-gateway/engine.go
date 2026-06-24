package main

import (
	"fmt"
	"log"

	"github.com/abolcerek/Trading-Matching-Simulator/internal/engine/types"
	"github.com/abolcerek/Trading-Matching-Simulator/internal/queue"
)

func (cfg *apiConfig) EngineConsumer() {
	ch, err := cfg.rabbitmqConnection.Channel()
	if err != nil {
		fmt.Printf("Error creating channel: %v", err)
		return
	}
	err = queue.SubscribeJSON(ch, "orders", "engine", "", cfg.handleOrdersEngine)
	if err != nil {
		log.Printf("Error subscribing to queue: %v", err)
	}
}

func (cfg *apiConfig) handleOrdersEngine(envelope types.Envelope) string {
	cfg.orderChannel <- envelope
	return queue.Ack
}