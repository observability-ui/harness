package ui

import (
	"strings"

	"github.com/charmbracelet/bubbles/viewport"
	"obsui/internal/process"
)

type ProcessTab struct {
	Name     string
	proc     *process.Process
	viewport viewport.Model
}

func NewProcessTab(name string, proc *process.Process, width, height int) ProcessTab {
	vp := viewport.New(width, height)
	vp.SetContent("")
	return ProcessTab{Name: name, proc: proc, viewport: vp}
}

func (pt *ProcessTab) Sync() {
	if pt.proc != nil {
		content := strings.Join(pt.proc.Output.Lines(), "\n")
		pt.viewport.SetContent(content)
		pt.viewport.GotoBottom()
	}
}

func (pt *ProcessTab) SetSize(width, height int) {
	pt.viewport.Width = width
	pt.viewport.Height = height
}

func (pt ProcessTab) View() string {
	return pt.viewport.View()
}
