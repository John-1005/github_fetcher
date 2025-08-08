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

func (c *Client) GetRepositories(username string, verbose bool) ([]Repository, error) {
	if username == "" {
		return nil, errors.New("username cannot be empty")
	}

	page := 1
	perPage := 60
	searchedRepos := []Repository{}

	for {
		url := c.BaseURL + username + "/repos?per_page=" + strconv.Itoa(perPage) + "&page=" + strconv.Itoa(page)

		if verbose {
			fmt.Printf("[verbose] Fetching page %d: %s\n", page, url)
		}

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

		if verbose {
			fmt.Printf("[verbose] Status code :%d\n", rsp.StatusCode)
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

		body, err := io.ReadAll(rsp.Body)
		rsp.Body.Close()
		if err != nil {
			return nil, fmt.Errorf("failed to read response body: %w", err)
		}

		var repos []Repository
		err = json.Unmarshal(body, &repos)
		if err != nil {
			return nil, fmt.Errorf("failed to unmarshal JSON: %w", err)
		}

		if verbose {
			fmt.Printf("[verbose] Page %d returned %d repos\n", page, len(repos))
		}

		searchedRepos = append(searchedRepos, repos...)

		if len(repos) < perPage {
			break
		}

		page++
	}

	if verbose {
		fmt.Printf("[verbose] Total repos fetched: %d\n", len(searchedRepos))
	}

	return searchedRepos, nil
}
