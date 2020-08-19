package main

import (
	"fmt"
	"log"
	"os"

	"github.com/joho/godotenv"
	"github.com/slack-go/slack"
)

var slackClient *slack.Client

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	slackAccessToken := os.Getenv("SLACK_ACCESS_TOKEN")

	slackClient = slack.New(
		slackAccessToken,
		slack.OptionDebug(true),
		slack.OptionLog(log.New(os.Stdout, "slack bot:", log.Lshortfile|log.LstdFlags)),
	)

	rtm := slackClient.NewRTM()
	go rtm.ManageConnection()

	for msg := range rtm.IncomingEvents {
		switch ev := msg.Data.(type) {
		case *slack.MessageEvent:
			go handleMsgFromSlack(ev)
		}
	}
}

func handleMsgFromSlack(event *slack.MessageEvent) {
	fmt.Printf("%v\n", event)
}
