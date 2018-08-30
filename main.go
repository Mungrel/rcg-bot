package main

import (
	"fmt"
	"io/ioutil"
	"strings"
	"time"

	"github.com/Mungrel/rcg-bot/bot"
)

const maxRetries = 10

func main() {
	accessToken, err := getAccessToken()
	if err != nil {
		panic(err)
	}

	b := bot.NewBot(accessToken)

	var comic *bot.Comic
	comicErr := bot.ErrRCG500
	retries := 0

	// Retry if we get that specific error
	for comicErr == bot.ErrRCG500 {
		comic, comicErr = b.GetComic()
		if err != nil && err != bot.ErrRCG500 {
			panic(err)
		}

		fmt.Println("failed, retrying...")

		// Lets not ddos them
		time.Sleep(500 * time.Millisecond)
		retries++
		if retries > maxRetries {
			panic("max retries exceeded")
		}
	}

	fmt.Printf("Image URL: %s\nPermalink: %s\n", comic.ComicURL, comic.Permalink)

	err = b.PostToAPI(comic)
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
