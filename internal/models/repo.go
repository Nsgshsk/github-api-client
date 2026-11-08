package models

import "time"

type BaseRepository struct {
	Name       string
	IsPrivate  bool
	URL        string
	StarsCount int
}

type Repository struct {
	BaseRepository
	Owner           BaseUser
	MirrorUrl       *string
	IsFork          bool
	CreatedAt       time.Time
	UpdatedAt       time.Time
	PushedAt        time.Time
	Homepage        *string
	RepoSize        int
	WatchersCount   int
	ForksCount      int
	IsArchived      bool
	IsDisabled      bool
	OpenIssuesCount int
	License         *License
	AllowsForking   bool
	IsTemplate      bool
	Topics          []string
	DefaultBranch   string
}
