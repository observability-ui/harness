package ui

import (
	"strings"

	"github.com/charmbracelet/lipgloss"
)

var (
	tabPadding = lipgloss.NewStyle().Padding(0, 2)

	activeNameStyle = lipgloss.NewStyle().
		Bold(true).
		Underline(true)

	inactiveNameStyle = lipgloss.NewStyle().
		Foreground(lipgloss.Color("8"))

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
	return tb.ViewWithIcons(width, nil)
}

func (tb TabBar) ViewWithIcons(width int, icons []string) string {
	var rendered []string
	for i, name := range tb.tabs {
		nameStyle := inactiveNameStyle
		if i == tb.active {
			nameStyle = activeNameStyle
		}
		var display string
		if i < len(icons) && icons[i] != "" {
			display = icons[i] + " " + nameStyle.Render(name)
		} else {
			display = nameStyle.Render(name)
		}
		rendered = append(rendered, tabPadding.Render(display))
	}
	row := strings.Join(rendered, " ")
	return tabBarStyle.Width(width).Render(row)
}
