package main

import (
	"fmt"
	"os"
	"strings"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/bonsai/aw.tui/graph"
	"github.com/bonsai/aw.tui/repos"
)

type model struct {
	repos  []repos.Repo
	graph  graph.Graph
	cursor int
	mode   int
	width  int
	height int
	status string
}

var (
	accent = lipgloss.NewStyle().Bold(true)
	muted  = lipgloss.NewStyle().Faint(true)
	box    = lipgloss.NewStyle().Border(lipgloss.RoundedBorder()).Padding(1)
)

func main() {
	p := tea.NewProgram(model{status: "loading recent repositories…"}, tea.WithAltScreen())
	if _, err := p.Run(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

func (m model) Init() tea.Cmd {
	return func() tea.Msg {
		rs, err := repos.FetchRecent("bonsai", 30)
		if err != nil { return loadMsg{err: err} }
		return loadMsg{repos: rs, graph: graph.Build(rs)}
	}
}

type loadMsg struct { repos []repos.Repo; graph graph.Graph; err error }

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch v := msg.(type) {
	case tea.WindowSizeMsg:
		m.width, m.height = v.Width, v.Height
		m.cursor = clampCursor(m.cursor, len(m.repos), m.visibleRows())
	case loadMsg:
		if v.err != nil { m.status = "load failed: " + v.err.Error(); return m, nil }
		m.repos, m.graph = v.repos, v.graph
		m.cursor = clampCursor(m.cursor, len(m.repos), m.visibleRows())
		m.status = fmt.Sprintf("%d repos • %d graph nodes • %d edges", len(m.repos), len(m.graph.Nodes), len(m.graph.Edges))
	case tea.KeyMsg:
		switch v.String() {
		case "q", "ctrl+c": return m, tea.Quit
		case "j", "down": if m.cursor < len(m.repos)-1 { m.cursor++ }
		case "k", "up": if m.cursor > 0 { m.cursor-- }
		case "1": m.mode = 0
		case "2": m.mode = 1
		case "3": m.mode = 2
		case "r": return m, m.Init()
		}
	}
	return m, nil
}

// visibleRows is the number of repository rows that fit in the fixed panel.
// We intentionally do not scroll: the selected row moves through the current
// viewport and the list is re-rendered around it.
func (m model) visibleRows() int {
	rows := m.height - 10
	if rows < 1 { return 1 }
	return rows
}

func clampCursor(cursor, total, visible int) int {
	if total <= 0 { return 0 }
	if cursor < 0 { return 0 }
	if cursor >= total { return total - 1 }
	return cursor
}

func (m model) View() string {
	if m.width == 0 || m.height == 0 { return "aw.tui" }

	title := accent.Render("🌱 AW.TUI") + "  " + muted.Render("repo-first agentic workflow console")
	nav := "[1] REPOS   [2] AGENTS   [3] GRAPH   [r] refresh   [q] quit"

	// Fixed viewport: never let the repository list grow beyond the terminal.
	// Instead, re-render a window centered on the selected repository.
	rows := m.visibleRows()
	start := 0
	if len(m.repos) > rows {
		start = m.cursor - rows/2
		if start < 0 { start = 0 }
		if start > len(m.repos)-rows { start = len(m.repos)-rows }
	}
	end := min(start+rows, len(m.repos))

	var list []string
	for i := start; i < end; i++ {
		prefix := "  "
		if i == m.cursor { prefix = "> " }
		list = append(list, prefix+m.repos[i].Name)
	}
	left := strings.Join(list, "\n")
	if left == "" { left = "  " + m.status }

	leftWidth := max(24, m.width/3-2)
	rightWidth := max(30, m.width-leftWidth-6)
	leftPanel := box.Width(leftWidth).Height(rows+2).Render("REPOSITORIES\n\n" + left)

	detail := ""
	if len(m.repos) > 0 {
		r := m.repos[m.cursor]
		detail = accent.Render(r.Name) + "\n\n" + r.Description + "\n\n" +
			"branch: " + r.DefaultBranch + "\n" + "language: " + r.Language
	}
	if m.mode == 1 {
		detail = "AGENTS\n\nAgent registry will be loaded from bonsai/.company.\n\nRoles → capabilities → repositories."
	}
	if m.mode == 2 {
		detail = fmt.Sprintf("SELF-ORGANIZING GRAPH\n\nNodes: %d\nEdges: %d\n\nThe graph is derived from the recent repository set.\nML/embedding ranking can be attached here next.", len(m.graph.Nodes), len(m.graph.Edges))
	}
	rightPanel := box.Width(rightWidth).Height(rows+2).Render(detail)
	body := lipgloss.JoinHorizontal(lipgloss.Top, leftPanel, rightPanel)

	footer := muted.Render(strings.TrimSpace(m.status))
	return title + "\n" + muted.Render(nav) + "\n\n" + body + "\n" + footer
}

func min(a, b int) int { if a < b { return a }; return b }
func max(a, b int) int { if a > b { return a }; return b }
