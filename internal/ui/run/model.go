package runui

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/charmbracelet/bubbles/spinner"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"

	"github.com/rohit746/tscgit/internal/run"
)

// Model drives the Bubble Tea UI for scripted run lessons.
type Model struct {
	script  *run.Script
	spinner spinner.Model
	results []run.StepResult
	done    bool
	width   int
}

type stepResultMsg struct {
	result run.StepResult
}

// NewModel creates a run lesson UI model.
func NewModel(script *run.Script) *Model {
	sp := spinner.New()
	sp.Spinner = spinner.Dot
	sp.Style = lipgloss.NewStyle().Foreground(lipgloss.Color("205"))
	return &Model{script: script, spinner: sp}
}

// Init starts the spinner and kicks off the first command.
func (m *Model) Init() tea.Cmd {
	if len(m.script.Steps) == 0 {
		m.done = true
		return tea.Quit
	}
	return tea.Batch(m.spinner.Tick, runStepCmd(m.script, 0))
}

// Update handles Bubble Tea messages.
func (m *Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.width = msg.Width
		return m, nil
	case spinner.TickMsg:
		if m.done {
			return m, nil
		}
		var cmd tea.Cmd
		m.spinner, cmd = m.spinner.Update(msg)
		return m, cmd
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c", "q", "esc":
			return m, tea.Quit
		case "enter":
			if m.done {
				return m, tea.Quit
			}
		}
		return m, nil
	case stepResultMsg:
		m.results = append(m.results, msg.result)
		if len(m.results) == len(m.script.Steps) {
			m.done = true
			return m, tea.Quit
		}
		return m, tea.Batch(m.spinner.Tick, runStepCmd(m.script, len(m.results)))
	default:
		return m, nil
	}
}

// View renders the UI.
func (m *Model) View() string {
	var b strings.Builder

	b.WriteString(titleStyle.Render(fmt.Sprintf("Run Lesson %s — %s", m.script.ID, m.script.Title)))
	b.WriteString("\n")
	if m.script.Description != "" {
		b.WriteString(descStyle.Render(m.script.Description))
		b.WriteString("\n\n")
	}

	for i, step := range m.script.Steps {
		var status string
		var detail string

		if i < len(m.results) {
			res := m.results[i]
			if res.Passed {
				status = passStyle.Render(fmt.Sprintf("%s %s", successGlyph, step.Command))
				detail = detailStyle.Render(summarizeStdout(res.Stdout))
			} else {
				status = failStyle.Render(fmt.Sprintf("%s %s", failGlyph, step.Command))
				failures := make([]string, len(res.Failures))
				copy(failures, res.Failures)
				detail = failDetailStyle.Render(strings.Join(failures, "; "))
			}
			timing := lipgloss.NewStyle().Faint(true).Render(res.Duration.Round(10 * time.Millisecond).String())
			b.WriteString(status)
			b.WriteString(" " + timing)
		} else {
			if !m.done && i == len(m.results) {
				status = pendingStyle.Render(fmt.Sprintf("%s %s", m.spinner.View(), step.Command))
			} else {
				status = pendingStyle.Render(fmt.Sprintf("%s %s", pendingGlyph, step.Command))
			}
			detail = detailStyle.Render(stepDescription(step))
			b.WriteString(status)
		}

		if detail != "" {
			b.WriteString("\n  ")
			b.WriteString(detail)
		}
		b.WriteString("\n\n")
	}

	passed := 0
	for _, res := range m.results {
		if res.Passed {
			passed++
		}
	}
	total := len(m.script.Steps)

	summary := fmt.Sprintf("%d/%d steps passed", passed, total)
	if m.done {
		if passed == total {
			b.WriteString(summaryPassStyle.Render(summary))
		} else {
			b.WriteString(summaryFailStyle.Render(summary))
		}
		b.WriteString("\n")
		b.WriteString(helpStyle.Render("Press Enter or q to exit."))
	} else {
		b.WriteString(summaryPendingStyle.Render(summary))
		b.WriteString("\n")
		b.WriteString(helpStyle.Render("Running… press q to cancel."))
	}

	return lipgloss.NewStyle().Width(m.width).Render(strings.TrimSuffix(b.String(), "\n"))
}

func summarizeStdout(out string) string {
	trimmed := strings.TrimSpace(out)
	if trimmed == "" {
		return "(no stdout)"
	}
	if len([]rune(trimmed)) > 200 {
		trimmed = string([]rune(trimmed)[:200]) + "…"
	}
	return trimmed
}

func stepDescription(step run.Step) string {
	if len(step.ExpectStdout) == 0 {
		return "Awaiting command completion."
	}
	return fmt.Sprintf("Expecting stdout to include: %s", strings.Join(step.ExpectStdout, ", "))
}

func runStepCmd(script *run.Script, index int) tea.Cmd {
	if index >= len(script.Steps) {
		return nil
	}
	step := script.Steps[index]
	return func() tea.Msg {
		ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
		defer cancel()
		result := run.RunStep(ctx, step)
		return stepResultMsg{result: result}
	}
}

// Results returns the collected step outcomes.
func (m *Model) Results() []run.StepResult {
	out := make([]run.StepResult, len(m.results))
	copy(out, m.results)
	return out
}

// AllPassed reports whether every step succeeded.
func (m *Model) AllPassed() bool {
	if len(m.results) != len(m.script.Steps) {
		return false
	}
	for _, res := range m.results {
		if !res.Passed {
			return false
		}
	}
	return true
}

var (
	successGlyph = lipgloss.NewStyle().Foreground(lipgloss.Color("84")).Render("✔")
	failGlyph    = lipgloss.NewStyle().Foreground(lipgloss.Color("203")).Render("✘")
	pendingGlyph = lipgloss.NewStyle().Foreground(lipgloss.Color("248")).Render("•")

	titleStyle = lipgloss.NewStyle().Bold(true).Foreground(lipgloss.Color("147")).MarginBottom(1)
	descStyle  = lipgloss.NewStyle().Foreground(lipgloss.Color("250"))

	passStyle    = lipgloss.NewStyle().Foreground(lipgloss.Color("84")).Bold(true)
	failStyle    = lipgloss.NewStyle().Foreground(lipgloss.Color("203")).Bold(true)
	pendingStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("244")).Bold(true)

	detailStyle     = lipgloss.NewStyle().Foreground(lipgloss.Color("245"))
	failDetailStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("203")).Bold(true)

	summaryPassStyle    = lipgloss.NewStyle().Foreground(lipgloss.Color("84")).Bold(true)
	summaryFailStyle    = lipgloss.NewStyle().Foreground(lipgloss.Color("203")).Bold(true)
	summaryPendingStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("244")).Bold(true)
	helpStyle           = lipgloss.NewStyle().Foreground(lipgloss.Color("244")).Faint(true)
)
