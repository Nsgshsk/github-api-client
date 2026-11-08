package models

type UserOverview struct {
	User      *User
	Top5Repos []BaseRepository
}
