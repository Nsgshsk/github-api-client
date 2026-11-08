package models

import "time"

type Repository struct {
	Name            string
	IsPrivate       bool
	Owner           BaseUser
	Url             string
	MirrorUrl       *string
	IsFork          bool
	CreatedAt       time.Time
	UpdatedAt       time.Time
	PushedAt        time.Time
	Homepage        *string
	RepoSize        int
	StarsCount      int
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
