package main

import (
	"fmt"
	"io/ioutil"
	"strings"

	"github.com/Mungrel/rcg-bot/bot"
)

func main() {
	accessToken, err := getAccessToken()
	if err != nil {
		panic(err)
	}

	bot := bot.NewBot(accessToken)

	comic, err := bot.GetComic()
	if err != nil {
		panic(err)
	}

	fmt.Printf("Image URL: %s\nPermalink: %s\n", comic.ComicURL, comic.Permalink)

	err = bot.PostToAPI(comic)
	if err != nil {
		panic(err)
	}

	fmt.Println("Success")
}

func getAccessToken() (string, error) {
	bytes, err := ioutil.ReadFile("./access_token")
	if err != nil {
		return "", err
	}

	token := string(bytes)
	return strings.TrimSuffix(token, "\n"), nil
}
