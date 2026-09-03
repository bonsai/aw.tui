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
	case loadMsg:
		if v.err != nil { m.status = "load failed: " + v.err.Error(); return m, nil }
		m.repos, m.graph, m.status = v.repos, v.graph, fmt.Sprintf("%d repos • %d graph nodes • %d edges", len(v.repos), len(v.graph.Nodes), len(v.graph.Edges))
	case tea.KeyMsg:
		switch v.String() {
		case "q", "ctrl+c": return m, tea.Quit
		case "j", "down": if len(m.repos) > 0 { m.cursor = min(m.cursor+1, len(m.repos)-1) }
		case "k", "up": if m.cursor > 0 { m.cursor-- }
		case "1": m.mode = 0
		case "2": m.mode = 1
		case "3": m.mode = 2
		case "r": return m, m.Init()
		}
	}
	return m, nil
}

func (m model) View() string {
	if m.width == 0 { return "aw.tui" }
	title := accent.Render("🌱 AW.TUI") + "  " + muted.Render("repo-first agentic workflow console")
	nav := "[1] REPOS   [2] AGENTS   [3] GRAPH   [r] refresh   [q] quit"
	left := ""
	for i, r := range m.repos {
		prefix := "  "
		if i == m.cursor { prefix = "> " }
		left += prefix + r.Name + "\n"
	}
	if left == "" { left = "  " + m.status }
	leftPanel := box.Width(max(24, m.width/3)).Render("REPOSITORIES\n\n" + left)

	detail := ""
	if len(m.repos) > 0 {
		r := m.repos[m.cursor]
		detail = accent.Render(r.Name) + "\n\n" + r.Description + "\n\n" +
			"branch: " + r.DefaultBranch + "\n" + "language: " + r.Language
	}
	if m.mode == 1 { detail = "AGENTS\n\nAgent registry will be loaded from bonsai/.company.\n\nRoles → capabilities → repositories." }
	if m.mode == 2 { detail = fmt.Sprintf("SELF-ORGANIZING GRAPH\n\nNodes: %d\nEdges: %d\n\nThe graph is derived from the recent repository set.\nML/embedding ranking can be attached here next.", len(m.graph.Nodes), len(m.graph.Edges)) }
	rightPanel := box.Width(max(30, m.width-m.width/3-6)).Render(detail)
	body := lipgloss.JoinHorizontal(lipgloss.Top, leftPanel, rightPanel)
	return title + "\n" + muted.Render(nav) + "\n\n" + body + "\n" + muted.Render("\n"+strings.TrimSpace(m.status))
}

func min(a,b int) int { if a < b { return a }; return b }
func max(a,b int) int { if a > b { return a }; return b }
