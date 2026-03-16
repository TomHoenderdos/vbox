package tui

import (
	"fmt"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/TomHoenderdos/vbox/internal/config"
	"github.com/TomHoenderdos/vbox/internal/vagrant"
	"github.com/charmbracelet/bubbles/table"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type project struct {
	dir    string
	config *config.Config
	status string
}

type psModel struct {
	table    table.Model
	projects []project
	quitting bool
}

type refreshMsg struct{}

func loadProjects() []project {
	var projects []project
	pattern := filepath.Join(config.ProjectsDir(), "*", config.ConfFile)
	matches, _ := filepath.Glob(pattern)

	for _, match := range matches {
		dir := filepath.Dir(match)
		cfg, err := config.Load(dir)
		if err != nil {
			continue
		}

		status := "stopped"
		if s, err := vagrant.Status(dir); err == nil && s == "running" {
			status = "running"
		}

		projects = append(projects, project{dir: dir, config: cfg, status: status})
	}
	return projects
}

func buildTable(projects []project) table.Model {
	columns := []table.Column{
		{Title: "#", Width: 3},
		{Title: "Project", Width: 18},
		{Title: "Status", Width: 12},
		{Title: "Profiles", Width: 20},
		{Title: "Ports", Width: 30},
	}

	var rows []table.Row
	for i, p := range projects {
		statusIcon := "○ stopped"
		if p.status == "running" {
			statusIcon = "● running"
		}

		var portSummary []string
		for _, port := range p.config.Ports {
			portSummary = append(portSummary, fmt.Sprintf("%s:%d", port.Label, port.Host))
		}

		rows = append(rows, table.Row{
			fmt.Sprintf("%d", i+1),
			p.config.Name,
			statusIcon,
			strings.Join(p.config.Profiles, " "),
			strings.Join(portSummary, ", "),
		})
	}

	t := table.New(
		table.WithColumns(columns),
		table.WithRows(rows),
		table.WithFocused(true),
		table.WithHeight(len(rows)+1),
	)

	s := table.DefaultStyles()
	s.Header = s.Header.BorderStyle(lipgloss.NormalBorder()).BorderForeground(lipgloss.Color("240")).BorderBottom(true)
	s.Selected = s.Selected.Foreground(lipgloss.Color("229")).Background(lipgloss.Color("57")).Bold(false)
	t.SetStyles(s)

	return t
}

func newPsModel() psModel {
	projects := loadProjects()
	return psModel{
		table:    buildTable(projects),
		projects: projects,
	}
}

func (m psModel) Init() tea.Cmd {
	return nil
}

func (m psModel) selectedProject() *project {
	idx := m.table.Cursor()
	if idx >= 0 && idx < len(m.projects) {
		return &m.projects[idx]
	}
	return nil
}

func (m psModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case refreshMsg:
		m.projects = loadProjects()
		m.table = buildTable(m.projects)
		return m, nil

	case tea.KeyMsg:
		switch msg.String() {
		case "q", "esc":
			m.quitting = true
			return m, tea.Quit

		case "u":
			if p := m.selectedProject(); p != nil {
				c := exec.Command("vagrant", "up")
				c.Dir = p.dir
				return m, tea.ExecProcess(c, func(err error) tea.Msg {
					return refreshMsg{}
				})
			}

		case "d":
			if p := m.selectedProject(); p != nil {
				c := exec.Command("vagrant", "halt")
				c.Dir = p.dir
				return m, tea.ExecProcess(c, func(err error) tea.Msg {
					return refreshMsg{}
				})
			}

		case "s":
			if p := m.selectedProject(); p != nil {
				c := exec.Command("vagrant", "ssh")
				c.Dir = p.dir
				return m, tea.ExecProcess(c, func(err error) tea.Msg {
					return refreshMsg{}
				})
			}

		case "c":
			if p := m.selectedProject(); p != nil {
				c := exec.Command("vagrant", "ssh", "-c", "cd /vagrant && claude --dangerously-skip-permissions")
				c.Dir = p.dir
				return m, tea.ExecProcess(c, func(err error) tea.Msg {
					return refreshMsg{}
				})
			}

		case "D":
			if p := m.selectedProject(); p != nil {
				c := exec.Command("vagrant", "destroy", "-f")
				c.Dir = p.dir
				return m, tea.ExecProcess(c, func(err error) tea.Msg {
					return refreshMsg{}
				})
			}
		}
	}

	var cmd tea.Cmd
	m.table, cmd = m.table.Update(msg)
	return m, cmd
}

func (m psModel) View() string {
	if m.quitting {
		return ""
	}

	var b strings.Builder
	b.WriteString(m.table.View() + "\n\n")

	help := lipgloss.NewStyle().Faint(true)
	b.WriteString(help.Render("↑/↓ navigate • u start • d stop • s ssh • c claude • D destroy • q quit"))

	return b.String()
}

// RunPsDashboard runs the interactive ps dashboard.
func RunPsDashboard() error {
	p := tea.NewProgram(newPsModel())
	_, err := p.Run()
	return err
}

// PrintPlainPS prints a non-interactive table of projects.
func PrintPlainPS() {
	projects := loadProjects()
	fmt.Printf("%-20s %-20s %-30s  %s\n", "PROJECT", "PROFILES", "PORTS", "STATUS")
	fmt.Printf("%-20s %-20s %-30s  %s\n", "-------", "--------", "-----", "------")

	for _, p := range projects {
		var portSummary []string
		for _, port := range p.config.Ports {
			portSummary = append(portSummary, fmt.Sprintf("%s:%d", port.Label, port.Host))
		}
		fmt.Printf("%-20s %-20s %-30s  %s\n",
			p.config.Name,
			strings.Join(p.config.Profiles, " "),
			strings.Join(portSummary, ", "),
			p.status,
		)
	}
}
