# github-api-client

A command-line client for the public GitHub API. It reads a list of usernames from a file, fetches each user's public profile, repositories, and per-repository language data, and prints a serial summary of all users as tables on stdout.

For every user the report contains:

- Basic profile information
- Top repositories by star count
- Language distribution across all public repositories, computed from GitHub's language byte counts and grouped into percentages

## Build

```bash
go build -o client ./cmd/client
```

Or run directly:

```bash
go run ./cmd/client usernames.txt
```

## Username file

One username per line. Blank lines and lines starting with `#` are ignored. Duplicates are removed case-insensitively, and file order is preserved.

```text
# users to analyze
octocat
torvalds
golang
```

Use `-` as the file argument to read usernames from stdin.

## Usage

```text
Usage: client [flags] <usernames-file>

Flags:
  --api-url string          override GitHub API base URL
  --auth-mode string        auth mode: auto, unauthenticated, authenticated (default "auto")
  --cache-ttl duration      in-memory cache TTL (default 5m0s)
  --exclude-archived        exclude archived repositories
  --exclude-forks           exclude forked repositories
  --no-cache                disable request caching
  --others-threshold float  group languages below this percentage as Others (default 5)
  --timeout duration        overall run timeout, 0 disables (default 2m0s)
  --top int                 number of top repositories to show (default 5)
```

Examples:

```bash
# unauthenticated, defaults
./client usernames.txt

# exclude forks and archived repositories, show top 10 repos
./client --exclude-forks --exclude-archived --top 10 usernames.txt

# group languages below 2% as "Others"
./client --others-threshold 2 usernames.txt
```

## Authentication and rate limits

By default (`--auth-mode=auto`) the client uses `GITHUB_TOKEN` when it is set and falls back to unauthenticated requests otherwise.

| Mode | Behavior |
|---|---|
| `auto` | Use `GITHUB_TOKEN` if present, otherwise unauthenticated |
| `unauthenticated` | Never send a token, even if `GITHUB_TOKEN` is set |
| `authenticated` | Require `GITHUB_TOKEN`; exit immediately with an error if it is missing |

```bash
export GITHUB_TOKEN=ghp_...
./client usernames.txt
```

Unauthenticated requests are limited to 60 requests per hour. Since the client fetches language statistics for every repository, a token is strongly recommended for more than a couple of users. Language requests and repository pages are cached in memory for the duration of the run (`--cache-ttl`, default 5 minutes; disable with `--no-cache`).

## Output

```text
User: octocat (The Octocat)
┌────────────────┬────────────────────────────┐
│ Field          │ Value                      │
├────────────────┼────────────────────────────┤
│ Username       │ octocat                    │
│ Name           │ The Octocat                │
│ Location       │ San Francisco              │
│ Followers      │ 16000                      │
│ Following      │ 9                          │
│ Public repos   │ 8                          │
│ Public gists   │ 8                          │
│ Repos analyzed │ 8                          │
│ Total stars    │ 2500                       │
│ Profile        │ https://github.com/octocat │
│ Member since   │ 2011-01-25                 │
└────────────────┴────────────────────────────┘

Top repositories
┌──────┬─────────────────────┬───────┬───────┬──────────┐
│ Rank │ Repository          │ Stars │ Forks │ Language │
├──────┼─────────────────────┼───────┼───────┼──────────┤
│    1 │ octocat/Hello-World │  2500 │  1200 │ Go       │
└──────┴─────────────────────┴───────┴───────┴──────────┘

Language distribution
┌──────────┬────────────┐
│ Language │ Percentage │
├──────────┼────────────┤
│ Go       │      62.4% │
│ Shell    │      22.1% │
│ Others   │      15.5% │
└──────────┴────────────┘
```

- Languages below `--others-threshold` are summed into `Others`, always shown last.
- The largest language is always shown even if it falls below the threshold.
- Percentages are derived from the byte counts returned by the GitHub languages endpoint.
- Errors for individual users are printed to stderr; the client continues with the remaining users and exits non-zero if any failed.

## Exit codes

| Code | Meaning |
|---|---|
| 0 | All users processed successfully |
| 1 | File, runtime, or per-user failure |
| 2 | Usage or authentication configuration error |
| 130 | Interrupted |

## Development

```bash
gofmt -l .
go build ./...
go vet ./...
go test -race ./...
```
