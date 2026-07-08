package ui

import (
	"context"
	"time"

	"obs/internal/component"
	"obs/internal/process"

	"github.com/charmbracelet/bubbles/help"
	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/spinner"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type stepUpdateMsg component.StepUpdate
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

type RequirementsCheckingMsg struct{}
type RequirementsPassedMsg struct{ Steps []*component.Step }
type RequirementsFailedMsg struct{ Err error }

type reqState int

const (
	reqChecking reqState = iota
	reqPassed
	reqFailed
)

type Model struct {
	tabBar              TabBar
	mainTab             MainTab
	processTabs         []ProcessTab
	manager             *process.Manager
	help                help.Model
	width               int
	height              int
	quitting            bool
	shutdown            bool
	updates             <-chan component.StepUpdate
	statusMsg           string
	reqStatus           reqState
	reqErr              error
	retryCh             chan<- struct{}
	cachedContentHeight int
}

func NewModel(mgr *process.Manager, updates <-chan component.StepUpdate, retryCh chan<- struct{}) Model {
	tabs := []string{"main"}
	mt := NewMainTab()
	h := help.New()
	h.ShortSeparator = " · "

	return Model{
		tabBar:    NewTabBar(tabs),
		mainTab:   mt,
		manager:   mgr,
		help:      h,
		updates:   updates,
		reqStatus: reqChecking,
		retryCh:   retryCh,
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
		m.recomputeLayout()

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
			if m.tabBar.Active() == 0 {
				m.mainTab.viewport.LineUp(1)
			} else if tab := m.activeProcessTab(); tab != nil {
				tab.viewport.LineUp(1)
			}
		case key.Matches(msg, keys.Down):
			if m.tabBar.Active() == 0 {
				m.mainTab.viewport.LineDown(1)
			} else if tab := m.activeProcessTab(); tab != nil {
				tab.viewport.LineDown(1)
			}
		case key.Matches(msg, keys.PageUp):
			if m.tabBar.Active() == 0 {
				m.mainTab.viewport.ViewUp()
			} else if tab := m.activeProcessTab(); tab != nil {
				tab.viewport.ViewUp()
			}
		case key.Matches(msg, keys.PageDown):
			if m.tabBar.Active() == 0 {
				m.mainTab.viewport.ViewDown()
			} else if tab := m.activeProcessTab(); tab != nil {
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
			stepName := tab.StepName
			if stepName == "" {
				continue
			}
			step := m.mainTab.GetStep(stepName)
			if step == nil || step.Status == component.StatusDone || step.Status == component.StatusStopped || step.Status == component.StatusFailed {
				continue
			}
			if !proc.Running() {
				procStatus := MapProcessStatus(proc.Status)
				var stepErr error
				if procStatus == component.StatusFailed {
					stepErr = proc.Err
				}
				m.mainTab.UpdateStep(stepName, procStatus, stepErr)
				m.mainTab.UpdateProcess(stepName, tab.Name, procStatus)
			} else if step.Status != component.StatusReady {
				m.mainTab.UpdateProcess(stepName, tab.Name, component.StatusStarted)
			}
		}
		cmds = append(cmds, tickOutputRefresh())

	case AddProcessTabMsg:
		found := false
		for i := range m.processTabs {
			if m.processTabs[i].Name == msg.Name {
				m.processTabs[i].SetProcess(msg.Proc)
				if msg.StepName != "" {
					m.processTabs[i].StepName = msg.StepName
				}
				found = true
				break
			}
		}
		if !found {
			m.addProcessTab(msg.Name, msg.Proc)
		}

	case processRestartedMsg:
		if msg.err != nil {
			m.statusMsg = "Restart failed: " + msg.err.Error()
		} else {
			for i := range m.processTabs {
				if m.processTabs[i].Name == msg.name {
					m.processTabs[i].SetProcess(msg.proc)
					m.mainTab.UpdateStep(m.processTabs[i].StepName, component.StatusStarted, nil)
					break
				}
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
			procNames := make([]string, 0, len(step.Processes))
			for _, spec := range step.Processes {
				procNames = append(procNames, spec.Name)
				pt := NewProcessTab(spec.Name, nil, m.width, m.contentHeight())
				pt.StepName = step.Name
				pt.DependsOn = step.DependsOn
				m.tabBar.Add(spec.Name)
				m.processTabs = append(m.processTabs, pt)
			}
			m.mainTab.AddStepWithProcesses(step.Name, procNames)
		}
		m.recomputeLayout()

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

	icons := m.tabIcons()
	tabBarView := m.tabBar.ViewWithIcons(m.width, icons)

	m.help.Width = m.width
	helpBar := m.help.View(keyMapForTab(m.tabBar.Active()))

	ch := m.contentHeight()

	var content string
	if m.tabBar.Active() == 0 {
		content = m.mainTab.ViewWithRequirements(m.width, ch, m.reqStatus, m.reqErr)
	} else {
		idx := m.tabBar.Active() - 1
		if idx < len(m.processTabs) {
			content = m.processTabs[idx].View()
		}
	}

	bottom := helpBar
	if m.statusMsg != "" {
		bottom = dimStyle.Width(m.width).Render(m.statusMsg) + "\n" + bottom
	}

	return tabBarView + "\n" + content + "\n" + bottom
}

func (m Model) contentHeight() int {
	return m.cachedContentHeight
}

func (m *Model) recomputeLayout() {
	tabBarHeight := lipgloss.Height(m.tabBar.View(m.width))
	bottomHeight := lipgloss.Height(m.help.View(keyMapForTab(0)))
	if m.statusMsg != "" {
		bottomHeight++
	}
	h := m.height - tabBarHeight - bottomHeight - 1
	if h < 1 {
		h = 1
	}
	m.cachedContentHeight = h
	m.mainTab.SetSize(m.width, h)
	for i := range m.processTabs {
		m.processTabs[i].SetSize(m.width, h)
	}
}

func (m *Model) addProcessTab(name string, proc *process.Process) {
	m.tabBar.Add(name)
	ch := m.contentHeight()
	pt := NewProcessTab(name, proc, m.width, ch)
	m.processTabs = append(m.processTabs, pt)
	m.recomputeLayout()
}

func (m Model) tabIcons() []string {
	icons := make([]string, m.tabBar.Count())
	spinnerView := m.mainTab.spinner.View()
	for i := range m.processTabs {
		tabIdx := i + 1
		if tabIdx >= len(icons) {
			break
		}
		stepName := m.processTabs[i].StepName
		if stepName == "" {
			continue
		}
		step := m.mainTab.GetStep(stepName)
		if step == nil {
			continue
		}
		icons[tabIdx] = StatusIcon(step.Status, spinnerView)
	}
	return icons
}

func (m *Model) activeProcessTab() *ProcessTab {
	if m.tabBar.Active() > 0 {
		idx := m.tabBar.Active() - 1
		if idx < len(m.processTabs) {
			return &m.processTabs[idx]
		}
	}
	return nil
}

func waitForUpdate(ch <-chan component.StepUpdate) tea.Cmd {
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
