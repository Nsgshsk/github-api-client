package models

type UserReport struct {
	User          *User
	TopRepos      []BaseRepository
	Languages     LanguageDistribution
	ReposAnalyzed int
	TotalStars    int
	Warnings      []string
}
