package main

import (
	"errors"
	"flag"
	"fmt"
	"io"
	"log"
	"os"
	"strings"

	"github.com/Nsgshsk/github-api-client/internal/display"
	"github.com/Nsgshsk/github-api-client/internal/github"
	"github.com/Nsgshsk/github-api-client/internal/input"
	"github.com/Nsgshsk/github-api-client/internal/service"
)

type authMode string

const (
	authAuto            authMode = "auto"
	authUnauthenticated authMode = "unauthenticated"
	authAuthenticated   authMode = "authenticated"
)

func main() {
	os.Exit(run(os.Args[1:], os.Stdout, os.Stderr))
}

func run(args []string, stdout, stderr io.Writer) int {
	logger := log.New(stderr, "", 0)

	flags := flag.NewFlagSet("client", flag.ContinueOnError)
	flags.SetOutput(stderr)
	flags.Usage = func() {
		logger.Print("Usage: client [flags] <usernames-file>")
		logger.Print("")
		logger.Print("Flags:")
		flags.PrintDefaults()
	}

	top := flags.Int("top", 5, "number of top repositories to show")
	othersThreshold := flags.Float64("others-threshold", 5, "group languages below this percentage as Others")
	excludeForks := flags.Bool("exclude-forks", false, "exclude forked repositories")
	excludeArchived := flags.Bool("exclude-archived", false, "exclude archived repositories")
	authModeFlag := flags.String("auth-mode", string(authAuto), "auth mode: auto, unauthenticated, authenticated")
	apiURL := flags.String("api-url", "", "override GitHub API base URL")

	if err := flags.Parse(args); err != nil {
		if errors.Is(err, flag.ErrHelp) {
			return 0
		}
		return 2
	}
	if flags.NArg() != 1 {
		flags.Usage()
		return 2
	}
	if *top < 0 {
		logger.Print("--top must be greater than or equal to 0")
		return 2
	}
	if *othersThreshold < 0 {
		logger.Print("--others-threshold must be greater than or equal to 0")
		return 2
	}

	mode, err := parseAuthMode(*authModeFlag)
	if err != nil {
		logger.Print(err)
		return 2
	}
	token, err := resolveToken(mode)
	if err != nil {
		logger.Print(err)
		return 2
	}

	usernames, err := input.ReadUsernames(flags.Arg(0))
	if err != nil {
		logger.Printf("read usernames: %v", err)
		return 1
	}

	options := []github.Option{}
	if *apiURL != "" {
		options = append(options, github.WithBaseURL(*apiURL))
	}
	if token != "" {
		options = append(options, github.WithToken(token))
	}
	client := github.NewClient(options...)

	reportService := service.NewReportService(client, service.ReportConfig{
		TopRepos:        *top,
		OthersThreshold: *othersThreshold,
		ExcludeForks:    *excludeForks,
		ExcludeArchived: *excludeArchived,
	})

	failed := false
	for _, username := range usernames {
		report, err := reportService.GetUserReport(username)
		if err != nil {
			printUserError(logger, username, err)
			failed = true
			continue
		}

		for _, warning := range report.Warnings {
			logger.Printf("%s: %s", username, warning)
			failed = true
		}

		if err := display.RenderUserReport(stdout, report); err != nil {
			logger.Printf("render report for %s: %v", username, err)
			return 1
		}
	}

	if failed {
		return 1
	}
	return 0
}

func parseAuthMode(value string) (authMode, error) {
	switch authMode(strings.ToLower(strings.TrimSpace(value))) {
	case "", authAuto:
		return authAuto, nil
	case authUnauthenticated:
		return authUnauthenticated, nil
	case authAuthenticated:
		return authAuthenticated, nil
	default:
		return "", fmt.Errorf("invalid --auth-mode %q: expected auto, unauthenticated, or authenticated", value)
	}
}

func resolveToken(mode authMode) (string, error) {
	token := strings.TrimSpace(os.Getenv("GITHUB_TOKEN"))
	switch mode {
	case authUnauthenticated:
		return "", nil
	case authAuthenticated:
		if token == "" {
			return "", errors.New("GITHUB_TOKEN must be set when --auth-mode=authenticated")
		}
		return token, nil
	default:
		return token, nil
	}
}

func printUserError(logger *log.Logger, username string, err error) {
	if github.IsRateLimited(err) {
		logger.Printf("%s: %v (hint: set GITHUB_TOKEN or use --auth-mode=authenticated)", username, err)
		return
	}
	logger.Printf("%s: %v", username, err)
}
