package verifyui

import (
	"fmt"
	"strings"
	"time"

	"github.com/charmbracelet/bubbles/spinner"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"

	"github.com/rohit746/tscgit/internal/gitutil"
	"github.com/rohit746/tscgit/internal/lessons"
	"github.com/rohit746/tscgit/internal/verify"
)

// Model drives the Bubble Tea verification view.
type Model struct {
	lesson  *lessons.Lesson
	repo    *gitutil.Repository
	spinner spinner.Model
	results []verify.Result
	width   int
	done    bool
}

type checkResultMsg struct {
	result verify.Result
}

// NewModel constructs the verification UI model.
func NewModel(lesson *lessons.Lesson, repo *gitutil.Repository) *Model {
	sp := spinner.New()
	sp.Spinner = spinner.Dot
	sp.Style = lipgloss.NewStyle().Foreground(lipgloss.Color("205"))

	return &Model{
		lesson:  lesson,
		repo:    repo,
		spinner: sp,
	}
}

// Init starts the spinner and kicks off the first check.
func (m *Model) Init() tea.Cmd {
	if len(m.lesson.Checks) == 0 {
		m.done = true
		return tea.Quit
	}
	return tea.Batch(m.spinner.Tick, runCheckCmd(m.lesson, m.repo, len(m.results)))
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
	case checkResultMsg:
		m.results = append(m.results, msg.result)
		if len(m.results) == len(m.lesson.Checks) {
			m.done = true
			return m, tea.Quit
		}
		return m, tea.Batch(m.spinner.Tick, runCheckCmd(m.lesson, m.repo, len(m.results)))
	default:
		return m, nil
	}
}

// View renders the verification UI.
func (m *Model) View() string {
	var b strings.Builder

	b.WriteString(titleStyle.Render(m.lesson.Title))
	b.WriteString("\n")
	b.WriteString(descStyle.Render(m.lesson.Description))
	b.WriteString("\n\n")

	for i, check := range m.lesson.Checks {
		var status string
		var detail string

		if i < len(m.results) {
			res := m.results[i]
			passed := res.Outcome.Passed && res.Outcome.Err == nil
			if res.Outcome.Err != nil {
				detail = res.Outcome.Err.Error()
			} else if res.Outcome.Message != "" {
				detail = res.Outcome.Message
			} else {
				detail = "Completed."
			}

			if passed {
				status = passStyle.Render(fmt.Sprintf("%s %s", successGlyph, check.Title))
				detail = detailStyle.Render(detail)
			} else {
				status = failStyle.Render(fmt.Sprintf("%s %s", failGlyph, check.Title))
				if res.Outcome.Err != nil {
					detail = failDetailStyle.Render(detail)
				} else {
					detail = warnDetailStyle.Render(detail)
				}
			}

			duration := lipgloss.NewStyle().Faint(true).Render(res.Duration.Round(10 * time.Millisecond).String())
			b.WriteString(status)
			b.WriteString(" " + duration)
		} else {
			if !m.done && i == len(m.results) {
				status = pendingStyle.Render(fmt.Sprintf("%s %s", m.spinner.View(), check.Title))
				detail = detailStyle.Render(check.Description)
			} else {
				status = pendingStyle.Render(fmt.Sprintf("%s %s", pendingGlyph, check.Title))
				detail = detailStyle.Render(check.Description)
			}
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
		if res.Outcome.Passed && res.Outcome.Err == nil {
			passed++
		}
	}
	total := len(m.lesson.Checks)

	summary := fmt.Sprintf("%d/%d checks passed", passed, total)
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
		b.WriteString(helpStyle.Render("Verifying… press q to cancel."))
	}

	return lipgloss.NewStyle().Width(m.width).Render(strings.TrimSuffix(b.String(), "\n"))
}

func runCheckCmd(lesson *lessons.Lesson, repo *gitutil.Repository, index int) tea.Cmd {
	if index >= len(lesson.Checks) {
		return nil
	}
	check := lesson.Checks[index]

	return func() tea.Msg {
		ctx, cancel := lessons.TimeoutContext()
		defer cancel()

		start := time.Now()
		outcome := check.Verify(ctx, repo)
		res := verify.Result{
			Check:    check,
			Outcome:  outcome,
			Duration: time.Since(start),
		}
		return checkResultMsg{result: res}
	}
}

// Results returns a snapshot of the check outcomes collected so far.
func (m *Model) Results() []verify.Result {
	out := make([]verify.Result, len(m.results))
	copy(out, m.results)
	return out
}

// AllPassed reports whether every check finished successfully.
func (m *Model) AllPassed() bool {
	if len(m.results) != len(m.lesson.Checks) {
		return false
	}
	for _, res := range m.results {
		if !res.Outcome.Passed || res.Outcome.Err != nil {
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
	warnDetailStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("214"))
	failDetailStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("203")).Bold(true)

	summaryPassStyle    = lipgloss.NewStyle().Foreground(lipgloss.Color("84")).Bold(true)
	summaryFailStyle    = lipgloss.NewStyle().Foreground(lipgloss.Color("203")).Bold(true)
	summaryPendingStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("244")).Bold(true)
	helpStyle           = lipgloss.NewStyle().Foreground(lipgloss.Color("244")).Faint(true)
)
