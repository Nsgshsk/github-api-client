package service

import (
	"fmt"

	"github.com/Nsgshsk/github-api-client/internal/github"
	"github.com/Nsgshsk/github-api-client/internal/models"
)

type ReportConfig struct {
	TopRepos        int
	OthersThreshold float64
	ExcludeForks    bool
	ExcludeArchived bool
}

type ReportService struct {
	client *github.Client
	cfg    ReportConfig
}

func NewReportService(client *github.Client, cfg ReportConfig) *ReportService {
	if cfg.TopRepos < 0 {
		cfg.TopRepos = 0
	}
	return &ReportService{
		client: client,
		cfg:    cfg,
	}
}

func (s *ReportService) GetUserReport(username string) (*models.UserReport, error) {
	user, err := s.client.GetUser(username)
	if err != nil {
		return nil, fmt.Errorf("fetch user: %w", err)
	}

	repos, err := s.client.GetUserRepos(username)
	if err != nil {
		return &models.UserReport{
			User:     user,
			Warnings: []string{fetchWarning("fetch repositories", err)},
		}, nil
	}

	filtered := make([]models.Repository, 0, len(repos))
	totalStars := 0
	for _, repo := range repos {
		if s.cfg.ExcludeForks && repo.IsFork {
			continue
		}
		if s.cfg.ExcludeArchived && repo.IsArchived {
			continue
		}
		filtered = append(filtered, repo)
		totalStars += repo.StarsCount
	}
	sortRepositories(filtered)

	totals := make(map[string]int64)
	warnings := make([]string, 0)
	for _, repo := range filtered {
		languages, err := s.client.GetRepoLanguages(repo.Name)
		if err != nil {
			warnings = append(warnings, fetchWarning("fetch languages for "+repo.Name, err))
			continue
		}
		for name, bytes := range languages {
			totals[name] += bytes
		}
	}

	return &models.UserReport{
		User:          user,
		TopRepos:      topRepositories(filtered, s.cfg.TopRepos),
		Languages:     aggregateLanguages(totals, s.cfg.OthersThreshold),
		ReposAnalyzed: len(filtered),
		TotalStars:    totalStars,
		Warnings:      warnings,
	}, nil
}

func fetchWarning(operation string, err error) string {
	if github.IsRateLimited(err) {
		return fmt.Sprintf("%s: %v (hint: set GITHUB_TOKEN or use --auth-mode=authenticated)", operation, err)
	}
	return fmt.Sprintf("%s: %v", operation, err)
}
