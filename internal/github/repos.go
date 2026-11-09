package github

import (
	"fmt"
	"net/url"

	"github.com/Nsgshsk/github-api-client/internal/github/dto"
	"github.com/Nsgshsk/github-api-client/internal/models"
)

func (c *Client) GetUserRepos(username string) ([]models.Repository, error) {
	path := fmt.Sprintf("/users/%s/repos?per_page=100", url.PathEscape(username))
	repos, err := getPaged[dto.RepoResponse](c, path)
	if err != nil {
		return nil, err
	}
	return mapRepos(repos), nil
}

func (c *Client) GetUserReposShort(username string) ([]models.BaseRepository, error) {
	path := fmt.Sprintf("/users/%s/repos?per_page=100", url.PathEscape(username))
	repos, err := getPaged[dto.RepoResponse](c, path)
	if err != nil {
		return nil, err
	}
	return mapBaseRepos(repos), nil
}
