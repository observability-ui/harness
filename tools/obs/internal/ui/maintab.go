package ui

import (
	"fmt"
	"strings"

	"obs/internal/component"

	"github.com/charmbracelet/bubbles/spinner"
	"github.com/charmbracelet/bubbles/viewport"
	"github.com/charmbracelet/lipgloss"
)

var (
	dimStyle  = lipgloss.NewStyle().Foreground(lipgloss.Color("8"))
	failStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("1"))
)

type statusDef struct {
	Icon       string
	Style      lipgloss.Style
	UseSpinner bool
}

var statusDefs = map[component.Status]statusDef{
	component.StatusPending: {Icon: "○"},
	component.StatusWaiting: {Icon: "◷", Style: lipgloss.NewStyle().Foreground(lipgloss.Color("3"))},
	component.StatusRunning: {UseSpinner: true},
	component.StatusStarted: {UseSpinner: true},
	component.StatusReady:   {Icon: "●", Style: lipgloss.NewStyle().Foreground(lipgloss.Color("2"))},
	component.StatusDone:    {Icon: "✓", Style: lipgloss.NewStyle().Foreground(lipgloss.Color("2"))},
	component.StatusStopped: {Icon: "■", Style: lipgloss.NewStyle().Foreground(lipgloss.Color("8"))},
	component.StatusFailed:  {Icon: "✗", Style: failStyle},
	component.StatusSkipped: {Icon: "⊘"},
}

func StatusIcon(status component.Status, spinnerView string) string {
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
	Status component.Status
}

type StepState struct {
	Name      string
	Status    component.Status
	Err       error
	Processes []ProcessInfo
	DependsOn []string
}

type MainTab struct {
	steps    []StepState
	spinner  spinner.Model
	viewport viewport.Model
}

func NewMainTab() MainTab {
	s := spinner.New()
	s.Spinner = spinner.Dot
	s.Style = lipgloss.NewStyle().Foreground(lipgloss.Color("6"))
	return MainTab{spinner: s}
}

func (mt *MainTab) SetSize(width, height int) {
	mt.viewport.Width = width
	mt.viewport.Height = height
}

func (mt *MainTab) AddStepWithProcesses(name string, processNames []string, dependsOn []string) {
	var procs []ProcessInfo
	for _, pn := range processNames {
		procs = append(procs, ProcessInfo{Name: pn, Status: component.StatusPending})
	}
	mt.steps = append(mt.steps, StepState{Name: name, Status: component.StatusPending, Processes: procs, DependsOn: dependsOn})
}

func (mt *MainTab) UpdateProcess(stepName, procName string, status component.Status) {
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

func (mt *MainTab) UpdateStep(name string, status component.Status, err error) {
	for i := range mt.steps {
		if mt.steps[i].Name == name {
			mt.steps[i].Status = status
			mt.steps[i].Err = err
			return
		}
	}
}

func (mt *MainTab) RefreshContent(width int, reqStatus reqState, reqErr error) {
	var lines []string

	switch reqStatus {
	case reqChecking:
		lines = append(lines, lipgloss.NewStyle().Bold(true).Render("Requirements:"))
		lines = append(lines, fmt.Sprintf("  %s Checking requirements…", mt.spinner.View()))
		lines = append(lines, "")
	case reqFailed:
		lines = append(lines, lipgloss.NewStyle().Bold(true).Render("Requirements:"))
		errMsg := "unknown error"
		if reqErr != nil {
			errMsg = reqErr.Error()
		}
		errStyle := failStyle.Width(width - 6)
		lines = append(lines, fmt.Sprintf("  %s %s", failStyle.Render("✗"), errStyle.Render(errMsg)))
		lines = append(lines, "")
	case reqPassed:
		lines = append(lines, lipgloss.NewStyle().Bold(true).Render("Requirements:"))
		lines = append(lines, fmt.Sprintf("  %s All requirements met", StatusIcon(component.StatusDone, "")))
		lines = append(lines, "")
	}

	if len(mt.steps) > 0 {
		lines = append(lines, lipgloss.NewStyle().Bold(true).Render("Recipe Steps:"))
		lines = append(lines, "")
	}

	spinnerView := mt.spinner.View()

	stepNames := make(map[string]bool)
	for _, s := range mt.steps {
		stepNames[s.Name] = true
	}

	children := make(map[string][]string)
	for _, step := range mt.steps {
		if len(step.DependsOn) > 0 && stepNames[step.DependsOn[0]] {
			parent := step.DependsOn[0]
			children[parent] = append(children[parent], step.Name)
		}
	}

	rendered := make(map[string]bool)
	for _, step := range mt.steps {
		if !rendered[step.Name] {
			lines = append(lines, mt.renderStepTree(step.Name, 0, spinnerView, children, rendered)...)
		}
	}

	mt.viewport.SetContent(strings.Join(lines, "\n"))
}

func (mt *MainTab) ViewWithRequirements() string {
	return mt.viewport.View()
}

func (mt *MainTab) renderStepTree(name string, depth int, spinnerView string, children map[string][]string, rendered map[string]bool) []string {
	rendered[name] = true
	step := mt.GetStep(name)
	if step == nil {
		return nil
	}

	prefix := "  "
	for range depth {
		prefix += dimStyle.Render("│") + " "
	}

	var lines []string
	icon := StatusIcon(step.Status, spinnerView)
	line := fmt.Sprintf("%s%s %s", prefix, icon, step.Name)
	if step.Status == component.StatusFailed && step.Err != nil {
		line += failStyle.Render(fmt.Sprintf(" — %v", step.Err))
	}
	lines = append(lines, line)

	if len(step.Processes) > 1 || (len(step.Processes) == 1 && step.Processes[0].Name != step.Name) {
		procPrefix := prefix + dimStyle.Render("│") + " "
		for _, proc := range step.Processes {
			procIcon := StatusIcon(proc.Status, spinnerView)
			lines = append(lines, fmt.Sprintf("%s%s %s", procPrefix, procIcon, proc.Name))
		}
	}

	for _, childName := range children[name] {
		if !rendered[childName] {
			lines = append(lines, mt.renderStepTree(childName, depth+1, spinnerView, children, rendered)...)
		}
	}

	return lines
}
