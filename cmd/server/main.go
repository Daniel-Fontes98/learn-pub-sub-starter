package main

import (
	"fmt"
	"log"

	"github.com/bootdotdev/learn-pub-sub-starter/internal/gamelogic"
	"github.com/bootdotdev/learn-pub-sub-starter/internal/pubsub"
	"github.com/bootdotdev/learn-pub-sub-starter/internal/routing"
	amqp "github.com/rabbitmq/amqp091-go"
)

func main() {
	fmt.Println("Starting Peril server...")
	connectionSring := "amqp://guest:guest@localhost:5672/"
	conn, err := amqp.Dial(connectionSring)
	if err != nil {
		log.Fatalf("Couldnt dial into: %v\nError: %v", connectionSring, err)
	}
	defer conn.Close()
	fmt.Println("Connection sucessuful!")

	ch, _, err := pubsub.DeclareAndBind(
		conn, 
		routing.ExchangePerilDirect,
		"game_logs", 
		routing.GameLogSlug+".*", 
		pubsub.SimpleQueueDurable,
	)
	if err != nil {
		log.Fatalf("Failed to declare and Bind to channel: %v", err)
	}

	gamelogic.PrintServerHelp()
	for {
		userInput := gamelogic.GetInput()
		if len(userInput) == 0 {
			continue
		}

		switch userInput[0] {
		case "pause":
			fmt.Println("Sending a pause message")
			err := pubsub.PublishJSON(ch, routing.ExchangePerilDirect, routing.PauseKey, routing.PlayingState{ IsPaused: true})
			if err != nil {
				log.Printf("Could not publish: %v", err)
			}
			
		case "resume":
			fmt.Println("Sending a resume message")
			err := pubsub.PublishJSON(ch, routing.ExchangePerilDirect, routing.PauseKey, routing.PlayingState{ IsPaused: false})
			if err != nil {
				log.Printf("Could not publish: %v", err)
			}
		
		case "quit":
			fmt.Println("Server shutting down")
			return
		default:
			fmt.Println("Invalid command")
		}
	}
}
