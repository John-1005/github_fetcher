package githubapi

type Repository struct {
	Name            string `json:"name"`
	StargazersCount int    `json:"stargazers_count"`
}
