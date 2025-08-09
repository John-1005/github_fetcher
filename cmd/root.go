package cmd

import (
	"fmt"
	"os"
	"strings"

	"github.com/John-1005/github_fetcher/internal/charm"
	"github.com/John-1005/github_fetcher/internal/githubapi"
	"github.com/spf13/cobra"
)

var sortOrder string
var limit int
var verbose bool
var noCache bool
var clearCache bool

var rootCmd = &cobra.Command{
	Use:   "githubfetcher [username]",
	Short: "Fetches github respositories for a given user",
	Long: `githubfetcher is a command-line tool to retreive and display
				public GitHub respositories for any specified user. It organizes
				the reposotories in an astethically pleasing table, sorted by
				the number of stargazers

				Usage Examples:
				githubfetcher octocat
				githubfetcher johnsmith --sort asc`,

	RunE: func(cmd *cobra.Command, args []string) error {

		if clearCache {
			client := githubapi.NewClient()
			client.ClearCache()
			os.Remove(client.CacheFilePath())
			fmt.Println("Cache cleared")
			return nil
		}

		if len(args) == 0 {
			return fmt.Errorf("Github username is required")
		}

		userName := args[0]
		order := strings.ToLower(sortOrder)

		client := githubapi.NewClient()

		repos, err := client.GetRepositories(userName, verbose, noCache)
		if err != nil {
			return err
		}

		sortedRepos := githubapi.SortRepositories(repos, order)

		if limit > 0 && len(sortedRepos) > limit {
			sortedRepos = sortedRepos[:limit]
		}

		fmt.Println(charm.RenderTable(sortedRepos))

		return nil
	},
}

func init() {
	rootCmd.Flags().StringVarP(&sortOrder, "sort", "s", "desc", "Sort repositories by stargazers: 'asc or 'desc'")
	rootCmd.Flags().IntVarP(&limit, "limit", "l", 0, "Limit the number for repositories requested")
	rootCmd.Flags().BoolVarP(&verbose, "verbose", "v", false, "Enables Verbose output")
	rootCmd.Flags().BoolVarP(&noCache, "no-cache", "", false, "Disables cache to always fetch fresh data")
	rootCmd.Flags().BoolVarP(&clearCache, "clear-cache", "", false, "Cleares saved cache")
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		os.Exit(1)
	}
}
