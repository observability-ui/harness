package ui

import (
	"fmt"
	"strings"

	"obs/internal/process"

	"github.com/charmbracelet/bubbles/viewport"
)

type ProcessTab struct {
	Name              string
	StepName          string
	DependsOn         []string
	proc              *process.Process
	viewport          viewport.Model
	lastRenderedCount int
}

func NewProcessTab(name string, proc *process.Process, width, height int) ProcessTab {
	vp := viewport.New(width, height)
	vp.SetContent("")
	return ProcessTab{Name: name, proc: proc, viewport: vp}
}

func (pt *ProcessTab) Sync() {
	if pt.proc == nil {
		return
	}
	count := pt.proc.Output.Len()
	if count == pt.lastRenderedCount {
		return
	}
	content := strings.Join(pt.proc.Output.Lines(), "\n")
	pt.viewport.SetContent(content)
	pt.viewport.GotoBottom()
	pt.lastRenderedCount = count
}

func (pt *ProcessTab) SetProcess(proc *process.Process) {
	pt.proc = proc
	pt.lastRenderedCount = 0
	pt.viewport.SetContent("")
}

func (pt *ProcessTab) Process() *process.Process {
	return pt.proc
}

func (pt *ProcessTab) SetSize(width, height int) {
	pt.viewport.Width = width
	pt.viewport.Height = height
}

func (pt *ProcessTab) View() string {
	if pt.proc == nil && len(pt.DependsOn) > 0 {
		var lines []string
		lines = append(lines, "")
		lines = append(lines, dimStyle.Render("Waiting for:"))
		for _, dep := range pt.DependsOn {
			lines = append(lines, dimStyle.Render(fmt.Sprintf("  • %s", dep)))
		}
		return strings.Join(lines, "\n")
	}
	return pt.viewport.View()
}
