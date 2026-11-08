package github

import (
	"net/http"
)

const baseURL = "https://api.github.com"

type Client struct {
	httpClient *http.Client
	baseURL    string
}

func NewClient() *Client {
	return &Client{
		httpClient: http.DefaultClient,
		baseURL:    baseURL,
	}
}
