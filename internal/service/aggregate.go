package service

import (
	"cmp"
	"slices"
	"strings"

	"github.com/Nsgshsk/github-api-client/internal/models"
)

const othersLabel = "Others"

func sortRepositories(repos []models.Repository) {
	slices.SortFunc(repos, func(a, b models.Repository) int {
		if a.StarsCount != b.StarsCount {
			return cmp.Compare(b.StarsCount, a.StarsCount)
		}
		return cmp.Compare(a.Name, b.Name)
	})
}

func topRepositories(repos []models.Repository, n int) []models.BaseRepository {
	if n <= 0 {
		return nil
	}
	if len(repos) > n {
		repos = repos[:n]
	}

	result := make([]models.BaseRepository, len(repos))
	for i, repo := range repos {
		result[i] = repo.BaseRepository
	}
	return result
}

func aggregateLanguages(totals map[string]int64, threshold float64) models.LanguageDistribution {
	result := make(models.LanguageDistribution, 0, len(totals))

	var total int64
	for _, bytes := range totals {
		total += bytes
	}
	if total <= 0 {
		return result
	}

	type share models.LanguageShare

	shares := make([]share, 0, len(totals))
	for name, bytes := range totals {
		shares = append(shares, share{
			Name:       name,
			Percentage: float64(bytes) / float64(total) * 100,
		})
	}
	slices.SortFunc(shares, func(a, b share) int {
		if a.Percentage != b.Percentage {
			return cmp.Compare(b.Percentage, a.Percentage)
		}
		return cmp.Compare(a.Name, b.Name)
	})

	grouping := threshold > 0
	var others float64
	kept := false
	for _, item := range shares {
		if strings.EqualFold(item.Name, othersLabel) {
			others += item.Percentage
			continue
		}

		keep := !grouping || !kept || item.Percentage >= threshold
		if !keep {
			others += item.Percentage
			continue
		}
		result = append(result, models.LanguageShare(item))
		kept = true
	}

	if others > 0 {
		result = append(result, models.LanguageShare{
			Name:       othersLabel,
			Percentage: others,
		})
	}

	return result
}
