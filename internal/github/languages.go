package github

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/Nsgshsk/github-api-client.git/internal/models"
)

func (c *Client) GetLanguages(repo string) (*models.Languages, error) {
	url := fmt.Sprintf("%s/repos/%s/languages", c.BaseURL, repo)

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

	var languages models.Languages
	if err = json.NewDecoder(resp.Body).Decode(&languages); err != nil {
		return nil, err
	}
	return &languages, nil
}
