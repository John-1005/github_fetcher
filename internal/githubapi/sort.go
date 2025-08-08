package githubapi

import (
	"sort"
)

func SortRepositories(repos []Repository, order string) []Repository {
	sort.Slice(repos, func(i, j int) bool {
		if order == "asc" {
			return repos[i].StargazersCount < repos[j].StargazersCount
		}

		return repos[i].StargazersCount > repos[j].StargazersCount
	})
	return repos
}
