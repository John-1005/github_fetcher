package charm

import (
	"github.com/John-1005/github_fetcher/internal/githubapi"
	"github.com/charmbracelet/lipgloss"
	"github.com/charmbracelet/lipgloss/table"
	"golang.org/x/text/message"
)

func RenderTable(repos []githubapi.Repository) string {

	green := lipgloss.Color("#04B575")
	yellow := lipgloss.Color("#FFDE21")
	lightGray := lipgloss.Color("808080")

	printer := message.NewPrinter(message.MatchLanguage("en"))

	maxNameLen := 0

	for _, repo := range repos {
		if len(repo.Name) > maxNameLen {
			maxNameLen = len(repo.Name)
		}
	}

	colWidthName := maxNameLen + 2

	headerStyle := lipgloss.NewStyle().
		Foreground(lightGray).
		Bold(true).
		Align(lipgloss.Center)

	repoNameStyle := lipgloss.NewStyle().
		Foreground(green).
		Padding(0, 1).
		Width(colWidthName)

	starStyle := lipgloss.NewStyle().
		Foreground(yellow).
		Padding(0, 1).
		Width(8).
		Align(lipgloss.Right)

	var rows [][]string

	for _, repo := range repos {
		rows = append(rows, []string{
			repo.Name,
			printer.Sprintf("%d", repo.StargazersCount),
		})
	}

	t := table.New().
		Border(lipgloss.NormalBorder()).
		BorderStyle(lipgloss.NewStyle().Foreground(lightGray)).
		Headers("Repository Name", "Stars").
		Rows(rows...).
		StyleFunc(func(row, col int) lipgloss.Style {
			if row == table.HeaderRow {
				return headerStyle
			}
			if col == 0 {
				return repoNameStyle
			}
			return starStyle
		})

	return t.Render()

}
