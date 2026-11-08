package service

import (
	"cmp"
	"slices"

	"github.com/Nsgshsk/github-api-client.git/internal/github"
	"github.com/Nsgshsk/github-api-client.git/internal/models"
)

type UserService struct {
	client *github.Client
}

func NewUserService() *UserService {
	return &UserService{
		client: github.NewClient(),
	}
}

func (s *UserService) GetUserOverview(username string) (*models.UserOverview, error) {
	user, err := s.client.GetUser(username)
	if err != nil {
		return nil, err
	}

	repos, err := s.client.GetUserReposShort(username)
	if err != nil {
		return nil, err
	}
	slices.SortFunc(repos, func(a, b models.BaseRepository) int {
		if a.StarsCount != b.StarsCount {
			return cmp.Compare(b.StarsCount, a.StarsCount)
		}
		return cmp.Compare(a.Name, b.Name)
	})

	overview := models.UserOverview{
		User:      user,
		Top5Repos: repos[:5],
	}
	return &overview, nil
}
