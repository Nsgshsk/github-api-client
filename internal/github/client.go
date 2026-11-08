package github

import (
	"net/http"
)

const baseURL = "https://api.github.com"

type GithubClient struct {
	HTTPClient *http.Client
	BaseURL    string
}

func NewGithubClient() *GithubClient {
	return &GithubClient{
		HTTPClient: http.DefaultClient,
		BaseURL:    baseURL,
	}
}
