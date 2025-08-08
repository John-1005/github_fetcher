package cmd

import (
	"fmt"
	"os"

	"github.com/John-1005/github_fetcher/internal/githubapi"
	"github.com/spf13/cobra"
)

var sortOrder string

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
		if len(args) == 0 {
			return fmt.Errorf("Github username is required")
		}

		userName := args[0]

		client := githubapi.NewClient()

		repos, err := client.GetRepositories(userName)
		if err != nil {
			return err
		}

		fmt.Printf("Fetched %d repositories for %s:\n", len(repos), userName)
		for _, r := range repos {
			fmt.Printf("- %s (%d stars)\n", r.Name, r.StargazersCount)
		}

		return nil
	},
}

func init() {
	rootCmd.Flags().StringVarP(&sortOrder, "sort", "s", "desc", "Sort repositories by stargazers: 'asc or 'desc'")
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		os.Exit(1)
	}
}
