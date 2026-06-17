package ui

import (
	"time"

	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/spinner"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"obsui/internal/process"
	"obsui/internal/recipe"
	"obsui/internal/types"
)

type stepUpdateMsg types.StepUpdate
type processOutputMsg struct{}
type shutdownCompleteMsg struct{}

type AddProcessTabMsg struct {
	Name string
	Proc *process.Process
}

type Model struct {
	tabBar      TabBar
	mainTab     MainTab
	processTabs []ProcessTab
	manager     *process.Manager
	width       int
	height      int
	quitting    bool
	shutdown    bool
	updates     <-chan types.StepUpdate
	statusBar   string
}

func NewModel(mgr *process.Manager, steps []*recipe.Step, updates <-chan types.StepUpdate) Model {
	tabs := []string{"main"}
	mt := NewMainTab()

	for _, step := range steps {
		mt.AddStep(step.Name)
	}

	return Model{
		tabBar:  NewTabBar(tabs),
		mainTab: mt,
		manager: mgr,
		updates: updates,
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
		contentHeight := m.height - 4 // tab bar + status bar
		for i := range m.processTabs {
			m.processTabs[i].SetSize(m.width, contentHeight)
		}

	case tea.KeyMsg:
		switch {
		case keyMatches(msg, keys.Quit):
			if m.shutdown {
				return m, tea.Quit
			}
			if !m.quitting {
				m.quitting = true
				m.statusBar = "Stopping processes… press q again to exit"
				return m, stopProcesses(m.manager)
			}
			return m, tea.Quit
		case keyMatches(msg, keys.NextTab):
			m.tabBar.Next()
		case keyMatches(msg, keys.PrevTab):
			m.tabBar.Prev()
		case keyMatches(msg, keys.Restart):
			// Only restart if on a process tab
			if m.tabBar.Active() > 0 {
				idx := m.tabBar.Active() - 1
				if idx < len(m.processTabs) {
					m.statusBar = "Restarting " + m.processTabs[idx].Name + "…"
				}
			}
		case keyMatches(msg, keys.Up):
			if m.tabBar.Active() > 0 {
				idx := m.tabBar.Active() - 1
				if idx < len(m.processTabs) {
					m.processTabs[idx].viewport.LineUp(1)
				}
			}
		case keyMatches(msg, keys.Down):
			if m.tabBar.Active() > 0 {
				idx := m.tabBar.Active() - 1
				if idx < len(m.processTabs) {
					m.processTabs[idx].viewport.LineDown(1)
				}
			}
		case keyMatches(msg, keys.PageUp):
			if m.tabBar.Active() > 0 {
				idx := m.tabBar.Active() - 1
				if idx < len(m.processTabs) {
					m.processTabs[idx].viewport.ViewUp()
				}
			}
		case keyMatches(msg, keys.PageDown):
			if m.tabBar.Active() > 0 {
				idx := m.tabBar.Active() - 1
				if idx < len(m.processTabs) {
					m.processTabs[idx].viewport.ViewDown()
				}
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
		for i := range m.processTabs {
			m.processTabs[i].Sync()
		}
		cmds = append(cmds, tickOutputRefresh())

	case AddProcessTabMsg:
		m.addProcessTab(msg.Name, msg.Proc)

	case shutdownCompleteMsg:
		m.shutdown = true
		m.statusBar = "All processes stopped. Press q to exit."
	}

	return m, tea.Batch(cmds...)
}

func (m Model) View() string {
	if m.width == 0 {
		return "initializing…"
	}

	tabBarView := m.tabBar.View(m.width)
	contentHeight := m.height - lipgloss.Height(tabBarView) - 2

	var content string
	if m.tabBar.Active() == 0 {
		content = m.mainTab.View(m.width, contentHeight)
	} else {
		idx := m.tabBar.Active() - 1
		if idx < len(m.processTabs) {
			content = m.processTabs[idx].View()
		}
	}

	statusStyle := lipgloss.NewStyle().
		Foreground(lipgloss.Color("8")).
		Width(m.width)
	status := statusStyle.Render(m.statusBar)

	return tabBarView + "\n" + content + "\n" + status
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

func keyMatches(msg tea.KeyMsg, binding key.Binding) bool {
	return key.Matches(msg, binding)
}

func waitForUpdate(ch <-chan types.StepUpdate) tea.Cmd {
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
