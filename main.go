package main

import (
	"io/ioutil"
	"strings"
	"time"

	"github.com/Mungrel/rcg-bot/bot"
)

var postDelay = 30 * time.Minute
var rcgBot *bot.Bot

func main() {
	accessToken, err := getAccessToken()
	if err != nil {
		panic(err)
	}

	rcgBot = bot.NewBot(accessToken)

	// infinite loop with a 30 minute sleep/delay
	for {

		err = rcgBot.Post()
		if err != nil {
			panic(err)
		}
		time.Sleep(postDelay)
	}
}

func getAccessToken() (string, error) {
	bytes, err := ioutil.ReadFile("./access_token")
	if err != nil {
		return "", err
	}

	token := string(bytes)
	return strings.TrimSuffix(token, "\n"), nil
}
