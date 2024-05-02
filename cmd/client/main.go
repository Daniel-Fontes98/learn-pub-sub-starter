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
		log.Fatalf("Couldn't get username: %v", err)
	}

	gs := gamelogic.NewGameState(username)
	ch, _, err := pubsub.DeclareAndBind(
		conn,
		routing.ExchangePerilTopic,
		routing.ArmyMovesPrefix+"."+gs.GetUsername(),
		routing.ArmyMovesPrefix+".*",
		pubsub.SimpleQueueTransient,
	)
	if err != nil {
		log.Fatalf("Couldn't bind to army move channel: %v", err)
	}

	err = pubsub.SubscribeJSON(
		conn, 
		routing.ExchangePerilDirect,
		routing.PauseKey+"."+gs.GetUsername(),
		routing.PauseKey,
		pubsub.SimpleQueueTransient,
		handlerPause(gs),
	)
	if err != nil {
		log.Fatalf("Couldn't subscribe to pause queue: %v", err)
	}

	err = pubsub.SubscribeJSON(
		conn,
		routing.ExchangePerilTopic,
		routing.ArmyMovesPrefix+"."+gs.GetUsername(),
		routing.ArmyMovesPrefix,
		pubsub.SimpleQueueTransient,
		handlerMove(gs),
	)
	if err != nil {
		log.Fatalf("Couldn't subscribe to pause queue: %v", err)
	}

	gamelogic.PrintClientHelp()
	for {
		userInput := gamelogic.GetInput()
		if len(userInput) == 0 {
			continue
		}

		switch userInput[0] {
		case "spawn":
			err := gs.CommandSpawn(userInput)
			if err != nil {
				log.Printf("Couldn't spawn a unit: %v", err)
			}

		case "move":
			am, err := gs.CommandMove(userInput)
			if err != nil {
				log.Printf("Couldn't move: %v", err)
			}

			err = pubsub.PublishJSON(ch, routing.ExchangePerilTopic, routing.ArmyMovesPrefix, am)
			if err != nil {
				log.Printf("Couldn't publish move: %v", err)
			}

		case "status":
			gs.CommandStatus()

		case "help":
			gamelogic.PrintClientHelp()

		case "spam":
			fmt.Println("Spamming not allowed yet!")

		case "quit":
			gamelogic.PrintQuit()
			return

		default:
			fmt.Println("Invalid command")
		}
	}
}


