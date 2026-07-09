package ui

import (
	"context"
	"errors"
	"fmt"
	"time"

	"obs/internal/component"
	"obs/internal/mixer"
	"obs/internal/process"

	"github.com/charmbracelet/bubbles/help"
	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/spinner"
	"github.com/charmbracelet/bubbles/textinput"
	"github.com/charmbracelet/bubbles/viewport"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type stepUpdateMsg struct {
	component.StepUpdate
	gen        int
	restartGen int
}
type processOutputMsg struct{}
type shutdownCompleteMsg struct{ gen int }

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
	ctx                 context.Context
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
	inputMode           bool
	inputField          textinput.Model
	inputFlag           string
	inputValues    map[string]string
	generation     int
	stepRestartGen map[string]int
}

func NewModel(ctx context.Context, mgr *process.Manager, updates <-chan component.StepUpdate, retryCh chan<- struct{}, inputValues map[string]string) Model {
	tabs := []string{"main"}
	mt := NewMainTab()
	h := help.New()
	h.ShortSeparator = " · "

	ti := textinput.New()
	ti.CharLimit = 256

	if inputValues == nil {
		inputValues = make(map[string]string)
	}

	return Model{
		tabBar:         NewTabBar(tabs),
		mainTab:        mt,
		ctx:            ctx,
		manager:        mgr,
		help:           h,
		updates:        updates,
		reqStatus:      reqChecking,
		retryCh:        retryCh,
		inputField:     ti,
		inputValues:    inputValues,
		stepRestartGen: make(map[string]int),
	}
}

