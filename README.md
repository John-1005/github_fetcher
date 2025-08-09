# GitHub Fetcher CLI

This is my first real solo project after progressing through Boot.dev. I really wanted to learn how to build a CLI tool and advance the knowledge of API calls, adding a cache and then how to add a table with style.

This is a command-line tool written in Go to fetch and display public repositories for any GitHub user. It has sorting and persistent caching to reduce API calls. 


## Features
- Fetch public repositories for any GitHub user using the GitHub API.
- Sort repositories by stargazer count (asc or desc).
- Limit the number of displayed repositories.
- Persistent caching between runs to avoid unnecessary API calls.
- Verbose mode to see detailed execution steps.
- --no-cache flag to skip cache for a single run.
- --clear-cache flag to wipe the cache file entirely.
- Pretty table output using Charmbracelet’s Lip Gloss styling.

## Installation
From source:
git clone https://github.com/<john-1005>/github_fetcher.git
cd github_fetcher
go build -o githubfetcher


## Usage
githubfetcher [username] [flags]

### Flags
--sort (-s)         Sort order: asc or desc (default desc)  
--limit (-l)        Limit the number of repositories displayed (default 0 = no limit)  
--verbose (-v)      Enable verbose output  
--no-cache          Disable cache for this run  
--clear-cache       Clear the persistent cache and exit  
--help (-h)         Show help message

## **Examples**

### Fetch repos for a user (default: sort by stars descending)

githubfetcher octocat

- Sort ascending by stars:

githubfetcher octocat --sort asc

- Limit to top 5 repos:

githubfetcher octocat --limit 5

- Verbose mode (see cache hits/misses, API calls, etc.):

githubfetcher octocat --verbose

- Skip cache for this run:

githubfetcher octocat --no-cache

- Clear the cache and exit:

githubfetcher --clear-cache

## Cache Behavior
- Cache file is stored at:
~/.cache/github_fetcher/cache.json

- Cache persists between runs.
- --no-cache skips reading/writing cache for a single run.
- --clear-cache wipes the cache file entirely.
- Verbose mode will show:
- [verbose] Cache found for <username> → cache hit
- [verbose] No cache found for <username> → cache miss

## **Example Verbose Output**
```text
[verbose] Cache found for octocat
┌────────────────────┬────────┐
│  Repository Name   │ Stars  │
├────────────────────┼────────┤
│ Spoon-Knife        │ 13,151 │
│ Hello-World        │  3,059 │
│ octocat.github.io  │    886 │
│ hello-worId        │    580 │
│ linguist           │    575 │
│ git-consortium     │    450 │
│ boysenberry-repo-1 │    350 │
│ test-repo1         │    339 │
└────────────────────┴────────┘
```


## Rate Limits
- Unauthenticated requests: 60 requests/hour per IP.
- Authenticated requests (future feature): 5,000 requests/hour.
- Verbose mode shows:
  - Total rate limit
  - Remaining requests
  - Reset time

---

## Dependencies
This project uses:
- Lip Gloss — for terminal styling.
- Lip Gloss Table — for table rendering.

End users:  
You do not need to install these manually. They are compiled into the binary when you build the project.

Developers building from source:  
Dependencies are managed automatically via Go modules. When you run `go build` or `go run`, Go will download them based on the go.mod file.






