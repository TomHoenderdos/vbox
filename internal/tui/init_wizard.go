package tui

import (
	"fmt"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/TomHoenderdos/vbox/internal/config"
	"github.com/TomHoenderdos/vbox/internal/profile"
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type initStep int

const (
	stepName initStep = iota
	stepProfiles
	stepPorts
	stepMemory
	stepCPUs
	stepAutoSync
	stepSummary
)

type initModel struct {
	step       initStep
	input      textinput.Model
	cwd        string
	config     *config.Config
	profiles   []profile.Info
	portIdx    int
	ports      []config.Port
	aborted    bool
	done       bool
	titleStyle lipgloss.Style
	dimStyle   lipgloss.Style
	boxStyle   lipgloss.Style
}

func newInitModel(cwd string) initModel {
	ti := textinput.New()
	ti.Focus()
	ti.Placeholder = filepath.Base(cwd)

	profiles, _ := profile.List()

	return initModel{
		step:       stepName,
		input:      ti,
		cwd:        cwd,
		config:     &config.Config{Memory: 2048, CPUs: 2, AutoSync: true},
		profiles:   profiles,
		titleStyle: lipgloss.NewStyle().Bold(true),
		dimStyle:   lipgloss.NewStyle().Faint(true),
		boxStyle:   lipgloss.NewStyle().Border(lipgloss.RoundedBorder()).Padding(0, 1),
	}
}

func (m initModel) Init() tea.Cmd {
	return textinput.Blink
}

func (m initModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c", "esc":
			m.aborted = true
			return m, tea.Quit
		case "enter":
			return m.advance()
		}
	}

	var cmd tea.Cmd
	m.input, cmd = m.input.Update(msg)
	return m, cmd
}

func (m initModel) advance() (tea.Model, tea.Cmd) {
	val := strings.TrimSpace(m.input.Value())

	switch m.step {
	case stepName:
		name := val
		if name == "" {
			name = filepath.Base(m.cwd)
		}
		m.config.Name = name
		m.step = stepProfiles
		m.input.SetValue("")
		m.input.Placeholder = "elixir"

	case stepProfiles:
		input := val
		if input == "" {
			input = "elixir"
		}
		selections := strings.Split(input, ",")
		var selected []string
		for _, sel := range selections {
			sel = strings.TrimSpace(sel)
			if num, err := strconv.Atoi(sel); err == nil {
				idx := num - 1
				if idx >= 0 && idx < len(m.profiles) {
					selected = append(selected, m.profiles[idx].Name)
				}
			} else {
				if profile.Exists(sel) {
					selected = append(selected, sel)
				}
			}
		}
		if len(selected) == 0 {
			selected = []string{"elixir"}
		}
		m.config.Profiles = selected

		m.ports, _ = profile.CollectPorts(selected)
		if len(m.ports) > 0 {
			m.step = stepPorts
			m.portIdx = 0
			m.input.SetValue("")
			m.input.Placeholder = strconv.Itoa(m.ports[0].Host)
		} else {
			m.step = stepMemory
			m.input.SetValue("")
			m.input.Placeholder = "2048"
		}

	case stepPorts:
		host := m.ports[m.portIdx].Host
		if val != "" {
			if v, err := strconv.Atoi(val); err == nil {
				host = v
			}
		}
		m.config.Ports = append(m.config.Ports, config.Port{
			Guest: m.ports[m.portIdx].Guest,
			Host:  host,
			Label: m.ports[m.portIdx].Label,
		})
		m.portIdx++
		if m.portIdx < len(m.ports) {
			m.input.SetValue("")
			m.input.Placeholder = strconv.Itoa(m.ports[m.portIdx].Host)
		} else {
			m.step = stepMemory
			m.input.SetValue("")
			m.input.Placeholder = "2048"
		}

	case stepMemory:
		mem := 2048
		if val != "" {
			if v, err := strconv.Atoi(val); err == nil {
				mem = v
			}
		}
		m.config.Memory = mem
		m.step = stepCPUs
		m.input.SetValue("")
		m.input.Placeholder = "2"

	case stepCPUs:
		cpus := 2
		if val != "" {
			if v, err := strconv.Atoi(val); err == nil {
				cpus = v
			}
		}
		m.config.CPUs = cpus
		m.step = stepAutoSync
		m.input.SetValue("")
		m.input.Placeholder = "Y"

	case stepAutoSync:
		if strings.ToLower(val) == "n" {
			m.config.AutoSync = false
		}
		m.step = stepSummary
		m.input.SetValue("")
		m.input.Placeholder = "Y"

	case stepSummary:
		if strings.ToLower(val) == "n" {
			m.aborted = true
			return m, tea.Quit
		}
		m.done = true
		return m, tea.Quit
	}

	return m, nil
}

