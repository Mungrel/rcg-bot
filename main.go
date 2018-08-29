package main

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"strings"
)

func main() {
	client := http.DefaultClient

	resp, err := client.Get("http://explosm.net/rcg")
	if err != nil {
		panic(err)
	}

	if resp.StatusCode != 200 {
		panic(fmt.Errorf("non-200 status code: %d", resp.StatusCode))
	}

	respBytes, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		panic(err)
	}

	comicURL, permalink := getURL(string(respBytes))
	fmt.Printf("Image URL: %s\nPermalink: %s\n", comicURL, permalink)

	err = postToAPI(comicURL, permalink)
	if err != nil {
		panic(err)
	}

	fmt.Println("Success")
}

func getURL(response string) (string, string) {
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

	return comicURL, permalink
}

const postURL = "https://graph.facebook.com/v3.1/680457985653773/photos"

func postToAPI(comicURL, permalink string) error {
	client := http.DefaultClient

	accessToken, err := getAccessToken()
	if err != nil {
		return err
	}

	params := url.Values{}
	params.Add("url", comicURL)
	params.Add("published", "true")
	params.Add("caption", permalink)

	url := postURL + "?" + params.Encode()

	req, err := http.NewRequest("POST", url, nil)
	if err != nil {
		return err
	}

	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", accessToken))

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

func getAccessToken() (string, error) {
	bytes, err := ioutil.ReadFile("./access_token")
	if err != nil {
		return "", err
	}

	token := string(bytes)
	return strings.TrimSuffix(token, "\n"), nil
}
