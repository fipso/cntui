package main

import (
	"fmt"
	"os"
	"strconv"

	"github.com/charmbracelet/bubbles/table"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type model struct {
	table table.Model
}

var baseStyle = lipgloss.NewStyle().
	BorderStyle(lipgloss.NormalBorder()).
	BorderForeground(lipgloss.Color("240"))

func pullRequests() tea.Msg {
	// Block till update
	<-requestsUpdated

	return true
}

func (m model) Init() tea.Cmd { return pullRequests }

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "esc":
			if m.table.Focused() {
				m.table.Blur()
			} else {
				m.table.Focus()
			}
		case "q", "ctrl+c":
			return m, tea.Quit
		case "enter":
			id := m.table.SelectedRow()[1]
			for _, req := range requests {
				if string(req.id) == id {
					exportAsCurl(&req)
					break
				}
			}

			return m, tea.Quit
		}
	case bool:

		var rows []table.Row

		for _, req := range requests {
			status := ""
			if req.res != nil {
				status = strconv.Itoa(req.res.Status)
			}

			row := table.Row{
				req.time.Format("15:04:05.99"), string(req.id),
				req.req.Method, req.req.URL, req.initiator.Type, status,
			}

			rows = append(rows, row)
		}

		m.table.SetRows(rows)
		// m.table.GotoBottom()
		return m, pullRequests
	}
	m.table, cmd = m.table.Update(msg)
	return m, cmd
}

func (m model) View() string {
	return baseStyle.Render(m.table.View()) + "\n"
}

func setupTUI() {
	columns := []table.Column{
		{Title: "Time", Width: 15},
		{Title: "ID", Width: 10},
		{Title: "Method", Width: 10},
		{Title: "URL", Width: 50},
		{Title: "Initiator", Width: 10},
		{Title: "Status", Width: 10},
	}

	t := table.New(
		table.WithColumns(columns),
		table.WithRows([]table.Row{}),
		table.WithFocused(true),
		table.WithHeight(25),
	)

	s := table.DefaultStyles()
	s.Header = s.Header.
		BorderStyle(lipgloss.NormalBorder()).
		BorderForeground(lipgloss.Color("240")).
		BorderBottom(true).
		Bold(false)
	s.Selected = s.Selected.
		Foreground(lipgloss.Color("229")).
		Background(lipgloss.Color("57")).
		Bold(false)
	t.SetStyles(s)

	m := model{t}
	if _, err := tea.NewProgram(m).Run(); err != nil {
		fmt.Println("Error running program:", err)
		os.Exit(1)
	}
}
