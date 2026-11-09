package input

import (
	"bufio"
	"errors"
	"fmt"
	"io"
	"os"
	"strings"
)

// ReadUsernames reads usernames from path. Use "-" to read from stdin.
func ReadUsernames(path string) ([]string, error) {
	if path == "-" {
		return ParseUsernames(os.Stdin)
	}

	file, err := os.Open(path)
	if err != nil {
		return nil, fmt.Errorf("open %s: %w", path, err)
	}
	defer file.Close()

	return ParseUsernames(file)
}

// ParseUsernames reads one username per line, ignoring blank lines and comments.
func ParseUsernames(r io.Reader) ([]string, error) {
	scanner := bufio.NewScanner(r)
	scanner.Buffer(make([]byte, 0, 64*1024), 1<<20)

	seen := make(map[string]struct{})
	var usernames []string

	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if line == "" || strings.HasPrefix(line, "#") {
			continue
		}

		key := strings.ToLower(line)
		if _, exists := seen[key]; exists {
			continue
		}
		seen[key] = struct{}{}
		usernames = append(usernames, line)
	}
	if err := scanner.Err(); err != nil {
		return nil, fmt.Errorf("read usernames: %w", err)
	}
	if len(usernames) == 0 {
		return nil, errors.New("no usernames found")
	}
	return usernames, nil
}
