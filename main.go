package main

import (
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"net/url"
	"os"
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

	imgURL := getURL(string(respBytes))
	fmt.Printf("Image URL: %s\n", imgURL)

	err = saveImage(imgURL, "tmp.png")
	if err != nil {
		panic(err)
	}

	fmt.Println("Image saved.")

	err = postToAPI(imgURL)
	if err != nil {
		panic(err)
	}
}

func saveImage(url, fileName string) error {
	response, err := http.Get(url)
	if err != nil {
		return err
	}

	defer response.Body.Close()

	file, err := os.Create(fileName)
	if err != nil {
		return err
	}

	_, err = io.Copy(file, response.Body)
	if err != nil {
		return err
	}

	return file.Close()
}

func getURL(response string) string {
	lines := strings.Split(response, "\n")
	var tag string
	for _, line := range lines {
		if strings.Contains(line, "<img src=\"//files.explosm.net/rcg/") {
			tag = line
			break
		}
	}

	src := strings.Split(tag, " ")[1]
	url := strings.Split(src, "=")[1]

	trimmedURL := strings.Trim(url, `"`)
	return "http:" + trimmedURL
}

const pageID = "680457985653773"
const postURL = "https://graph.facebook.com/v3.1/680457985653773/photos"

type PostBody struct {
	URL       string `json:"url"`
	Published string `json:"published"`
}

func postToAPI(comicURL string) error {
	client := http.DefaultClient

	accessToken, err := getAccessToken()
	if err != nil {
		return err
	}

	params := url.Values{}
	params.Add("url", comicURL)
	params.Add("published", "true")

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
