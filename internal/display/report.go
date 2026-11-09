package display

import (
	"errors"
	"fmt"
	"io"
	"strconv"
	"strings"

	"github.com/Nsgshsk/github-api-client/internal/models"
)

const maxCellRunes = 80

func RenderUserReport(w io.Writer, report *models.UserReport) error {
	if report == nil || report.User == nil {
		return errors.New("render report: nil report")
	}

	user := report.User
	if user.Name != "" {
		fmt.Fprintf(w, "User: %s (%s)\n", user.Username, user.Name)
	} else {
		fmt.Fprintf(w, "User: %s\n", user.Username)
	}

	if err := renderBasicInfo(w, report); err != nil {
		return err
	}

	fmt.Fprintln(w, "\nTop repositories")
	if err := renderTopRepos(w, report.TopRepos); err != nil {
		return err
	}

	fmt.Fprintln(w, "\nLanguage distribution")
	if err := renderLanguages(w, report.Languages); err != nil {
		return err
	}

	fmt.Fprintln(w)
	return nil
}

func renderBasicInfo(w io.Writer, report *models.UserReport) error {
	user := report.User
	table := NewTable(
		Column{Header: "Field"},
		Column{Header: "Value"},
	)

	table.AddRow("Username", user.Username)
	addIfNotEmpty(table, "Name", user.Name)
	addIfNotEmpty(table, "Location", user.Location)
	if user.Company != nil {
		addIfNotEmpty(table, "Company", *user.Company)
	}
	addIfNotEmpty(table, "Blog", user.Blog)
	if user.TwitterUsername != "" {
		addIfNotEmpty(table, "Twitter", "@"+strings.TrimPrefix(user.TwitterUsername, "@"))
	}
	addIfNotEmpty(table, "Email", user.Email)
	if user.Hireable {
		table.AddRow("Hireable", "Yes")
	}
	if user.Bio != nil {
		addIfNotEmpty(table, "Bio", truncate(*user.Bio, maxCellRunes))
	}
	table.AddRow("Followers", strconv.Itoa(user.FollowersCount))
	table.AddRow("Following", strconv.Itoa(user.FollowingCount))
	table.AddRow("Public repos", strconv.Itoa(user.PublicReposCount))
	table.AddRow("Public gists", strconv.Itoa(user.PublicGistsCount))
	table.AddRow("Repos analyzed", strconv.Itoa(report.ReposAnalyzed))
	table.AddRow("Total stars", strconv.Itoa(report.TotalStars))
	addIfNotEmpty(table, "Profile", user.Homepage)
	if !user.CreatedAt.IsZero() {
		table.AddRow("Member since", user.CreatedAt.Format("2006-01-02"))
	}

	return table.Render(w)
}

func renderTopRepos(w io.Writer, repos []models.BaseRepository) error {
	table := NewTable(
		Column{Header: "Rank", Align: AlignRight},
		Column{Header: "Repository"},
		Column{Header: "Stars", Align: AlignRight},
		Column{Header: "Forks", Align: AlignRight},
		Column{Header: "Language"},
	)

	if len(repos) == 0 {
		table.AddRow("-", "No public repositories", "-", "-", "-")
		return table.Render(w)
	}

	for i, repo := range repos {
		table.AddRow(
			strconv.Itoa(i+1),
			repo.Name,
			strconv.Itoa(repo.StarsCount),
			strconv.Itoa(repo.ForksCount),
			valueOrDash(repo.PrimaryLanguage),
		)
	}
	return table.Render(w)
}

func renderLanguages(w io.Writer, languages models.LanguageDistribution) error {
	table := NewTable(
		Column{Header: "Language"},
		Column{Header: "Percentage", Align: AlignRight},
	)

	if len(languages) == 0 {
		table.AddRow("No language data", "")
		return table.Render(w)
	}

	for _, language := range languages {
		table.AddRow(language.Name, fmt.Sprintf("%.1f%%", language.Percentage))
	}
	return table.Render(w)
}

func addIfNotEmpty(table *Table, field, value string) {
	if strings.TrimSpace(value) != "" {
		table.AddRow(field, value)
	}
}

func valueOrDash(value string) string {
	if strings.TrimSpace(value) == "" {
		return "-"
	}
	return value
}

func truncate(value string, max int) string {
	runes := []rune(value)
	if len(runes) <= max {
		return value
	}
	if max <= 1 {
		return string(runes[:max])
	}
	return string(runes[:max-1]) + "…"
}
