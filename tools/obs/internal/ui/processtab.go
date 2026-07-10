package ui

import (
	"fmt"
	"strings"

	"obs/internal/process"

	"github.com/charmbracelet/bubbles/viewport"
	"github.com/charmbracelet/x/ansi"
)

type processTab struct {
	Name              string
	StepName          string
	DependsOn         []string
	proc              *process.Process
	viewport          viewport.Model
	lastRenderedCount int
	wrapMode          bool
}

func newProcessTab(name string, proc *process.Process, width, height int) processTab {
	vp := viewport.New(width, height)
	vp.SetContent("")
	return processTab{Name: name, proc: proc, viewport: vp}
}

func (pt *processTab) Sync() {
	if pt.proc == nil {
		return
	}
	count := pt.proc.Output.Len()
	if count == pt.lastRenderedCount {
		return
	}
	content := strings.Join(pt.proc.Output.Lines(), "\n")
	if pt.wrapMode && pt.viewport.Width > 0 {
		content = ansi.Hardwrap(content, pt.viewport.Width, false)
	}
	pt.viewport.SetContent(content)
	pt.viewport.GotoBottom()
	pt.lastRenderedCount = count
}

func (pt *processTab) ToggleWrap() {
	pt.wrapMode = !pt.wrapMode
	pt.lastRenderedCount = 0
}

func (pt *processTab) WrapMode() bool {
	return pt.wrapMode
}

func (pt *processTab) SetProcess(proc *process.Process) {
	pt.proc = proc
	pt.lastRenderedCount = 0
	pt.viewport.SetContent("")
}

func (pt *processTab) Process() *process.Process {
	return pt.proc
}

func (pt *processTab) SetSize(width, height int) {
	pt.viewport.Width = width
	pt.viewport.Height = height
}

func (pt *processTab) View() string {
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
