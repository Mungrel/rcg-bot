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

// Post gets a comic and posts it to the page
func (bot *Bot) Post() error {
	var comic *Comic
	comicErr := errRCG500
	retries := 0

	// Retry if we get that specific error
	for comicErr == errRCG500 {
		comic, comicErr = bot.getComic()
		if comicErr != nil && comicErr != errRCG500 {
			return comicErr
		}

		fmt.Println("failed, retrying...")

		// Lets not ddos them
		time.Sleep(500 * time.Millisecond)
		retries++
		if retries > maxRetries {
			return errors.New("max retries exceeded")
		}
	}

	fmt.Printf("Image URL: %s\nPermalink: %s\n", comic.ComicURL, comic.Permalink)

	err := bot.postToAPI(comic)
	if err != nil {
		return err
	}

	fmt.Println("Success")
	return nil
}

// GetComic gets a Comic's data from explosm.net/rcg
func (bot *Bot) getComic() (*Comic, error) {
	resp, err := bot.client.Get("http://explosm.net/rcg")
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

const postURL = "https://graph.facebook.com/v3.1/680457985653773/photos"

func (bot *Bot) postToAPI(comic *Comic) error {
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
