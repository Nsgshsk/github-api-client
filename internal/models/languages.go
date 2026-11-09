package models

type Languages map[string]int64

type LanguageShare struct {
	Name       string
	Percentage float64
}

type LanguageDistribution []LanguageShare
