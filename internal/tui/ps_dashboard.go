package tui

import (
	"bufio"
	"fmt"
	"io"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/TomHoenderdos/vbox/internal/config"
	"github.com/TomHoenderdos/vbox/internal/vagrant"
	"github.com/charmbracelet/bubbles/viewport"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type project struct {
	dir    string
	config *config.Config
	status string
}

// Messages
type batchOutputMsg struct{ lines []string }
type commandDoneMsg struct{ err error }
type refreshMsg struct{}

type psModel struct {
	projects  []project
	cursor    int
	output    viewport.Model
	outputLog []string
	running   bool
	cmdLabel  string
	width     int
	height    int
	quitting  bool
}

var (
	greenDot  = lipgloss.NewStyle().Foreground(lipgloss.Color("42")).Render("●")
	dimDot    = lipgloss.NewStyle().Faint(true).Render("○")
	greenText = lipgloss.NewStyle().Foreground(lipgloss.Color("42"))
	dimText   = lipgloss.NewStyle().Faint(true)
	boldText  = lipgloss.NewStyle().Bold(true)
	headerStyle = lipgloss.NewStyle().Bold(true).Foreground(lipgloss.Color("240"))
	selectedStyle = lipgloss.NewStyle().Background(lipgloss.Color("57")).Foreground(lipgloss.Color("229"))
	paneStyle = lipgloss.NewStyle().
			Border(lipgloss.RoundedBorder()).
			BorderForeground(lipgloss.Color("240")).
			Padding(0, 1)
	helpStyle = lipgloss.NewStyle().Faint(true)
)

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

func newPsModel() psModel {
	projects := loadProjects()
	vp := viewport.New(40, 20)
	vp.SetContent("")

	return psModel{
		projects: projects,
		output:   vp,
		width:    120,
		height:   24,
	}
}

func (m psModel) Init() tea.Cmd {
	return nil
}

func (m psModel) selectedProject() *project {
	if m.cursor >= 0 && m.cursor < len(m.projects) {
		return &m.projects[m.cursor]
	}
	return nil
}

func streamCommand(dir string, args ...string) tea.Cmd {
	return func() tea.Msg {
		cmd := exec.Command(args[0], args[1:]...)
		cmd.Dir = dir

		pr, pw := io.Pipe()
		cmd.Stdout = pw
		cmd.Stderr = pw

		if err := cmd.Start(); err != nil {
			return commandDoneMsg{err: err}
		}

		go func() {
			cmd.Wait()
			pw.Close()
		}()

		scanner := bufio.NewScanner(pr)
		var lines []string
		for scanner.Scan() {
			lines = append(lines, scanner.Text())
		}

		return batchOutputMsg{lines: lines}
	}
}

func (m *psModel) appendOutput(line string) {
	m.outputLog = append(m.outputLog, line)
	content := strings.Join(m.outputLog, "\n")
	m.output.SetContent(content)
	m.output.GotoBottom()
}

func (m *psModel) clearOutput() {
	m.outputLog = nil
	m.output.SetContent("")
}

func (m psModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height
		rightWidth := m.width - leftPaneWidth(m.width) - 5
		if rightWidth < 10 {
			rightWidth = 10
		}
		paneH := m.height - 4
		if paneH < 3 {
			paneH = 3
		}
		m.output.Width = rightWidth
		m.output.Height = paneH - 4
		return m, nil

	case batchOutputMsg:
		for _, line := range msg.lines {
			m.appendOutput(line)
		}
		m.running = false
		m.projects = loadProjects()
		return m, nil

	case commandDoneMsg:
		m.running = false
		if msg.err != nil {
			m.appendOutput(fmt.Sprintf("Error: %v", msg.err))
		}
		m.projects = loadProjects()
		return m, nil

	case refreshMsg:
		m.projects = loadProjects()
		return m, tea.ClearScreen

	case tea.KeyMsg:
		switch msg.String() {
		case "q", "esc", "ctrl+c":
			m.quitting = true
			return m, tea.Quit

		case "up", "k":
			if m.cursor > 0 {
				m.cursor--
			}
			return m, nil

		case "down", "j":
			if m.cursor < len(m.projects)-1 {
				m.cursor++
			}
			return m, nil

		case "u":
			if p := m.selectedProject(); p != nil && !m.running {
				m.running = true
				m.cmdLabel = "vagrant up: " + p.config.Name
				m.clearOutput()
				m.appendOutput("$ vagrant up (" + p.config.Name + ")")
				m.appendOutput("")
				return m, streamCommand(p.dir, "vagrant", "up")
			}

		case "d":
			if p := m.selectedProject(); p != nil && !m.running {
				m.running = true
				m.cmdLabel = "vagrant halt: " + p.config.Name
				m.clearOutput()
				m.appendOutput("$ vagrant halt (" + p.config.Name + ")")
				m.appendOutput("")
				return m, streamCommand(p.dir, "vagrant", "halt")
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
			if p := m.selectedProject(); p != nil && !m.running {
				m.running = true
				m.cmdLabel = "vagrant destroy: " + p.config.Name
				m.clearOutput()
				m.appendOutput("$ vagrant destroy -f (" + p.config.Name + ")")
				m.appendOutput("")
				return m, streamCommand(p.dir, "vagrant", "destroy", "-f")
			}
		}
	}

	return m, nil
}

func leftPaneWidth(termWidth int) int {
	w := termWidth / 2
	if w > 60 {
		w = 60
	}
	if w < 40 {
		w = 40
	}
	return w
}

func (m psModel) renderTable(width, height int) string {
	var b strings.Builder

	// Header
	header := fmt.Sprintf(" %-3s %-14s %-12s %s", "#", "Project", "Status", "Profiles")
	if len(header) > width {
		header = header[:width]
	}
	b.WriteString(headerStyle.Render(header) + "\n")
	b.WriteString(headerStyle.Render(strings.Repeat("─", width-2)) + "\n")

	// Rows
	for i, p := range m.projects {
		num := fmt.Sprintf("%-3d", i+1)
		name := truncate(p.config.Name, 14)
		profiles := truncate(strings.Join(p.config.Profiles, " "), width-35)

		var status string
		if p.status == "running" {
			status = greenDot + " " + greenText.Render("running")
		} else {
			status = dimDot + " " + dimText.Render("stopped")
		}

		// Pad name and profiles to fixed widths (using visible chars only)
		line := fmt.Sprintf("  %s %-14s %s  %s", num, name, padRight(status, 12, 10), profiles)

		if i == m.cursor {
			prefix := fmt.Sprintf("  %-3d %-14s ", i+1, name)
			suffix := fmt.Sprintf("  %s", profiles)
			line = selectedStyle.Render(prefix) + padRight(status, 12, 10) + selectedStyle.Render(padToWidth(suffix, width-2-len(prefix)-12))
		}

		b.WriteString(line + "\n")
	}

	// Fill remaining height
	rendered := strings.Count(b.String(), "\n")
	for rendered < height {
		b.WriteString("\n")
		rendered++
	}

	return b.String()
}

// padRight pads a string that contains ANSI codes.
// visibleLen is the known visible character count of s.
// targetVisible is the desired visible width.
func padRight(s string, targetVisible, visibleLen int) string {
	if visibleLen >= targetVisible {
		return s
	}
	return s + strings.Repeat(" ", targetVisible-visibleLen)
}

func padToWidth(s string, width int) string {
	if len(s) >= width {
		return s[:width]
	}
	return s + strings.Repeat(" ", width-len(s))
}

func truncate(s string, maxLen int) string {
	if len(s) <= maxLen {
		return s
	}
	if maxLen <= 3 {
		return s[:maxLen]
	}
	return s[:maxLen-1] + "…"
}

func (m psModel) View() string {
	if m.quitting {
		return ""
	}

	paneH := m.height - 4
	if paneH < 5 {
		paneH = 5
	}

	leftW := leftPaneWidth(m.width)
	innerLeftW := leftW - 4 // border + padding

	// Left pane: custom table
	tableContent := m.renderTable(innerLeftW, paneH-2)
	left := paneStyle.
		Width(innerLeftW).
		Height(paneH - 2).
		Render(tableContent)

	// Right pane: output
	rightW := m.width - leftW - 5
	if rightW < 10 {
		rightW = 10
	}

	outputTitle := "Output"
	if m.cmdLabel != "" {
		outputTitle = m.cmdLabel
		if m.running {
			outputTitle += " ..."
		}
	}

	m.output.Width = rightW
	m.output.Height = paneH - 4
	outputHeader := boldText.Copy().Foreground(lipgloss.Color("42")).Render(outputTitle)
	rightContent := outputHeader + "\n" + m.output.View()

	right := paneStyle.
		Width(rightW).
		Height(paneH - 2).
		Render(rightContent)

	layout := lipgloss.JoinHorizontal(lipgloss.Top, left, " ", right)

	helpText := helpStyle.Render("↑/↓ navigate • u start • d stop • s ssh • c claude • D destroy • q quit")

	return layout + "\n" + helpText
}

// RunPsDashboard runs the interactive ps dashboard.
func RunPsDashboard() error {
	p := tea.NewProgram(newPsModel(), tea.WithAltScreen())
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
