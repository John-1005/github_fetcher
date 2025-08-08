package githubapi

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"strconv"
	"time"
)

type Client struct {
	BaseURL string
	//Eventual caching added
}

func NewClient() Client {
	return Client{
		BaseURL: "https://api.github.com/users/",
	}
}

func (c *Client) GetRepositories(username string) ([]Repository, error) {
	if username == "" {
		return nil, errors.New("username cannot be empty")
	}

	url := c.BaseURL + username + "/repos"

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to make request: %w", err)
	}

	req.Header.Set("User-Agent", "githubfetcher")

	client := &http.Client{}

	rsp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	if rsp.StatusCode == http.StatusForbidden {
		remaining := rsp.Header.Get("X-RateLimit-Remaining")

		reset := rsp.Header.Get("X-RateLimit-Reset")

		if remaining == "0" {
			resetUnix, err := strconv.ParseInt(reset, 10, 64)
			if err != nil {
				return nil, fmt.Errorf("Github API rate limit exceeded, Resets at UNIX time: %s", reset)
			}

			resetTime := time.Unix(resetUnix, 0)
			return nil, fmt.Errorf("GitHub API rate limit exceeded, Resets at %s", resetTime)
		}
		return nil, fmt.Errorf("access forbidden, unauthorized or over rate limit")
	}

	if rsp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("unexpected status code: %d", rsp.StatusCode)
	}

	defer rsp.Body.Close()

	body, err := io.ReadAll(rsp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response body: %w", err)
	}

	var repos []Repository
	err = json.Unmarshal(body, &repos)
	if err != nil {
		return nil, fmt.Errorf("failed to unmarshal JSON: %w", err)
	}

	return repos, nil
}
