package ui

import (
	"fmt"
	"strings"

	"github.com/charmbracelet/bubbles/spinner"
	"github.com/charmbracelet/lipgloss"
	"obsui/internal/recipe"
)

var (
	doneStyle   = lipgloss.NewStyle().Foreground(lipgloss.Color("2")) // green
	failStyle   = lipgloss.NewStyle().Foreground(lipgloss.Color("1")) // red
	statusIcons = map[recipe.Status]string{
		recipe.StatusPending: "○",
		recipe.StatusRunning: "", // replaced by spinner
		recipe.StatusDone:    "✓",
		recipe.StatusFailed:  "✗",
		recipe.StatusSkipped: "⊘",
	}
)

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
	var lines []string
	lines = append(lines, lipgloss.NewStyle().Bold(true).Render("Recipe Steps:"))
	lines = append(lines, "")

	for _, step := range mt.steps {
		icon := statusIcons[step.Status]
		line := ""
		switch step.Status {
		case recipe.StatusRunning:
			line = fmt.Sprintf("  %s %s", mt.spinner.View(), step.Name)
		case recipe.StatusDone:
			line = fmt.Sprintf("  %s %s", doneStyle.Render(icon), step.Name)
		case recipe.StatusFailed:
			line = fmt.Sprintf("  %s %s", failStyle.Render(icon), step.Name)
			if step.Err != nil {
				line += failStyle.Render(fmt.Sprintf(" — %v", step.Err))
			}
		default:
			line = fmt.Sprintf("  %s %s", icon, step.Name)
		}
		lines = append(lines, line)
	}

	return strings.Join(lines, "\n")
}
