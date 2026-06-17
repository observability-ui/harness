package ui

import (
	"strings"

	"github.com/charmbracelet/lipgloss"
)

var (
	activeTabStyle = lipgloss.NewStyle().
		Bold(true).
		Foreground(lipgloss.Color("15")).
		Background(lipgloss.Color("4")).
		Padding(0, 2)

	inactiveTabStyle = lipgloss.NewStyle().
		Foreground(lipgloss.Color("7")).
		Padding(0, 2)

	tabBarStyle = lipgloss.NewStyle().
		BorderBottom(true).
		BorderStyle(lipgloss.NormalBorder()).
		BorderForeground(lipgloss.Color("8"))
)

type TabBar struct {
	tabs   []string
	active int
}

func NewTabBar(tabs []string) TabBar {
	return TabBar{tabs: tabs, active: 0}
}

func (tb *TabBar) Add(name string) {
	tb.tabs = append(tb.tabs, name)
}

func (tb *TabBar) Next() {
	if len(tb.tabs) > 0 {
		tb.active = (tb.active + 1) % len(tb.tabs)
	}
}

func (tb *TabBar) Prev() {
	if len(tb.tabs) > 0 {
		tb.active = (tb.active - 1 + len(tb.tabs)) % len(tb.tabs)
	}
}

func (tb *TabBar) Active() int    { return tb.active }
func (tb *TabBar) Count() int     { return len(tb.tabs) }
func (tb *TabBar) ActiveName() string {
	if tb.active < len(tb.tabs) {
		return tb.tabs[tb.active]
	}
	return ""
}

func (tb TabBar) View(width int) string {
	var rendered []string
	for i, name := range tb.tabs {
		style := inactiveTabStyle
		if i == tb.active {
			style = activeTabStyle
		}
		rendered = append(rendered, style.Render(name))
	}
	row := strings.Join(rendered, " ")
	return tabBarStyle.Width(width).Render(row)
}
