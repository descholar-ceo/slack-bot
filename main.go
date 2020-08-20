package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"

	"github.com/joho/godotenv"
	"github.com/slack-go/slack"
)

// Res is from internet
type Res map[string]string

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
	user, err := slackClient.GetUserInfo(event.User)

	myCom := retrieveStaticCommands("welcomeMsg")

	fmt.Println("command res" + myCom)

	channelID, timestamp, err := slackClient.PostMessage(
		user.ID,
		slack.MsgOptionText(myCom, false),
		slack.MsgOptionAsUser(true),
	)

	if err != nil {
		fmt.Printf("Ooops! There is an error: %s\n", err)
		return
	}
	fmt.Printf("Message successfully sent to channel %s at %s", channelID, timestamp)
}

// function to retrieve static command from api
func retrieveStaticCommands(command string) string {
	var result Res
	var res string

	resp, err := http.Get(os.Getenv("STATIC_COMMANDS_API"))
	if err != nil {
		fmt.Printf("Ooops! Something went wrong %v\n", err)
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)

	if json.Valid(body) {
		json.Unmarshal(body, &result)
	} else {
		// result["error"]=["there is an error"]
	}

	// iterating over the result
	for k, v := range result {
		if k == command {
			res = v
		} else {
			res = "Sorry! I don't understand that command"
		}
	}

	return res
}
