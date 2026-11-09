package models

import "time"

type BaseRepository struct {
	Name            string
	IsPrivate       bool
	URL             string
	StarsCount      int
	ForksCount      int
	PrimaryLanguage string
	Description     *string
	UpdatedAt       time.Time
	IsFork          bool
	IsArchived      bool
}

type Repository struct {
	BaseRepository
	Owner           BaseUser
	MirrorUrl       *string
	CreatedAt       time.Time
	PushedAt        time.Time
	Homepage        *string
	RepoSize        int
	WatchersCount   int
	IsDisabled      bool
	OpenIssuesCount int
	License         *License
	AllowsForking   bool
	IsTemplate      bool
	Topics          []string
	DefaultBranch   string
}
