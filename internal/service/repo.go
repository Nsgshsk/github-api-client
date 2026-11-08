package service

import (
	"cmp"
	"log"
	"slices"

	"github.com/Nsgshsk/github-api-client.git/internal/github"
	"github.com/Nsgshsk/github-api-client.git/internal/models"
)

const maxRepos = 10

type RepoService struct {
	client *github.Client
}

func NewRepoService() *RepoService {
	return &RepoService{
		client: github.NewClient(),
	}
}

func (s *RepoService) GetUserReposOverview(username string) (*models.UserReposOverview, error) {
	repos, err := s.client.GetUserRepos(username)
	if err != nil {
		return nil, err
	}
	slices.SortFunc(repos, func(a, b models.Repository) int {
		if a.StarsCount != b.StarsCount {
			return cmp.Compare(b.StarsCount, a.StarsCount)
		}
		return cmp.Compare(a.Name, b.Name)
	})
	if len(repos) > maxRepos {
		repos = repos[:maxRepos]
		log.Printf("There are more than %d repos, limiting to %d most stared repos", maxRepos, maxRepos)
	}

	langPercentage := make([]models.LangPercentage, len(repos))
	for i, repo := range repos {
		if repo.IsPrivate {
			log.Printf("%s is private repo", repo.Name)
			continue
		}
		langs, err := s.client.GetRepoLanguages(repo.Name)
		if err != nil {
			log.Println(err)
			langPercentage[i] = nil
			continue
		}
		langPercentage[i] = mapLangPercentage(&langs)
	}
	return &models.UserReposOverview{
		Repositories: repos,
		Languages:    langPercentage,
	}, nil
}
