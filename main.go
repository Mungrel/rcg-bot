package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"strings"
	"time"

	"github.com/Mungrel/rcg-bot/bot"
)

var postDelay = 30 * time.Minute

const (
	modePost  = "post"
	modeTop10 = "top10"
)

func main() {
	accessToken, err := getAccessToken()
	if err != nil {
		panic(err)
	}

	rcgBot := bot.NewBot(accessToken)

	mode := flag.String("mode", modePost, "Mode to run bot in. (post)")
	flag.Parse()

	if mode != nil && *mode == modeTop10 {
		fmt.Println("Running in Top10 mode")
		err = rcgBot.Top10()
		if err != nil {
			panic(err)
		}

		return
	}

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
