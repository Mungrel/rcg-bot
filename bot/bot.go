package bot

import (
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"strings"
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

// GetComic gets a Comic's data from explosm.net/rcg
func (bot *Bot) GetComic() (*Comic, error) {
	resp, err := bot.client.Get("http://explosm.net/rcg")
	if err != nil {
		return nil, err
	}

	if resp.StatusCode != 200 {
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

const postURL = "https://graph.facebook.com/v3.1/680457985653773/photos"

// PostToAPI posts a comic to the FB Graph API
func (bot *Bot) PostToAPI(comic *Comic) error {
	client := http.DefaultClient

	params := url.Values{}
	params.Add("url", comic.ComicURL)
	params.Add("caption", comic.Permalink)
	params.Add("published", "true")

	url := postURL + "?" + params.Encode()

	req, err := http.NewRequest("POST", url, nil)
	if err != nil {
		return err
	}

	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", bot.AccessToken))

	resp, err := client.Do(req)
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
