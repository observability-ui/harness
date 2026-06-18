package ui

import (
	"context"
	"time"

	"github.com/charmbracelet/bubbles/help"
	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/spinner"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"obs/internal/process"
	"obs/internal/recipe"
)

type stepUpdateMsg recipe.StepUpdate
type processOutputMsg struct{}
type shutdownCompleteMsg struct{}

type AddProcessTabMsg struct {
	StepName string
	Name     string
	Proc     *process.Process
}

type processRestartedMsg struct {
	name string
	proc *process.Process
	err  error
}

type readinessResultMsg struct {
	processName string
	ready       bool
}

type RequirementsCheckingMsg struct{}
type RequirementsPassedMsg struct{ Steps []*recipe.Step }
type RequirementsFailedMsg struct{ Err error }

type reqState int

const (
	reqChecking reqState = iota
	reqPassed
	reqFailed
)

type Model struct {
	tabBar      TabBar
	mainTab     MainTab
	processTabs []ProcessTab
	manager     *process.Manager
	help        help.Model
	width       int
	height      int
	quitting    bool
	shutdown    bool
	updates     <-chan recipe.StepUpdate
	statusMsg      string
	reqStatus      reqState
	reqErr         error
	retryCh        chan<- struct{}
	processToStep  map[string]string
}

func NewModel(mgr *process.Manager, steps []*recipe.Step, updates <-chan recipe.StepUpdate, retryCh chan<- struct{}) Model {
	tabs := []string{"main"}
	mt := NewMainTab()

	for _, step := range steps {
		mt.AddStep(step.Name)
	}

	h := help.New()
	h.ShortSeparator = " · "

	return Model{
		tabBar:        NewTabBar(tabs),
		mainTab:       mt,
		manager:       mgr,
		help:          h,
		updates:       updates,
		reqStatus:     reqChecking,
		retryCh:       retryCh,
		processToStep: make(map[string]string),
	}
}

func (m Model) Init() tea.Cmd {
	return tea.Batch(
		m.mainTab.spinner.Tick,
		waitForUpdate(m.updates),
		tickOutputRefresh(),
	)
}

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmds []tea.Cmd

	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height
		contentHeight := m.height - 4
		for i := range m.processTabs {
			m.processTabs[i].SetSize(m.width, contentHeight)
		}

	case tea.KeyMsg:
		switch {
		case key.Matches(msg, keys.Quit):
			if m.shutdown {
				return m, tea.Quit
			}
			if !m.quitting {
				m.quitting = true
				m.statusMsg = "Stopping processes… press q again to exit"
				return m, stopProcesses(m.manager)
			}
			return m, tea.Quit
		case key.Matches(msg, keys.NextTab):
			m.tabBar.Next()
		case key.Matches(msg, keys.PrevTab):
			m.tabBar.Prev()
		case key.Matches(msg, keys.Restart):
			if m.tabBar.Active() == 0 && m.reqStatus != reqChecking {
				m.reqStatus = reqChecking
				m.reqErr = nil
				m.statusMsg = ""
				m.mainTab = NewMainTab()
				m.processTabs = nil
				m.processToStep = make(map[string]string)
				m.tabBar = NewTabBar([]string{"main"})
				cmds = append(cmds, m.mainTab.spinner.Tick)
				select {
				case m.retryCh <- struct{}{}:
				default:
				}
			} else if tab := m.activeProcessTab(); tab != nil {
				name := tab.Name
				mgr := m.manager
				cmds = append(cmds, func() tea.Msg {
					proc, err := mgr.RestartProcess(context.Background(), name)
					return processRestartedMsg{name: name, proc: proc, err: err}
				})
			}
		case key.Matches(msg, keys.Up):
			if tab := m.activeProcessTab(); tab != nil {
				tab.viewport.LineUp(1)
			}
		case key.Matches(msg, keys.Down):
			if tab := m.activeProcessTab(); tab != nil {
				tab.viewport.LineDown(1)
			}
		case key.Matches(msg, keys.PageUp):
			if tab := m.activeProcessTab(); tab != nil {
				tab.viewport.ViewUp()
			}
		case key.Matches(msg, keys.PageDown):
			if tab := m.activeProcessTab(); tab != nil {
				tab.viewport.ViewDown()
			}
		}

	case spinner.TickMsg:
		var cmd tea.Cmd
		m.mainTab.spinner, cmd = m.mainTab.spinner.Update(msg)
		cmds = append(cmds, cmd)

	case stepUpdateMsg:
		m.mainTab.UpdateStep(msg.StepName, msg.Status, msg.Err)
		cmds = append(cmds, waitForUpdate(m.updates))

	case processOutputMsg:
		if tab := m.activeProcessTab(); tab != nil {
			tab.Sync()
		}
		// Monitor process lifecycle: exit detection + readiness probing
		for i := range m.processTabs {
			tab := &m.processTabs[i]
			proc := tab.Process()
			if proc == nil {
				continue
			}
			stepName, ok := m.processToStep[tab.Name]
			if !ok {
				continue
			}
			step := m.mainTab.GetStep(stepName)
			if step == nil || step.Status == recipe.StatusDone || step.Status == recipe.StatusStopped || step.Status == recipe.StatusFailed {
				continue
			}
			if !proc.Running() {
				if proc.Status == process.ProcessFailed {
					m.mainTab.UpdateStep(stepName, recipe.StatusFailed, proc.Err)
				} else if proc.Status == process.ProcessStopped {
					m.mainTab.UpdateStep(stepName, recipe.StatusStopped, nil)
				} else {
					m.mainTab.UpdateStep(stepName, recipe.StatusDone, nil)
				}
				continue
			}
			if step.Status == recipe.StatusStarted && len(proc.Spec.Ports) > 0 {
				pName := tab.Name
				ports := proc.Spec.Ports
				cmds = append(cmds, func() tea.Msg {
					return readinessResultMsg{processName: pName, ready: process.ProbePorts(ports)}
				})
			}
		}
		cmds = append(cmds, tickOutputRefresh())

	case readinessResultMsg:
		if msg.ready {
			if stepName, ok := m.processToStep[msg.processName]; ok {
				step := m.mainTab.GetStep(stepName)
				if step != nil && step.Status == recipe.StatusStarted {
					m.mainTab.UpdateStep(stepName, recipe.StatusReady, nil)
				}
			}
		}

	case AddProcessTabMsg:
		m.addProcessTab(msg.Name, msg.Proc)
		if msg.StepName != "" {
			m.processToStep[msg.Name] = msg.StepName
		}

	case processRestartedMsg:
		if msg.err != nil {
			m.statusMsg = "Restart failed: " + msg.err.Error()
		} else {
			for i := range m.processTabs {
				if m.processTabs[i].Name == msg.name {
					m.processTabs[i].SetProcess(msg.proc)
					break
				}
			}
			if stepName, ok := m.processToStep[msg.name]; ok {
				m.mainTab.UpdateStep(stepName, recipe.StatusStarted, nil)
			}
		}

	case RequirementsCheckingMsg:
		m.reqStatus = reqChecking
		m.reqErr = nil
		m.mainTab = NewMainTab()
		m.processTabs = nil
		m.tabBar = NewTabBar([]string{"main"})
		cmds = append(cmds, m.mainTab.spinner.Tick)

	case RequirementsPassedMsg:
		m.reqStatus = reqPassed
		for _, step := range msg.Steps {
			m.mainTab.AddStep(step.Name)
		}

	case RequirementsFailedMsg:
		m.reqStatus = reqFailed
		m.reqErr = msg.Err

	case shutdownCompleteMsg:
		m.shutdown = true
	}

	return m, tea.Batch(cmds...)
}

