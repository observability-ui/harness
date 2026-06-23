package ui

import (
	"fmt"
	"strings"

	"github.com/charmbracelet/bubbles/spinner"
	"github.com/charmbracelet/lipgloss"
	"obs/internal/recipe"
)

var failStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("1"))

type statusDef struct {
	Icon       string
	Style      lipgloss.Style
	UseSpinner bool
}

var statusDefs = map[recipe.Status]statusDef{
	recipe.StatusPending: {Icon: "○"},
	recipe.StatusWaiting: {Icon: "◷", Style: lipgloss.NewStyle().Foreground(lipgloss.Color("3"))},
	recipe.StatusRunning: {UseSpinner: true},
	recipe.StatusStarted: {UseSpinner: true},
	recipe.StatusReady:   {Icon: "●", Style: lipgloss.NewStyle().Foreground(lipgloss.Color("2"))},
	recipe.StatusDone:    {Icon: "✓", Style: lipgloss.NewStyle().Foreground(lipgloss.Color("2"))},
	recipe.StatusStopped: {Icon: "■", Style: lipgloss.NewStyle().Foreground(lipgloss.Color("8"))},
	recipe.StatusFailed:  {Icon: "✗", Style: failStyle},
	recipe.StatusSkipped: {Icon: "⊘"},
}

func StatusIcon(status recipe.Status, spinnerView string) string {
	def := statusDefs[status]
	if def.UseSpinner {
		return spinnerView
	}
	if def.Icon == "" {
		return ""
	}
	if def.Style.GetForeground() != (lipgloss.NoColor{}) {
		return def.Style.Render(def.Icon)
	}
	return def.Icon
}

type ProcessInfo struct {
	Name   string
	Status recipe.Status
}

type StepState struct {
	Name      string
	Status    recipe.Status
	Err       error
	Processes []ProcessInfo
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

func (mt *MainTab) AddStepWithProcesses(name string, processNames []string) {
	var procs []ProcessInfo
	for _, pn := range processNames {
		procs = append(procs, ProcessInfo{Name: pn, Status: recipe.StatusPending})
	}
	mt.steps = append(mt.steps, StepState{Name: name, Status: recipe.StatusPending, Processes: procs})
}

func (mt *MainTab) UpdateProcess(stepName, procName string, status recipe.Status) {
	step := mt.GetStep(stepName)
	if step == nil {
		return
	}
	for i := range step.Processes {
		if step.Processes[i].Name == procName {
			step.Processes[i].Status = status
			return
		}
	}
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
		lines = append(lines, fmt.Sprintf("  %s All requirements met", StatusIcon(recipe.StatusDone, "")))
		lines = append(lines, "")
	}

	if len(mt.steps) > 0 {
		lines = append(lines, lipgloss.NewStyle().Bold(true).Render("Recipe Steps:"))
		lines = append(lines, "")
	}

	treeStyle := lipgloss.NewStyle().Foreground(lipgloss.Color("8"))
	for _, step := range mt.steps {
		icon := StatusIcon(step.Status, mt.spinner.View())
		line := fmt.Sprintf("  %s %s", icon, step.Name)
		if step.Status == recipe.StatusFailed && step.Err != nil {
			line += failStyle.Render(fmt.Sprintf(" — %v", step.Err))
		}
		lines = append(lines, line)
		for _, proc := range step.Processes {
			procIcon := StatusIcon(proc.Status, mt.spinner.View())
			lines = append(lines, fmt.Sprintf("  %s %s %s", treeStyle.Render("│"), procIcon, proc.Name))
		}
	}

	return strings.Join(lines, "\n")
}
