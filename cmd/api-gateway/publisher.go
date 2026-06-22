package main

import (
	"fmt"

	"github.com/abolcerek/Trading-Matching-Simulator/internal/engine"
	"github.com/abolcerek/Trading-Matching-Simulator/internal/queue"
)



func (cfg *apiConfig) Publisher(events <-chan []engine.Event) {
	ch, err := cfg.rabbitmqConnection.Channel()
	if err != nil {
		fmt.Printf("Error creating channel for rabbitmq: %v", err)
		return
	}
	for event := range events {
		err = queue.PublishJSON(ch, "events", "", event)
		if err != nil {
			fmt.Printf("Error publishing to exchange: %v", err)
			return
		}
	}
}