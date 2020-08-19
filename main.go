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
		// slack.OptionDebug(true),
		// slack.OptionLog(log.New(os.Stdout, "slack bot:", log.Lshortfile|log.LstdFlags)),
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
	attachment := slack.Attachment{
		Pretext: "Hello descholar",
		Text:    "I am happy to see you here!",
		// Uncomment the following part to send a field too

		Fields: []slack.AttachmentField{
			slack.AttachmentField{
				Title: "a",
				Value: "no",
			},
		},
	}

	// fmt.Printf("This is the event obj %v: \n", event.Username)
	channelID, timestamp, err := slackClient.PostMessage(
		event.User,
		slack.MsgOptionText("Hello there!", true),
		slack.MsgOptionAttachments(attachment),
	)

	if err != nil {
		fmt.Printf("Ooops! There is an error: %s\n", err)
		return
	}
	fmt.Printf("Message successfully sent to channel %s at %s", channelID, timestamp)
}
