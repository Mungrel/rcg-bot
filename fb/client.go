package fb

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
)

const baseURL = "https://graph.facebook.com/v3.1/680457985653773"

// Client represents a client for interacting with the FB API.
type Client struct {
	accessToken string
	client      *http.Client
}

// NewClient creates a new Client.
func NewClient(accessToken string) *Client {
	return &Client{
		accessToken: accessToken,
		client:      http.DefaultClient,
	}
}

// Post makes an auth'd post request to the specified relative URL with the specified query params.
func (fb *Client) Post(relativeURL string, params url.Values) error {
	encodedURL := baseURL + relativeURL + "?" + params.Encode()
	return fb.doRequest(http.MethodPost, encodedURL)
}

func (fb *Client) doRequest(method, encodedURL string) error {
	req, err := http.NewRequest("POST", encodedURL, nil)
	if err != nil {
		return err
	}

	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", fb.accessToken))

	resp, err := fb.client.Do(req)
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
