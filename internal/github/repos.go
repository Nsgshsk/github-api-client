package github

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/Nsgshsk/github-api-client.git/internal/github/dto"
	"github.com/Nsgshsk/github-api-client.git/internal/models"
)

func (c *Client) GetRepos(username string) ([]models.Repository, error) {
	url := fmt.Sprintf("%s/users/%s/repos", c.BaseURL, username)

	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		return nil, err
	}

	req.Header.Set("Accept", "application/vnd.github.v3+json")
	resp, err := c.HTTPClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("unexpected status code: %d", resp.StatusCode)
	}

	var repos []dto.RepoResponse
	if err = json.NewDecoder(resp.Body).Decode(&repos); err != nil {
		return nil, err
	}
	return mapRepos(repos)
}
