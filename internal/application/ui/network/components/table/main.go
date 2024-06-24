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

package MITMApplicationNetworkUITableComponent

import (
	"fmt"
	"strings"

	"infinite-mitm/pkg/theme"

	"github.com/charmbracelet/bubbles/table"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type TableModel struct {
	height int
	ready  bool

	tableModel table.Model

	rowPositionIDMap map[string]int
	focused   bool
}

type tableRowEvent  struct {
	ID           string
	Prefix       string
	Method       string
	Host         string
	PathAndQuery string

	Status       int
	ContentType  string
}

type TableRowPushMsg   tableRowEvent
type TableRowUpdateMsg tableRowEvent

const (
	PruneRowsCommand = "ctrl+r"
)

func NewNetworkModel(height int) TableModel {
	m := TableModel{
		height: height,
		tableModel: table.New(),
		rowPositionIDMap: map[string]int{},
		focused: true,
	}

	return m
}

func CreateColums(width int) []table.Column {
	columnWidths := []int{2, 5, 8, 8, 25, 35}
	remainingWidth := width
	for _, width := range columnWidths {
		remainingWidth -= width
	}

	columnWidths = append(columnWidths, remainingWidth)
	actualWidths := make([]int, len(columnWidths))
	for i, percent := range columnWidths {
		actualWidths[i] = (width * percent) / 100
	}

	columns := []table.Column{
		{Title: "✎", Width: actualWidths[0]},
		{Title: "#", Width: actualWidths[1]},
		{Title: "Method", Width: actualWidths[2]},
		{Title: "Result", Width: actualWidths[3]},
		{Title: "Host", Width: actualWidths[4]},
		{Title: "Path", Width: actualWidths[5]},
		{Title: "Content Type", Width: actualWidths[6]},
	}

	return columns
}

func (m *TableModel) Focus() {
	m.focused = true
	m.tableModel.Focus()
}

func (m *TableModel) Blur() {
	m.focused = false
	m.tableModel.Blur()
}

func (m *TableModel) SetHeight(height int) {
	m.height = height
	m.tableModel.SetHeight(m.height - 3)
}

func (m TableModel) GetSelectedRowData() table.Row {
	return m.tableModel.SelectedRow()
}

func (m TableModel) GetRowPosition(id string) int {
	position, exists := m.rowPositionIDMap[id]; if !exists {
		return -1
	}

	return position
}

func (m TableModel) GetRowPositionMap() map[string]int {
	return m.rowPositionIDMap
}

func (m TableModel) getNextRowPosition() int {
	return len(m.rowPositionIDMap) + 1
}

func (m *TableModel) pruneRows() {
	m.rowPositionIDMap = map[string]int{}
	m.tableModel.SetRows([]table.Row{})
}

func (m *TableModel) draw(width int) {
	if m.ready {
		return
	}

	columns := CreateColums(width)
	rows := []table.Row{}
	s := table.DefaultStyles()

	s.Header = s.Header.
		BorderStyle(lipgloss.NormalBorder()).
		BorderForeground(theme.ColorGrey).
		BorderBottom(true).
		Bold(false)

	s.Selected = s.Selected.
		Foreground(theme.ColorLightYellow).
		Background(theme.ColorNeonBlue).
		Bold(false)

	m.tableModel.SetColumns(columns)
	m.tableModel.SetRows(rows)
	m.tableModel.SetStyles(s)

	m.ready = true
}

func (m TableModel) Update(msg tea.Msg) (TableModel, tea.Cmd) {
	var cmd tea.Cmd
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		if !m.ready {
			m.draw(msg.Width)
			return m, cmd
		}

		m.tableModel.SetColumns(CreateColums(msg.Width))
	case TableRowPushMsg:
		if !m.ready {
			return m, cmd
		}

		position := m.getNextRowPosition()
		statusCode := "..."
		if msg.Status != 0 {
			statusCode = fmt.Sprintf("%d", msg.Status)
		}

		contentType := "..."
		if msg.ContentType != "" || msg.Status != 0 {
			contentType = strings.Split(msg.ContentType, ";")[0]
		}

		m.rowPositionIDMap[msg.ID] = position
		m.tableModel.SetRows(append(m.tableModel.Rows(), table.Row([]string{
			msg.Prefix,
			fmt.Sprintf("%d", position),
			msg.Method,
			statusCode,
			msg.Host,
			msg.PathAndQuery,
			contentType,
		})))
	case TableRowUpdateMsg:
		if !m.ready {
			return m, cmd
		}

		rows := m.tableModel.Rows()
		index := m.GetRowPosition(msg.ID) - 1

		if index == -1 {
			return m.Update(TableRowPushMsg(msg))
		} else if index >= 0 && index < len(rows) {
			rows[index][0] = msg.Prefix
			rows[index][3] = fmt.Sprintf("%d", msg.Status)
			rows[index][6] = strings.Split(msg.ContentType, ";")[0]
		}
	case tea.KeyMsg:
		if m.ready && m.focused {
			switch msg.String() {
			case PruneRowsCommand:
				m.pruneRows()
			}
		}
	}

	if m.ready {
		m.tableModel, cmd = m.tableModel.Update(msg)
	}

	return m, cmd
}

func (m TableModel) View() string {
	return lipgloss.NewStyle().
		Height(m.height).
		MaxHeight(m.height).
		Margin(1, 2).
		Render(m.tableModel.View())
}
