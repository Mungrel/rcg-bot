package bot

import (
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"strings"
	"time"

	"github.com/Mungrel/rcg-bot/fb"
)

// Bot is a convenience struct for calling bot-related methods
type Bot struct {
	fbClient *fb.Client
	client   *http.Client
}

// NewBot creates a new Bot object
func NewBot(accessToken string) *Bot {
	return &Bot{
		fbClient: fb.NewClient(accessToken),
		client:   http.DefaultClient,
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

	err = bot.postComic(comic)
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

	const comicURLTagPrefix = `<a href="https://rcg-cdn.explosm.net/comics`
	const comicURLTagSuffix = `.png" class="custom-social-button" title="Download Image" download target="_blank">`
	for _, line := range lines {
		if strings.HasPrefix(line, comicURLTagPrefix) && strings.HasSuffix(line, comicURLTagSuffix) {
			comicURLTag = line
		} else if strings.Contains(line, "<input id=\"permalink\"") {
			permalinkTag = line
		}
	}

	src := strings.Split(comicURLTag, " ")[1]
	url := strings.Split(src, "=")[1]

	comicURL := strings.Trim(url, `"`)

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

const photosURL = "/photos"

func (bot *Bot) postComic(comic *Comic) error {
	params := url.Values{}
	params.Add("url", comic.ComicURL)
	params.Add("caption", comic.Permalink)
	params.Add("published", "true")

	return bot.fbClient.Post(photosURL, params)
}
