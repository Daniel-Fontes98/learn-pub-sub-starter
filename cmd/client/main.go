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
	fmt.Println("Starting Peril client...")
	connectionSring := "amqp://guest:guest@localhost:5672/"
	conn, err := amqp.Dial(connectionSring)
	if err != nil {
		log.Fatalf("Couldnt dial into: %v\nError: %v", connectionSring, err)
	}
	defer conn.Close()
	fmt.Println("Connection sucessuful!")

	username, err := gamelogic.ClientWelcome()
	if err != nil {
		log.Fatalf("Couldn't get username")
	}

	pubsub.DeclareAndBind(
		conn,
		routing.ExchangePerilDirect,
		routing.PauseKey+"."+username,
		routing.PauseKey, 
		pubsub.SimpleQueueTransient,
	)
	state := gamelogic.NewGameState(username)

	gamelogic.PrintClientHelp()
	for {
		userInput := gamelogic.GetInput()
		if len(userInput) == 0 {
			continue
		}

		switch userInput[0] {
		case "spawn":
			err := state.CommandSpawn(userInput)
			if err != nil {
				log.Printf("Couldn't spawn a unit: %v", err)
			}

		case "move":
			_, err := state.CommandMove(userInput)
			if err != nil {
				log.Printf("Couldn't move: %v", err)
			}

		case "status":
			state.CommandStatus()

		case "help":
			gamelogic.PrintClientHelp()

		case "spam":
			fmt.Println("Spamming not allowed yet!")

		case "quit":
			gamelogic.PrintQuit()
			return

		default:
			fmt.Println("Invalid commad")
		}
	}
}
