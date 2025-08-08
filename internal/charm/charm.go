package charm

import (
	"strconv"

	"github.com/John-1005/github_fetcher/internal/githubapi"
	"github.com/charmbracelet/bubbles/table"
	"github.com/charmbracelet/lipgloss"
)

func RenderTable(repos []githubapi.Repository) string {

	greenText := lipgloss.NewStyle().Foreground(lipgloss.Color("#04B575"))  //Green
	yellowText := lipgloss.NewStyle().Foreground(lipgloss.Color("#FFDE21")) //Star Yellow
	style := lipgloss.NewStyle().Padding(2)

	columns := []table.Column{
		{Title: "Repository Name", Width: 30},
		{Title: "Stars", Width: 4},
	}

	rows := []table.Row{}

	for _, repo := range repos {
		rows = append(rows, table.Row{
			greenText.Render(repo.Name),
			yellowText.Render(strconv.Itoa(repo.StargazersCount)),
		})
	}

	t := table.New(
		table.WithColumns(columns),
		table.WithRows(rows),
		table.WithFocused(true),
		table.WithHeight(7),
	)

	return style.Render(t.View())

}
