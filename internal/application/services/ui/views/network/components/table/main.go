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

package MITMApplicationUIServiceNetworkComponentTable

import (
	"fmt"
	"strings"

	"github.com/charmbracelet/bubbles/table"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type TableModel struct {
	width int
	tableModel table.Model

	rowPositionIDMap *map[string]int
	focused   bool
}

type TableRowPush struct {
	ID           string
	WithProxy    string
	Method       string
	Host         string
	PathAndQuery string
}

type TableRowUpdate struct {
	ID          string
	WithProxy   string
	Status      int
	ContentType string
}

func NewNetworkModel(width int) TableModel {
	columnWidths := []int{
		int(0.02 * float64(width)),  // 2% 		- ✎
		int(0.04 * float64(width)),  // 4%		- #
		int(0.07 * float64(width)),  // 7%		- Method
		int(0.07 * float64(width)),  // 7%		- Result
		int(0.30 * float64(width)),  // 30%		- Host
		int(0.30 * float64(width)),  // 30%		- Path
	}

	remainingWidth := width
	for _, width := range columnWidths {
		remainingWidth -= width
	}

	columnWidths = append(columnWidths, remainingWidth)
	columns := []table.Column{
		{Title: "✎", Width: columnWidths[0]},
		{Title: "#", Width: columnWidths[1]},
		{Title: "Method", Width: columnWidths[2]},
		{Title: "Result", Width: columnWidths[3]},
		{Title: "Host", Width: columnWidths[4]},
		{Title: "Path", Width: columnWidths[5]},
		{Title: "Content Type", Width: columnWidths[6]},
	}

	rows := []table.Row{}

	t := table.New(
		table.WithColumns(columns),
		table.WithRows(rows),
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

	m := TableModel{
		width: width,
		tableModel: t,
		rowPositionIDMap: &map[string]int{},
		focused: false,
	}

	return m
}

func (m *TableModel) Focus() {
	m.focused = true
	m.tableModel.Focus()
}

func (m *TableModel) Blur() {
	m.focused = false
	m.tableModel.Blur()
}

func (m *TableModel) GetNextRowPosition() int {
	rows := m.tableModel.Rows()
	position := len(rows) + 1
	return position
}

func (m *TableModel) GetSelectedRowData() table.Row {
	return m.tableModel.SelectedRow()
}

func (m *TableModel) SetWidth(width int) {
	m.width = width
}

func (m *TableModel) AssignRowPosition(id string, position int) {
	(*m.rowPositionIDMap)[id] = position
}

func (m *TableModel) GetRowPosition(id string) int {
	return (*m.rowPositionIDMap)[id]
}

func (m *TableModel) GetRowPositionMap() *map[string]int {
	return m.rowPositionIDMap
}

func (m TableModel) Update(msg tea.Msg) (TableModel, tea.Cmd) {
	var cmd tea.Cmd
	switch msg := msg.(type) {
	case TableRowPush:
		position := m.GetNextRowPosition()
		m.AssignRowPosition(msg.ID, position)
		m.tableModel.SetRows(append(m.tableModel.Rows(), table.Row([]string{
			msg.WithProxy,
			fmt.Sprintf("%d", position),
			msg.Method,
			"...",
			msg.Host,
			msg.PathAndQuery,
			"...",
		})))
	case TableRowUpdate:
		rows := m.tableModel.Rows()
		position := m.GetRowPosition(msg.ID)
		index := position - 1

		if index >= 0 && index < len(rows) {
			index := position - 1
			contentType := msg.ContentType
			target := rows[index]
			target[0] = msg.WithProxy
			target[3] = fmt.Sprintf("%d", msg.Status)

			if contentType == "" {
				target[6] = contentType
			} else {
				explode := strings.Split(contentType, ";")
				target[6] = explode[0]
			}
		}

		if m.focused {
			m.tableModel.Focus()
		}
	}

	m.tableModel, cmd = m.tableModel.Update(msg)
	return m, cmd
}

func (m TableModel) View() string {
	return m.tableModel.View()
}
