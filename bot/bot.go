package bot

import (
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"strings"
	"time"
)

// Bot is a convenience struct for calling bot-related methods
type Bot struct {
	AccessToken string
	client      *http.Client
}

// NewBot creates a new Bot object
func NewBot(accessToken string) *Bot {
	return &Bot{
		AccessToken: accessToken,
		client:      http.DefaultClient,
	}
}

// Comic is a wrapper for its related data
type Comic struct {
	ComicURL  string
	Permalink string
}

var errRCG500 = errors.New("non-200 status code")

const maxRetries = 10

// doPost gets a comic and posts it to the page
func (bot *Bot) doPost() error {
	var comic *Comic

	comic, err := bot.getComic()
	if err != nil {
		return err
	}

	fmt.Printf("Image URL: %s\nPermalink: %s\nTime: %s\n", comic.ComicURL, comic.Permalink, time.Now().Format(time.RFC3339))

	err = bot.postToAPI(comic)
	if err != nil {
		return err
	}

	fmt.Println("Success")
	return nil
}

func (bot *Bot) Post() error {
	return repeat(0, bot.doPost)
}

func repeat(retries int, f func() error) error {
	if retries > maxRetries {
		return errors.New("max retries exceeded")
	}
	err := f()
	if err != nil {
		fmt.Printf("failed: %s\nretrying...\n", err.Error())
		time.Sleep(500 * time.Millisecond)
		retries++
		repeat(retries, f)
	}

	return nil
}

// GetComic gets a Comic's data from explosm.net/rcg
func (bot *Bot) getComic() (*Comic, error) {
	req, err := http.NewRequest("GET", "http://explosm.net/rcg?promo=false", nil)
	if err != nil {
		return nil, err
	}

	req.Header.Set("User-Agent", "cyanide_bot_69")

	resp, err := bot.client.Do(req)
	if err != nil {
		return nil, err
	}

	if resp.StatusCode == 500 {
		return nil, errRCG500
	} else if resp.StatusCode != 200 {
		return nil, fmt.Errorf("non-200 status code: %d", resp.StatusCode)
	}

	respBytes, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	response := string(respBytes)

	lines := strings.Split(response, "\n")
	var comicURLTag string
	var permalinkTag string
	for _, line := range lines {
		if strings.Contains(line, "<img src=\"//files.explosm.net/rcg/") {
			comicURLTag = line
		} else if strings.Contains(line, "<input id=\"permalink\"") {
			permalinkTag = line
		}
	}

	src := strings.Split(comicURLTag, " ")[1]
	url := strings.Split(src, "=")[1]

	comicURL := "http:" + strings.Trim(url, `"`)

	input := strings.Split(permalinkTag, " ")[3]
	value := strings.Split(input, "=")[1]

	permalink := strings.Trim(value, `"`)

	if comicURL == "" {
		return nil, errors.New("failed to parse comic url")
	}

	if permalink == "" {
		return nil, errors.New("failed to parse permalink")
	}

	return &Comic{
		ComicURL:  comicURL,
		Permalink: permalink,
	}, nil
}

const (
	baseURL = "https://graph.facebook.com/v3.1/680457985653773"
	photos  = "/photos"
)

func (bot *Bot) post(relativeURL string, queryParams map[string]string) error {
	params := url.Values{}
	for key, value := range queryParams {
		params.Add(key, value)
	}

	encodedURL := baseURL + relativeURL + "?" + params.Encode()

	req, err := http.NewRequest("POST", encodedURL, nil)
	if err != nil {
		return err
	}

	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", bot.AccessToken))

	resp, err := bot.client.Do(req)
	if err != nil {
		return err
	}

	respBody, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return err
	}

	if resp.StatusCode > 299 {
		return fmt.Errorf("\nbad response %d\nresponse: %s", resp.StatusCode, string(respBody))
	}

	return nil
}

func (bot *Bot) postToAPI(comic *Comic) error {
	params := map[string]string{
		"url":       comic.ComicURL,
		"caption":   comic.Permalink,
		"published": "true",
	}

	return bot.post(photos, params)
}
