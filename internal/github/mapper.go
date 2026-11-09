package github

import (
	"github.com/Nsgshsk/github-api-client/internal/github/dto"
	"github.com/Nsgshsk/github-api-client/internal/models"
)

func mapUser(resp dto.UserResponse) *models.User {
	return &models.User{
		BaseUser: models.BaseUser{
			Username: resp.Login,
			Homepage: resp.HtmlUrl,
		},
		Email:            resp.Email,
		IsPublic:         resp.UserViewType == "public",
		Name:             resp.Name,
		Company:          resp.Company,
		Blog:             resp.Blog,
		TwitterUsername:  resp.TwitterUsername,
		Location:         resp.Location,
		Hireable:         resp.Hireable,
		Bio:              resp.Bio,
		PublicReposCount: resp.PublicRepos,
		PublicGistsCount: resp.PublicGists,
		FollowersCount:   resp.Followers,
		FollowingCount:   resp.Following,
		CreatedAt:        resp.CreatedAt,
		UpdatedAt:        resp.UpdatedAt,
	}
}

func mapRepo(resp dto.RepoResponse) *models.Repository {
	owner := models.BaseUser{
		Username: resp.Owner.Login,
		Homepage: resp.Owner.HtmlUrl,
	}

	var license *models.License
	if resp.License != nil {
		license = &models.License{
			Name: resp.License.Name,
			URL:  resp.License.Url,
		}
	}

	return &models.Repository{
		BaseRepository: models.BaseRepository{
			Name:            resp.FullName,
			IsPrivate:       resp.Private,
			URL:             resp.HtmlUrl,
			StarsCount:      resp.StargazersCount,
			ForksCount:      resp.ForksCount,
			PrimaryLanguage: languageName(resp.Language),
			Description:     resp.Description,
			UpdatedAt:       resp.UpdatedAt,
			IsFork:          resp.Fork,
			IsArchived:      resp.Archived,
		},
		Owner:           owner,
		MirrorUrl:       resp.MirrorUrl,
		CreatedAt:       resp.CreatedAt,
		PushedAt:        resp.PushedAt,
		Homepage:        resp.Homepage,
		RepoSize:        resp.Size,
		WatchersCount:   resp.WatchersCount,
		IsDisabled:      resp.Disabled,
		OpenIssuesCount: resp.OpenIssuesCount,
		License:         license,
		AllowsForking:   resp.AllowForking,
		IsTemplate:      resp.IsTemplate,
		Topics:          resp.Topics,
		DefaultBranch:   resp.DefaultBranch,
	}
}

func languageName(language *string) string {
	if language == nil {
		return ""
	}
	return *language
}

func mapRepos(resp []dto.RepoResponse) []models.Repository {
	repos := make([]models.Repository, len(resp))
	for i, repo := range resp {
		repos[i] = *mapRepo(repo)
	}
	return repos
}

func mapBaseRepo(resp dto.RepoResponse) *models.BaseRepository {
	return &models.BaseRepository{
		Name:            resp.FullName,
		IsPrivate:       resp.Private,
		URL:             resp.HtmlUrl,
		StarsCount:      resp.StargazersCount,
		ForksCount:      resp.ForksCount,
		PrimaryLanguage: languageName(resp.Language),
		Description:     resp.Description,
		UpdatedAt:       resp.UpdatedAt,
		IsFork:          resp.Fork,
		IsArchived:      resp.Archived,
	}
}

func mapBaseRepos(resp []dto.RepoResponse) []models.BaseRepository {
	repos := make([]models.BaseRepository, len(resp))
	for i, repo := range resp {
		repos[i] = *mapBaseRepo(repo)
	}
	return repos
}
