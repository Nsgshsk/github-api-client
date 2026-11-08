package github

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/Nsgshsk/github-api-client.git/internal/github/dto"
	"github.com/Nsgshsk/github-api-client.git/internal/models"
)

func (c *Client) GetUser(username string) (*models.User, error) {
	url := fmt.Sprintf("%s/users/%s", c.baseURL, username)

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

	var user dto.UserResponse
	if err = json.NewDecoder(resp.Body).Decode(&user); err != nil {
		return nil, err
	}
	return mapUser(user)
}
