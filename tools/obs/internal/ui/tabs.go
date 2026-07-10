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

	tabBarStyle = lipgloss.NewStyle().
			BorderBottom(true).
			BorderStyle(lipgloss.NormalBorder()).
			BorderForeground(lipgloss.Color("8"))
)

type tabBar struct {
	tabs   []string
	active int
}

func newTabBar(tabs []string) tabBar {
	return tabBar{tabs: tabs, active: 0}
}

func (tb *tabBar) Add(name string) {
	tb.tabs = append(tb.tabs, name)
}

func (tb *tabBar) Next() {
	if len(tb.tabs) > 0 {
		tb.active = (tb.active + 1) % len(tb.tabs)
	}
}

func (tb *tabBar) Prev() {
	if len(tb.tabs) > 0 {
		tb.active = (tb.active - 1 + len(tb.tabs)) % len(tb.tabs)
	}
}

func (tb *tabBar) Active() int { return tb.active }
func (tb *tabBar) Count() int  { return len(tb.tabs) }

func (tb *tabBar) View(width int) string {
	return tb.ViewWithIcons(width, nil)
}

func (tb *tabBar) ViewWithIcons(width int, icons []string) string {
	var rendered []string
	for i, name := range tb.tabs {
		nameStyle := dimStyle
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
