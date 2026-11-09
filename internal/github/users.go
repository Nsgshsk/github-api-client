package github

import (
	"fmt"
	"net/url"

	"github.com/Nsgshsk/github-api-client/internal/github/dto"
	"github.com/Nsgshsk/github-api-client/internal/models"
)

func (c *Client) GetUser(username string) (*models.User, error) {
	var user dto.UserResponse
	if _, err := c.do(fmt.Sprintf("/users/%s", url.PathEscape(username)), &user); err != nil {
		return nil, err
	}
	return mapUser(user), nil
}
