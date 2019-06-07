package fb

import (
	"encoding/json"
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
	return fb.doRequest(http.MethodPost, encodedURL, nil)
}

// GetAbsoluteURL makes an auth'd get request to the specified absolute URL and
// marshals the result into the provided entities param.
func (fb *Client) GetAbsoluteURL(absoluteURL string, entities interface{}) error {
	return fb.doRequest(http.MethodGet, absoluteURL, entities)
}

// Get makes an auth'd get request to the specified relative URL and
// marshals the result into the provided entities param.
func (fb *Client) Get(relativeURL string, params url.Values, entities interface{}) error {
	encodedURL := baseURL + relativeURL + "?" + params.Encode()
	return fb.doRequest(http.MethodGet, encodedURL, entities)
}

func (fb *Client) doRequest(method, encodedURL string, entities interface{}) error {
	req, err := http.NewRequest(method, encodedURL, nil)
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

	if entities != nil {
		return json.Unmarshal(respBody, entities)
	}

	return nil
}
