package main

import (
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
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
