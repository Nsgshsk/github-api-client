package github

import (
	"net/http"
)

const baseURL = "https://api.github.com"

type Client struct {
	HTTPClient *http.Client
	BaseURL    string
}

func NewClient() *Client {
	return &Client{
		HTTPClient: http.DefaultClient,
		BaseURL:    baseURL,
	}
}
