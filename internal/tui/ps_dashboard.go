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
	"github.com/charmbracelet/bubbles/table"
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
type outputLineMsg string
type commandDoneMsg struct{ err error }
type refreshMsg struct{}

type psModel struct {
	table      table.Model
	projects   []project
	output     viewport.Model
	outputLog  []string
	running    bool
	cmdLabel   string
	width      int
	height     int
	quitting   bool
	tableStyle lipgloss.Style
	paneStyle  lipgloss.Style
}

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

func buildTable(projects []project, height int) table.Model {
	columns := []table.Column{
		{Title: "#", Width: 3},
		{Title: "Project", Width: 16},
		{Title: "Status", Width: 20},
		{Title: "Profiles", Width: 16},
	}

	greenStyle := lipgloss.NewStyle().Foreground(lipgloss.Color("42"))
	dimStyle := lipgloss.NewStyle().Faint(true)

	var rows []table.Row
	for i, p := range projects {
		statusIcon := dimStyle.Render("○ stopped")
		if p.status == "running" {
			statusIcon = greenStyle.Render("● running")
		}

		rows = append(rows, table.Row{
			fmt.Sprintf("%d", i+1),
			p.config.Name,
			statusIcon,
			strings.Join(p.config.Profiles, " "),
		})
	}

	h := len(rows) + 1
	if height > 3 && h > height-3 {
		h = height - 3
	}

	t := table.New(
		table.WithColumns(columns),
		table.WithRows(rows),
		table.WithFocused(true),
		table.WithHeight(h),
	)

	s := table.DefaultStyles()
	s.Header = s.Header.BorderStyle(lipgloss.NormalBorder()).BorderForeground(lipgloss.Color("240")).BorderBottom(true)
	s.Selected = s.Selected.Foreground(lipgloss.Color("229")).Background(lipgloss.Color("57")).Bold(false)
	t.SetStyles(s)

	return t
}

func newPsModel() psModel {
	projects := loadProjects()

	vp := viewport.New(40, 20)
	vp.SetContent("")

	return psModel{
		table:    buildTable(projects, 20),
		projects: projects,
		output:   vp,
		width:    120,
		height:   24,
		tableStyle: lipgloss.NewStyle().
			Border(lipgloss.RoundedBorder()).
			BorderForeground(lipgloss.Color("240")).
			Padding(0, 1),
		paneStyle: lipgloss.NewStyle().
			Border(lipgloss.RoundedBorder()).
			BorderForeground(lipgloss.Color("240")).
			Padding(0, 1),
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

// runCommand starts a command in the background, streaming output line by line.
func runCommand(dir string, args ...string) tea.Cmd {
	return func() tea.Msg {
		cmd := exec.Command(args[0], args[1:]...)
		cmd.Dir = dir

		// Combine stdout and stderr
		stdout, _ := cmd.StdoutPipe()
		cmd.Stderr = cmd.Stdout

		if err := cmd.Start(); err != nil {
			return commandDoneMsg{err: err}
		}

		// Read output in a goroutine — but we can't send tea.Msg from here
		// Instead, return the first approach: read all output via pipe
		// We'll use a different pattern: return a batch of commands
		go func() {
			cmd.Wait()
		}()

		// Actually, we need to stream. Let's use a different approach.
		// We'll read lines and the program will poll via channel.
		scanner := bufio.NewScanner(stdout)
		var lines []string
		for scanner.Scan() {
			lines = append(lines, scanner.Text())
		}
		cmd.Wait()

		// Return all lines at once (simpler than true streaming for now)
		return batchOutputMsg{lines: lines}
	}
}

type batchOutputMsg struct {
	lines []string
}

// streamCommand starts a command and streams output line by line via a channel.
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

func (m *psModel) recalcLayout() {
	// Left pane: table, Right pane: output
	// Account for borders (2 each side) and padding (1 each side) = 4 per pane, plus 1 gap
	leftWidth := 62
	rightWidth := m.width - leftWidth - 1
	if rightWidth < 20 {
		rightWidth = 20
	}

	paneHeight := m.height - 4 // room for help bar
	if paneHeight < 5 {
		paneHeight = 5
	}

	m.output.Width = rightWidth - 4
	m.output.Height = paneHeight - 2

	m.table = buildTable(m.projects, paneHeight-2)
}

func (m psModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height
		m.recalcLayout()
		return m, nil

	case batchOutputMsg:
		for _, line := range msg.lines {
			m.appendOutput(line)
		}
		m.running = false
		// Refresh project status after command completes
		m.projects = loadProjects()
		m.table = buildTable(m.projects, m.height-6)
		return m, nil

	case commandDoneMsg:
		m.running = false
		if msg.err != nil {
			m.appendOutput(fmt.Sprintf("Error: %v", msg.err))
		}
		m.projects = loadProjects()
		m.table = buildTable(m.projects, m.height-6)
		return m, nil

	case tea.KeyMsg:
		switch msg.String() {
		case "q", "esc":
			m.quitting = true
			return m, tea.Quit

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
				// SSH needs a real terminal — use ExecProcess
				c := exec.Command("vagrant", "ssh")
				c.Dir = p.dir
				return m, tea.ExecProcess(c, func(err error) tea.Msg {
					return refreshMsg{}
				})
			}

		case "c":
			if p := m.selectedProject(); p != nil {
				// Claude Code needs a real terminal — use ExecProcess
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

	// Handle refreshMsg (from ExecProcess returns for ssh/code)
	if _, ok := msg.(refreshMsg); ok {
		m.projects = loadProjects()
		m.table = buildTable(m.projects, m.height-6)
		return m, tea.ClearScreen
	}

	var cmd tea.Cmd
	m.table, cmd = m.table.Update(msg)
	return m, cmd
}

func (m psModel) View() string {
	if m.quitting {
		return ""
	}

	paneHeight := m.height - 4
	if paneHeight < 5 {
		paneHeight = 5
	}

	// Left pane: table
	leftWidth := 62
	left := m.tableStyle.
		Width(leftWidth - 4).
		Height(paneHeight - 2).
		Render(m.table.View())

	// Right pane: output
	rightWidth := m.width - leftWidth - 1
	if rightWidth < 20 {
		rightWidth = 20
	}

	outputTitle := "Output"
	if m.cmdLabel != "" {
		outputTitle = m.cmdLabel
		if m.running {
			outputTitle += " ..."
		}
	}

	titleStyle := lipgloss.NewStyle().Bold(true).Foreground(lipgloss.Color("42"))
	outputHeader := titleStyle.Render(outputTitle)

	rightContent := outputHeader + "\n" + m.output.View()
	right := m.paneStyle.
		Width(rightWidth - 4).
		Height(paneHeight - 2).
		Render(rightContent)

	// Join horizontally
	layout := lipgloss.JoinHorizontal(lipgloss.Top, left, " ", right)

	help := lipgloss.NewStyle().Faint(true)
	helpText := help.Render("↑/↓ navigate • u start • d stop • s ssh • c claude • D destroy • q quit")

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
