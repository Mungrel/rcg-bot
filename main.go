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

	err = rcgBot.Post()
	if err != nil {
		panic(err)
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