func (m Model) Init() tea.Cmd {
	return tea.Batch(
		m.mainTab.spinner.Tick,
		waitForUpdate(m.updates, m.generation),
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
		m.mainTab.RefreshContent(m.width, m.reqStatus, m.reqErr)

	case tea.KeyMsg:
		if m.inputMode {
			return m.handleInputKey(msg)
		}

		switch {
		case key.Matches(msg, keys.Quit):
			if m.shutdown {
				return m, tea.Quit
			}
			if !m.quitting {
				m.quitting = true
				m.statusMsg = "Stopping processes… press q again to exit"
				return m, stopProcesses(m.manager, m.generation)
			}
			return m, tea.Quit
		case key.Matches(msg, keys.NextTab):
			m.tabBar.Next()
		case key.Matches(msg, keys.PrevTab):
			m.tabBar.Prev()
		case key.Matches(msg, keys.Restart):
			if m.tabBar.Active() == 0 && m.reqStatus != reqChecking {
				m.statusMsg = ""
				cmds = append(cmds, m.resetForRetry())
				select {
				case m.retryCh <- struct{}{}:
				default:
				}
			} else if tab := m.activeProcessTab(); tab != nil {
				name := tab.Name
				mgr := m.manager
				ctx := m.ctx
				cmds = append(cmds, func() tea.Msg {
					proc, err := mgr.RestartProcess(ctx, name)
					return processRestartedMsg{name: name, proc: proc, err: err}
				})
			}
		case key.Matches(msg, keys.Up):
			if vp := m.activeViewport(); vp != nil {
				vp.ScrollUp(1)
			}
		case key.Matches(msg, keys.Down):
			if vp := m.activeViewport(); vp != nil {
				vp.ScrollDown(1)
			}
		case key.Matches(msg, keys.PageUp):
			if vp := m.activeViewport(); vp != nil {
				vp.PageUp()
			}
		case key.Matches(msg, keys.PageDown):
			if vp := m.activeViewport(); vp != nil {
				vp.PageDown()
			}
		}

	case spinner.TickMsg:
		var cmd tea.Cmd
		m.mainTab.spinner, cmd = m.mainTab.spinner.Update(msg)
		m.mainTab.RefreshContent(m.width, m.reqStatus, m.reqErr)
		cmds = append(cmds, cmd)

	case stepUpdateMsg:
		if msg.gen == m.generation && msg.restartGen >= m.stepRestartGen[msg.StepName] {
			m.mainTab.UpdateStep(msg.StepName, msg.Status, msg.Err)
		}
		m.syncProcessDisplayStatus()
		m.mainTab.RefreshContent(m.width, m.reqStatus, m.reqErr)
		if msg.restartGen == 0 {
			cmds = append(cmds, waitForUpdate(m.updates, m.generation))
		} else if msg.Status == component.StatusReady {
			if tab := m.findProcessTab(msg.StepName); tab != nil {
				if proc := tab.Process(); proc != nil {
					gen := m.generation
					restartGen := msg.restartGen
					stepName := msg.StepName
					cmds = append(cmds, func() tea.Msg {
						status := process.WaitForExit(proc)
						return stepUpdateMsg{
							StepUpdate: component.StepUpdate{StepName: stepName, Status: status},
							gen:        gen,
							restartGen: restartGen,
						}
					})
				}
			}
		}

	case processOutputMsg:
		if tab := m.activeProcessTab(); tab != nil {
			tab.Sync()
		}
		m.syncProcessDisplayStatus()
		m.mainTab.RefreshContent(m.width, m.reqStatus, m.reqErr)
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
			if msg.StepName != "" {
				m.processTabs[len(m.processTabs)-1].StepName = msg.StepName
			}
		}

	case processRestartedMsg:
		if msg.err != nil {
			m.statusMsg = "Restart failed: " + msg.err.Error()
		} else {
			for i := range m.processTabs {
				if m.processTabs[i].Name == msg.name {
					stepName := m.processTabs[i].StepName
					m.processTabs[i].SetProcess(msg.proc)
					m.stepRestartGen[stepName]++
					m.mainTab.UpdateStep(stepName, component.StatusStarted, nil)
					m.mainTab.RefreshContent(m.width, m.reqStatus, m.reqErr)
					gen := m.generation
					restartGen := m.stepRestartGen[stepName]
					proc := msg.proc
					ports := proc.Spec.Ports
					cmds = append(cmds, func() tea.Msg {
						status := process.WaitForReady(proc, ports)
						return stepUpdateMsg{
							StepUpdate: component.StepUpdate{StepName: stepName, Status: status},
							gen:        gen,
							restartGen: restartGen,
						}
					})
					break
				}
			}
		}

	case RequirementsCheckingMsg:
		m.reqStatus = reqChecking
		m.reqErr = nil
		if len(m.mainTab.steps) > 0 || len(m.processTabs) > 0 {
			m.mainTab = NewMainTab()
			m.processTabs = nil
			m.tabBar = NewTabBar([]string{"main"})
		}
		m.recomputeLayout()
		m.mainTab.RefreshContent(m.width, m.reqStatus, m.reqErr)
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
			m.mainTab.AddStepWithProcesses(step.Name, procNames, step.DependsOn)
		}
		m.recomputeLayout()
		m.mainTab.RefreshContent(m.width, m.reqStatus, m.reqErr)

	case RequirementsFailedMsg:
		if mfe, ok := errors.AsType[*mixer.MissingFlagError](msg.Err); ok {
			m.inputMode = true
			m.inputFlag = mfe.Flag
			m.inputField.Placeholder = mfe.Usage
			m.inputField.SetValue("")
			m.inputField.Width = m.width - len(m.inputFlag) - 12
			m.inputField.Focus()
			m.reqStatus = reqFailed
			m.reqErr = msg.Err
			return m, m.inputField.Cursor.BlinkCmd()
		}
		m.reqStatus = reqFailed
		m.reqErr = msg.Err
		m.mainTab.RefreshContent(m.width, m.reqStatus, m.reqErr)

	case shutdownCompleteMsg:
		if msg.gen == m.generation {
			m.shutdown = true
		}

	default:
		if m.inputMode {
			var cmd tea.Cmd
			m.inputField, cmd = m.inputField.Update(msg)
			cmds = append(cmds, cmd)
		}
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

	var content string
	if m.tabBar.Active() == 0 {
		content = m.mainTab.ViewWithRequirements()
		if m.inputMode {
			prompt := fmt.Sprintf("\n  Enter --%s: %s", m.inputFlag, m.inputField.View())
			content += prompt
		}
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
	m.help.Width = m.width
	bottomHeight := lipgloss.Height(m.help.View(keyMapForTab(0)))
	if ph := lipgloss.Height(m.help.View(keyMapForTab(1))); ph > bottomHeight {
		bottomHeight = ph
	}
	if m.statusMsg != "" {
		bottomHeight++
	}
	h := max(m.height-tabBarHeight-bottomHeight-2, 1)
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

func (m *Model) findProcessTab(stepName string) *ProcessTab {
	for i := range m.processTabs {
		if m.processTabs[i].StepName == stepName {
			return &m.processTabs[i]
		}
	}
	return nil
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

func (m *Model) resetForRetry() tea.Cmd {
	m.reqStatus = reqChecking
	m.reqErr = nil
	m.quitting = false
	m.shutdown = false
	m.generation++
	m.stepRestartGen = make(map[string]int)
	m.mainTab = NewMainTab()
	m.processTabs = nil
	m.tabBar = NewTabBar([]string{"main"})
	m.recomputeLayout()
	return m.mainTab.spinner.Tick
}

func (m *Model) activeViewport() *viewport.Model {
	if m.tabBar.Active() == 0 {
		return &m.mainTab.viewport
	}
	if tab := m.activeProcessTab(); tab != nil {
		return &tab.viewport
	}
	return nil
}

func (m *Model) handleInputKey(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	if key.Matches(msg, keys.Quit) {
		m.inputMode = false
		m.inputField.Blur()
		m.quitting = true
		m.statusMsg = "Stopping processes… press q again to exit"
		return *m, stopProcesses(m.manager, m.generation)
	}
	var cmds []tea.Cmd
	switch msg.Type {
	case tea.KeyEnter:
		val := m.inputField.Value()
		if val != "" {
			m.inputValues[m.inputFlag] = val
			m.inputMode = false
			m.inputField.Blur()
			cmd := m.resetForRetry()
			m.statusMsg = fmt.Sprintf("Set --%s=%s", m.inputFlag, val)
			cmds = append(cmds, cmd)
			select {
			case m.retryCh <- struct{}{}:
			default:
			}
		}
	case tea.KeyEsc:
		m.inputMode = false
		m.inputField.Blur()
	default:
		var cmd tea.Cmd
		m.inputField, cmd = m.inputField.Update(msg)
		cmds = append(cmds, cmd)
	}
	return *m, tea.Batch(cmds...)
}

func (m *Model) syncProcessDisplayStatus() {
	for i := range m.processTabs {
		tab := &m.processTabs[i]
		proc := tab.Process()
		if proc == nil || tab.StepName == "" {
			continue
		}
		step := m.mainTab.GetStep(tab.StepName)
		if step == nil {
			continue
		}
		for j := range step.Processes {
			if step.Processes[j].Name == tab.Name {
				if !proc.Running() {
					step.Processes[j].Status = process.MapExitStatus(proc)
				} else {
					step.Processes[j].Status = component.StatusStarted
				}
				break
			}
		}
	}
}

func waitForUpdate(ch <-chan component.StepUpdate, gen int) tea.Cmd {
	return func() tea.Msg {
		update, ok := <-ch
		if !ok {
			return nil
		}
		return stepUpdateMsg{StepUpdate: update, gen: gen}
	}
}

func tickOutputRefresh() tea.Cmd {
	return tea.Tick(time.Millisecond*100, func(t time.Time) tea.Msg {
		return processOutputMsg{}
	})
}

func stopProcesses(mgr *process.Manager, gen int) tea.Cmd {
	return func() tea.Msg {
		mgr.StopAll()
		return shutdownCompleteMsg{gen: gen}
	}
}

