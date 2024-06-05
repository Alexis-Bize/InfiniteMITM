// Copyright 2024 Alexis Bize
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//		https://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package MITMApplicationUIServiceTestTable

import (
	"fmt"
	"os"
	"strings"

	"github.com/charmbracelet/bubbles/table"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type tab string

const (
	requestTab  tab = "Request"
	responseTab tab = "Response"
)

type model struct {
	requestsTable table.Model
	tabs          []tab
	activeTab     tab
	action        string
	details       map[tab]string
	selectedIndex int
	windowWidth   int
	windowHeight  int
}

func Create() {
	tabs := []tab{requestTab, responseTab}

	details := map[tab]string{
		requestTab:  "Request details will appear here.",
		responseTab: "Response details will appear here.",
	}

	m := model{
		tabs:          tabs,
		activeTab:     requestTab,
		action:        "Press q to quit, tab to switch tabs, enter to view details",
		details:       details,
		selectedIndex: -1,
	}

	p := tea.NewProgram(m)
	if _, err := p.Run(); err != nil {
		fmt.Printf("Error: %v", err)
		os.Exit(1)
	}
}

func (m model) Init() tea.Cmd {
	return nil
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "q", "ctrl+c":
			return m, tea.Quit
		case "tab":
			m.switchTab()
		case "enter":
			m.selectRequest()
		}
	case tea.WindowSizeMsg:
		m.windowWidth = msg.Width
		m.windowHeight = msg.Height
		m.updateTable()
	}

	var cmd tea.Cmd
	m.requestsTable, cmd = m.requestsTable.Update(msg)
	return m, cmd
}

func (m *model) updateTable() {
	if m.windowWidth == 0 {
		return
	}

	columnWidth := m.windowWidth / 3

	columns := []table.Column{
		{Title: "Method", Width: columnWidth / 3},
		{Title: "Status", Width: columnWidth / 6},
		{Title: "URL", Width: columnWidth},
	}

	rows := []table.Row{
		{"GET", "200", "https://example.com"},
		{"POST", "200", "https://example.com/login"},
	}

	m.requestsTable = table.New(
		table.WithColumns(columns),
		table.WithRows(rows),
		table.WithFocused(true),
	)
}

func (m *model) switchTab() {
	if m.activeTab == requestTab {
		m.activeTab = responseTab
	} else {
		m.activeTab = requestTab
	}
}

func (m *model) selectRequest() {
	row := m.requestsTable.SelectedRow()
	if len(row) > 0 {
		m.selectedIndex = m.requestsTable.Cursor()
		m.details[requestTab] = fmt.Sprintf("Request details for %s %s", row[0], row[1])
		m.details[responseTab] = fmt.Sprintf("Response details for %s %s", row[0], row[1])
	}
}

func (m model) View() string {
	// Define styles
	tableStyle := lipgloss.NewStyle().Border(lipgloss.NormalBorder()).Padding(1, 2)
	tabStyle := lipgloss.NewStyle().MarginRight(1).Padding(0, 1)
	activeTabStyle := lipgloss.NewStyle().MarginRight(1).Padding(0, 1).Bold(true).Underline(true)
	detailsStyle := lipgloss.NewStyle().Border(lipgloss.NormalBorder()).Padding(1, 2)
	actionStyle := lipgloss.NewStyle().Border(lipgloss.NormalBorder()).Padding(0, 1).MarginTop(0)

	// Render tabs
	var renderedTabs []string
	for _, t := range m.tabs {
		if t == m.activeTab {
			renderedTabs = append(renderedTabs, activeTabStyle.Render(string(t)))
		} else {
			renderedTabs = append(renderedTabs, tabStyle.Render(string(t)))
		}
	}
	tabsView := strings.Join(renderedTabs, " ")

	// Combine views
	requestsView := tableStyle.Render(m.requestsTable.View())
	detailsView := detailsStyle.Render(tabsView + "\n\n" + m.details[m.activeTab])
	actionView := actionStyle.Render(m.action)

	return lipgloss.JoinHorizontal(lipgloss.Top, requestsView, detailsView) + "\n" + actionView
}
