package github

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/Nsgshsk/github-api-client.git/internal/github/dto"
	"github.com/Nsgshsk/github-api-client.git/internal/models"
)

func (c *Client) GetUserRepos(username string) ([]models.Repository, error) {
	url := fmt.Sprintf("%s/users/%s/repos", c.baseURL, username)

	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		return nil, err
	}

	req.Header.Set("Accept", "application/vnd.github+json")
	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer func() {
		err := resp.Body.Close()
		if err != nil {
			log.Printf("failed to close response body: %s", err)
		}
	}()
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("unexpected status code: %d", resp.StatusCode)
	}

	var repos []dto.RepoResponse
	if err = json.NewDecoder(resp.Body).Decode(&repos); err != nil {
		return nil, err
	}
	return mapRepos(repos)
}

func (c *Client) GetUserReposShort(username string) ([]models.BaseRepository, error) {
	url := fmt.Sprintf("%s/users/%s/repos", c.baseURL, username)

	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		return nil, err
	}

	req.Header.Set("Accept", "application/vnd.github+json")
	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer func() {
		err := resp.Body.Close()
		if err != nil {
			log.Printf("failed to close response body: %s", err)
		}
	}()
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("unexpected status code: %d", resp.StatusCode)
	}

	var repos []dto.RepoResponse
	if err = json.NewDecoder(resp.Body).Decode(&repos); err != nil {
		return nil, err
	}
	return mapBaseRepos(repos)
}
