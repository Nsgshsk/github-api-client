package github

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/Nsgshsk/github-api-client.git/internal/github/dto"
	"github.com/Nsgshsk/github-api-client.git/internal/models"
)

func (c *Client) GetRepoLanguages(repo string) (models.Languages, error) {
	url := fmt.Sprintf("%s/repos/%s/languages", c.baseURL, repo)

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

	var languages dto.LanguagesResponse
	if err = json.NewDecoder(resp.Body).Decode(&languages); err != nil {
		return nil, err
	}
	return models.Languages(languages), nil
}
