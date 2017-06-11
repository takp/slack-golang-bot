package main

import (
	"fmt"

	"github.com/joho/godotenv"
	"github.com/nlopes/slack"
	"log"
	"os"
)

func main() {
	if os.Getenv("GOENV") != "production" {
		err := godotenv.Load()
		if err != nil {
			log.Fatal("Error loading .env file")
		}
	}

	api := slack.New(os.Getenv("SLACK_TOKEN"))
	// If you set debugging, it will log all requests to the console
	// Useful when encountering issues
	api.SetDebug(true)

	rtm := api.NewRTM()
	go rtm.ManageConnection()

	for msg := range rtm.IncomingEvents {
		fmt.Print("Event Received: ")
		switch ev := msg.Data.(type) {
		case *slack.HelloEvent:
			log.Println("Hello!")

		case *slack.ConnectedEvent:
			fmt.Println("Infos:", ev.Info)
			fmt.Println("Connection counter:", ev.ConnectionCount)
			// Replace #general with your Channel ID
			rtm.SendMessage(rtm.NewOutgoingMessage("Hello world", "#general"))

		case *slack.MessageEvent:
			fmt.Printf("Message: %v\n", ev)
			rtm.SendMessage(rtm.NewOutgoingMessage(ev.Text, ev.Channel))

		case *slack.PresenceChangeEvent:
			fmt.Printf("Presence Change: %v\n", ev)

		case *slack.LatencyReport:
			fmt.Printf("Current latency: %v\n", ev.Value)

		case *slack.RTMError:
			fmt.Printf("Error: %s\n", ev.Error())

		case *slack.InvalidAuthEvent:
			fmt.Printf("Invalid credentials")
			return

		default:

			// Ignore other events..
			// fmt.Printf("Unexpected: %v\n", msg.Data)
		}
	}
}