func (m initModel) View() string {
	var b strings.Builder

	b.WriteString("vbox init - interactive setup\n\n")

	switch m.step {
	case stepName:
		b.WriteString(m.titleStyle.Render("Project name") + "\n")
		b.WriteString(m.input.View() + "\n")
		b.WriteString(m.dimStyle.Render("(default: " + filepath.Base(m.cwd) + ")"))

	case stepProfiles:
		b.WriteString(m.titleStyle.Render("Select profiles") + "\n\n")
		for i, p := range m.profiles {
			b.WriteString(fmt.Sprintf("  %2d) %-12s %s\n", i+1, p.Name, p.Description))
		}
		b.WriteString("\n" + m.input.View() + "\n")
		b.WriteString(m.dimStyle.Render("(comma-separated names or numbers)"))

	case stepPorts:
		p := m.ports[m.portIdx]
		b.WriteString(m.titleStyle.Render("Port configuration") + "\n\n")
		b.WriteString(fmt.Sprintf("  %s (guest :%d) -> host port:\n", p.Label, p.Guest))
		b.WriteString(m.input.View() + "\n")
		b.WriteString(m.dimStyle.Render(fmt.Sprintf("(default: %d)", p.Host)))

	case stepMemory:
		b.WriteString(m.titleStyle.Render("Memory (MB)") + "\n")
		b.WriteString(m.input.View() + "\n")
		b.WriteString(m.dimStyle.Render("(default: 2048)"))

	case stepCPUs:
		b.WriteString(m.titleStyle.Render("CPUs") + "\n")
		b.WriteString(m.input.View() + "\n")
		b.WriteString(m.dimStyle.Render("(default: 2)"))

	case stepAutoSync:
		b.WriteString(m.titleStyle.Render("Auto file sync?") + "\n")
		b.WriteString(m.input.View() + "\n")
		b.WriteString(m.dimStyle.Render("(Y/n)"))

	case stepSummary:
		var summary strings.Builder
		summary.WriteString(fmt.Sprintf("  Project:    %s\n", m.config.Name))
		summary.WriteString(fmt.Sprintf("  Profiles:   %s\n", strings.Join(m.config.Profiles, " ")))
		for _, p := range m.config.Ports {
			summary.WriteString(fmt.Sprintf("  %-12s :%-5d -> :%-5d\n", p.Label, p.Guest, p.Host))
		}
		summary.WriteString(fmt.Sprintf("  Memory:     %dMB / %d CPUs\n", m.config.Memory, m.config.CPUs))
		summary.WriteString(fmt.Sprintf("  Auto sync:  %t", m.config.AutoSync))

		b.WriteString(m.titleStyle.Render("Summary") + "\n")
		b.WriteString(m.boxStyle.Render(summary.String()) + "\n\n")
		b.WriteString("Continue? " + m.input.View() + "\n")
		b.WriteString(m.dimStyle.Render("(Y/n)"))
	}

	return b.String()
}

// RunInitWizard runs the interactive init wizard and returns the resulting config.
func RunInitWizard(cwd string) (*config.Config, error) {
	m := newInitModel(cwd)
	p := tea.NewProgram(m)
	result, err := p.Run()
	if err != nil {
		return nil, err
	}

	final := result.(initModel)
	if final.aborted {
		return nil, fmt.Errorf("aborted")
	}
	return final.config, nil
}
