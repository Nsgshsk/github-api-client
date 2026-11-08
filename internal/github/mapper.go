package github

import (
	"github.com/Nsgshsk/github-api-client.git/internal/github/dto"
	"github.com/Nsgshsk/github-api-client.git/internal/models"
)

func mapUser(resp dto.UserResponse) (*models.User, error) {
	isPublic := false
	if resp.UserViewType == "public" {
		isPublic = true
	}
	return &models.User{
		Username:         resp.Login,
		Homepage:         resp.HtmlUrl,
		Email:            resp.Email,
		IsPublic:         isPublic,
		Name:             resp.Name,
		Location:         resp.Location,
		Hirable:          resp.Hireable,
		Bio:              resp.Bio,
		PublicReposCount: resp.PublicRepos,
		PublicGistsCount: resp.PublicGists,
		FollowersCount:   resp.Followers,
		FollowingCount:   resp.Following,
		CreatedAt:        resp.CreatedAt,
		UpdatedAt:        resp.UpdatedAt,
	}, nil
}

func mapRepo(resp dto.RepoResponse) (*models.Repository, error) {
	owner := models.BaseUser{
		Username: resp.Owner.Login,
		Homepage: resp.Owner.HtmlUrl,
	}
	var license models.License
	if resp.License != nil {
		license = models.License{
			Name: resp.Name,
			Url:  resp.License.Url,
		}
	}

	return &models.Repository{
		Name:            resp.FullName,
		IsPrivate:       resp.Private,
		Owner:           owner,
		URL:             resp.HtmlUrl,
		MirrorUrl:       resp.MirrorUrl,
		IsFork:          resp.Fork,
		CreatedAt:       resp.CreatedAt,
		UpdatedAt:       resp.UpdatedAt,
		PushedAt:        resp.PushedAt,
		Homepage:        resp.Homepage,
		RepoSize:        resp.Size,
		StarsCount:      resp.StargazersCount,
		WatchersCount:   resp.WatchersCount,
		ForksCount:      resp.Forks,
		IsArchived:      resp.Archived,
		IsDisabled:      resp.Archived,
		OpenIssuesCount: resp.OpenIssuesCount,
		License:         &license,
		AllowsForking:   resp.AllowForking,
		IsTemplate:      resp.IsTemplate,
		Topics:          resp.Topics,
		DefaultBranch:   resp.DefaultBranch,
	}, nil
}

func mapRepos(resp []dto.RepoResponse) ([]models.Repository, error) {
	repos := make([]models.Repository, len(resp))
	for i, repo := range resp {
		repo, _ := mapRepo(repo)
		repos[i] = *repo
	}
	return repos, nil
}
