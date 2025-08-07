package cmd

import (
	"fmt"
	"github.com/spf13/cobra"
	"os"
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

		fmt.Printf("Fetching reposotories for: %s\n", userName)

		return nil
	},
}

func init() {
	rootCmd.Flags().StringVarP(&sortOrder, "sort", "s", "", "Sort repositories by stargazers")
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		os.Exit(1)
	}
}
