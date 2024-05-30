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

package InfiniteMITMApplicationUIServiceTable

import (
	"fmt"
	"strconv"

	"github.com/charmbracelet/bubbles/table"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

var program *tea.Program
var modelInstance *model

var baseStyle = lipgloss.NewStyle().
	BorderStyle(lipgloss.NormalBorder()).
	BorderForeground(lipgloss.Color("240"))

type model struct {
	table table.Model
}

type IOUpdate table.Row

func (m *model) Init() tea.Cmd {
	return nil
}

func (m *model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd
	switch msg := msg.(type) {
	case IOUpdate:
		m.table.SetRows(append(m.table.Rows(), table.Row(msg)))
	case tea.KeyMsg:
		switch msg.String() {
		case "esc":
			if m.table.Focused() {
				m.table.Blur()
			} else {
				m.table.Focus()
			}
		case "q", "ctrl+c":
			return m, tea.Quit
		case "enter":
			return m, tea.Batch(
				tea.Printf("Let's go to %s!", m.table.SelectedRow()[1]),
			)
		}
	}

	m.table, cmd = m.table.Update(msg)
	return m, cmd
}

func (m *model) View() string {
	return baseStyle.Render(m.table.View()) + "\n"
}

func Create() *tea.Program {
	columns := []table.Column{
		{Title: "#", Width: 4},
		{Title: "Method", Width: 10},
		{Title: "Status Code", Width: 15},
		{Title: "URI", Width: 120},
	}

	rows := []table.Row{}

	t := table.New(
		table.WithColumns(columns),
		table.WithRows(rows),
		table.WithFocused(true),
	)

	s := table.DefaultStyles()
	s.Header = s.Header.
		BorderStyle(lipgloss.NormalBorder()).
		BorderForeground(lipgloss.Color("240")).
		BorderBottom(true).
		Bold(false)
	s.Selected = s.Selected.
		Foreground(lipgloss.Color("229")).
		Background(lipgloss.Color("57")).
		Bold(false)
	t.SetStyles(s)

	modelInstance = &model{t}
	program = tea.NewProgram(modelInstance)

	return program
}

func Push(method string, statusCode int, url string) {
	rows := modelInstance.table.Rows()
	position := len(rows) + 1
	program.Send(IOUpdate(table.Row{
		fmt.Sprintf("%d", position),
		method,
		strconv.Itoa(statusCode),
		url,
	}))
}
