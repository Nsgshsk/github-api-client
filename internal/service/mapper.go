package service

import "github.com/Nsgshsk/github-api-client.git/internal/models"

func mapLangPercentage(langs *models.Languages) models.LangPercentage {
	result := make(models.LangPercentage, len(*langs))
	var sum float32 = 0
	for _, lang := range *langs {
		sum += float32(lang)
	}
	for idx, lang := range *langs {
		result[idx] += (float32(lang) / sum) * 100
	}
	return result
}
