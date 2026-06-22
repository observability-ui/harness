package ui

import (
	"fmt"
	"strings"

	"github.com/charmbracelet/bubbles/spinner"
	"github.com/charmbracelet/lipgloss"
	"obs/internal/recipe"
)

var (
	waitStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("3")) // yellow
	doneStyle    = lipgloss.NewStyle().Foreground(lipgloss.Color("2")) // green
	stoppedStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("8")) // grey
	failStyle    = lipgloss.NewStyle().Foreground(lipgloss.Color("1")) // red
	readyStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("2")) // green
	statusIcons = map[recipe.Status]string{
		recipe.StatusPending: "○",
		recipe.StatusWaiting: "◷", // waiting on dependency
		recipe.StatusRunning: "",  // replaced by spinner
		recipe.StatusStarted: "",  // replaced by spinner (running process)
		recipe.StatusReady:   "●", // ports accepting connections
		recipe.StatusDone:    "✓", // completed and exited
		recipe.StatusStopped: "■", // stopped by user
		recipe.StatusFailed:  "✗",
		recipe.StatusSkipped: "⊘",
	}
)

func StatusIcon(status recipe.Status, spinnerView string) string {
	switch status {
	case recipe.StatusWaiting:
		return waitStyle.Render(statusIcons[status])
	case recipe.StatusRunning, recipe.StatusStarted:
		return spinnerView
	case recipe.StatusReady:
		return readyStyle.Render(statusIcons[status])
	case recipe.StatusDone:
		return doneStyle.Render(statusIcons[status])
	case recipe.StatusStopped:
		return stoppedStyle.Render(statusIcons[status])
	case recipe.StatusFailed:
		return failStyle.Render(statusIcons[status])
	default:
		return statusIcons[status]
	}
}

type StepState struct {
	Name   string
	Status recipe.Status
	Err    error
}

type MainTab struct {
	steps   []StepState
	spinner spinner.Model
}

func NewMainTab() MainTab {
	s := spinner.New()
	s.Spinner = spinner.Dot
	s.Style = lipgloss.NewStyle().Foreground(lipgloss.Color("6"))
	return MainTab{spinner: s}
}

func (mt *MainTab) AddStep(name string) {
	mt.steps = append(mt.steps, StepState{Name: name, Status: recipe.StatusPending})
}

func (mt *MainTab) GetStep(name string) *StepState {
	for i := range mt.steps {
		if mt.steps[i].Name == name {
			return &mt.steps[i]
		}
	}
	return nil
}

func (mt *MainTab) UpdateStep(name string, status recipe.Status, err error) {
	for i := range mt.steps {
		if mt.steps[i].Name == name {
			mt.steps[i].Status = status
			mt.steps[i].Err = err
			return
		}
	}
}

func (mt MainTab) View(width, height int) string {
	return mt.ViewWithRequirements(width, height, 1, nil) // reqPassed = 1
}

func (mt MainTab) ViewWithRequirements(width, height int, reqStatus reqState, reqErr error) string {
	var lines []string

	// Requirements section
	switch reqStatus {
	case reqChecking:
		lines = append(lines, lipgloss.NewStyle().Bold(true).Render("Requirements:"))
		lines = append(lines, fmt.Sprintf("  %s Checking requirements…", mt.spinner.View()))
		lines = append(lines, "")
	case reqFailed:
		lines = append(lines, lipgloss.NewStyle().Bold(true).Render("Requirements:"))
		lines = append(lines, fmt.Sprintf("  %s %s", failStyle.Render("✗"), failStyle.Render(reqErr.Error())))
		lines = append(lines, "")
	case reqPassed:
		lines = append(lines, lipgloss.NewStyle().Bold(true).Render("Requirements:"))
		lines = append(lines, fmt.Sprintf("  %s All requirements met", doneStyle.Render("✓")))
		lines = append(lines, "")
	}

	if len(mt.steps) > 0 {
		lines = append(lines, lipgloss.NewStyle().Bold(true).Render("Recipe Steps:"))
		lines = append(lines, "")
	}

	for _, step := range mt.steps {
		icon := StatusIcon(step.Status, mt.spinner.View())
		line := fmt.Sprintf("  %s %s", icon, step.Name)
		if step.Status == recipe.StatusFailed && step.Err != nil {
			line += failStyle.Render(fmt.Sprintf(" — %v", step.Err))
		}
		lines = append(lines, line)
	}

	return strings.Join(lines, "\n")
}
