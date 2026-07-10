package ui

import (
	"fmt"
	"strings"

	"obs/internal/task"

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

var statusDefs = map[task.Status]statusDef{
	task.StatusPending: {Icon: "○"},
	task.StatusWaiting: {Icon: "◷", Style: lipgloss.NewStyle().Foreground(lipgloss.Color("3"))},
	task.StatusRunning: {UseSpinner: true},
	task.StatusStarted: {UseSpinner: true},
	task.StatusReady:   {Icon: "●", Style: lipgloss.NewStyle().Foreground(lipgloss.Color("2"))},
	task.StatusDone:    {Icon: "✓", Style: lipgloss.NewStyle().Foreground(lipgloss.Color("2"))},
	task.StatusStopped: {Icon: "■", Style: lipgloss.NewStyle().Foreground(lipgloss.Color("8"))},
	task.StatusFailed:  {Icon: "✗", Style: failStyle},
	task.StatusSkipped: {Icon: "⊘"},
}

func statusIcon(status task.Status, spinnerView string) string {
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

type processInfo struct {
	Name   string
	Status task.Status
}

type stepState struct {
	Name      string
	Status    task.Status
	Err       error
	Processes []processInfo
	DependsOn []string
}

type mainTab struct {
	steps    []stepState
	projects []task.ProjectInfo
	spinner  spinner.Model
	viewport viewport.Model
}

func newMainTab() mainTab {
	s := spinner.New()
	s.Spinner = spinner.Dot
	s.Style = lipgloss.NewStyle().Foreground(lipgloss.Color("6"))
	return mainTab{spinner: s}
}

func (mt *mainTab) SetSize(width, height int) {
	mt.viewport.Width = width
	mt.viewport.Height = height
}

func (mt *mainTab) SetProjects(projects []task.ProjectInfo) {
	mt.projects = projects
}

func (mt *mainTab) AddStepWithProcesses(name string, processNames []string, dependsOn []string) {
	var procs []processInfo
	for _, pn := range processNames {
		procs = append(procs, processInfo{Name: pn, Status: task.StatusPending})
	}
	mt.steps = append(mt.steps, stepState{Name: name, Status: task.StatusPending, Processes: procs, DependsOn: dependsOn})
}

func (mt *mainTab) updateProcess(stepName, procName string, status task.Status) {
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

func (mt *mainTab) GetStep(name string) *stepState {
	for i := range mt.steps {
		if mt.steps[i].Name == name {
			return &mt.steps[i]
		}
	}
	return nil
}

func (mt *mainTab) UpdateStep(name string, status task.Status, err error) {
	for i := range mt.steps {
		if mt.steps[i].Name == name {
			mt.steps[i].Status = status
			mt.steps[i].Err = err
			return
		}
	}
}

func (mt *mainTab) RefreshContent(width int, reqStatus reqState, reqErr error) {
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
		lines = append(lines, fmt.Sprintf("  %s All requirements met", statusIcon(task.StatusDone, "")))
		lines = append(lines, "")
	}

	if len(mt.projects) > 0 {
		lines = append(lines, lipgloss.NewStyle().Bold(true).Render("Projects:"))
		for _, p := range mt.projects {
			label := p.Branch
			if p.IsImage {
				label = "image:" + p.Branch
			}
			lines = append(lines, fmt.Sprintf("  %-25s %s", p.Name, dimStyle.Render(label)))
		}
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

func (mt *mainTab) ViewWithRequirements() string {
	return mt.viewport.View()
}

func (mt *mainTab) renderStepTree(name string, depth int, spinnerView string, children map[string][]string, rendered map[string]bool) []string {
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
	icon := statusIcon(step.Status, spinnerView)
	line := fmt.Sprintf("%s%s %s", prefix, icon, step.Name)
	if step.Status == task.StatusFailed && step.Err != nil {
		line += failStyle.Render(fmt.Sprintf(" — %v", step.Err))
	}
	lines = append(lines, line)

	if len(step.Processes) > 1 || (len(step.Processes) == 1 && step.Processes[0].Name != step.Name) {
		procPrefix := prefix + dimStyle.Render("│") + " "
		for _, proc := range step.Processes {
			procIcon := statusIcon(proc.Status, spinnerView)
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