func (m Model) View() string {
	if m.width == 0 {
		return "initializing…"
	}

	tabBarView := m.tabBar.View(m.width)

	m.help.Width = m.width
	helpBar := m.help.View(keyMapForTab(m.tabBar.Active()))

	bottomHeight := lipgloss.Height(helpBar)
	if m.statusMsg != "" {
		bottomHeight++
	}
	contentHeight := m.height - lipgloss.Height(tabBarView) - bottomHeight - 1

	var content string
	if m.tabBar.Active() == 0 {
		content = m.mainTab.ViewWithRequirements(m.width, contentHeight, m.reqStatus, m.reqErr)
	} else {
		idx := m.tabBar.Active() - 1
		if idx < len(m.processTabs) {
			content = m.processTabs[idx].View()
		}
	}

	bottom := helpBar
	if m.statusMsg != "" {
		bottom = statusBarStyle.Width(m.width).Render(m.statusMsg) + "\n" + bottom
	}

	return tabBarView + "\n" + content + "\n" + bottom
}

func (m *Model) addProcessTab(name string, proc *process.Process) {
	contentHeight := m.height - 4
	if contentHeight < 10 {
		contentHeight = 24
	}
	pt := NewProcessTab(name, proc, m.width, contentHeight)
	m.processTabs = append(m.processTabs, pt)
	m.tabBar.Add(name)
}

var statusBarStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("8"))

func (m *Model) activeProcessTab() *ProcessTab {
	if m.tabBar.Active() > 0 {
		idx := m.tabBar.Active() - 1
		if idx < len(m.processTabs) {
			return &m.processTabs[idx]
		}
	}
	return nil
}

func waitForUpdate(ch <-chan recipe.StepUpdate) tea.Cmd {
	return func() tea.Msg {
		update, ok := <-ch
		if !ok {
			return nil
		}
		return stepUpdateMsg(update)
	}
}

func tickOutputRefresh() tea.Cmd {
	return tea.Tick(time.Millisecond*100, func(t time.Time) tea.Msg {
		return processOutputMsg{}
	})
}

func stopProcesses(mgr *process.Manager) tea.Cmd {
	return func() tea.Msg {
		mgr.StopAll()
		return shutdownCompleteMsg{}
	}
}
