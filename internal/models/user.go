package models

import "time"

type BaseUser struct {
	Username string
	Homepage string
}
type User struct {
	BaseUser
	Email            string
	IsPublic         bool
	Name             string
	Location         string
	Hirable          bool
	Bio              *string
	PublicReposCount int
	PublicGistsCount int
	FollowersCount   int
	FollowingCount   int
	CreatedAt        time.Time
	UpdatedAt        time.Time
}
