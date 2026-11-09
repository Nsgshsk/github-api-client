package github

import (
	"github.com/Nsgshsk/github-api-client/internal/github/dto"
	"github.com/Nsgshsk/github-api-client/internal/models"
)

func (c *Client) GetRepoLanguages(fullName string) (models.Languages, error) {
	var languages dto.LanguagesResponse
	if _, err := c.do("/repos/"+escapeRepo(fullName)+"/languages", &languages); err != nil {
		return nil, err
	}
	return models.Languages(languages), nil
}
