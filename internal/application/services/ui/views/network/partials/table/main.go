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

package MITMApplicationUIServiceNetworkPartialTable

import (
	"fmt"
	"strings"

	"github.com/charmbracelet/bubbles/table"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type TableModel struct {
	Width int
	TableModel table.Model
	RowPositionIDMap map[string]int
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
	columns := []table.Column{
		{Title: "✎", Width: 2},
		{Title: "#", Width: 5},
		{Title: "Method", Width: 10},
		{Title: "Result", Width: 10},
		{Title: "Host", Width: 40},
		{Title: "Path", Width: 50},
		{Title: "Content Type", Width: 40},
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

	m := TableModel{
		Width: width,
		TableModel: t,
		RowPositionIDMap: make(map[string]int),
	}

	return m
}

func (m TableModel) GetNextRowPosition() int {
	rows := m.TableModel.Rows()
	position := len(rows) + 1
	return position
}

func (m *TableModel) SetWidth(width int) {
	m.Width = width
}

func (m TableModel) Init() tea.Cmd {
	return nil
}

func (m *TableModel) AssignRowPosition(id string, position int) {
	m.RowPositionIDMap[id] = position
}

func (m *TableModel) GetRowPosition(id string) int {
	return m.RowPositionIDMap[id]
}

func (m TableModel) Update(msg tea.Msg) (TableModel, tea.Cmd) {
	var cmd tea.Cmd
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.SetWidth(msg.Width)
	case TableRowPush:
		position := m.GetNextRowPosition()
		m.AssignRowPosition(msg.ID, position)
		m.TableModel.SetRows(append(m.TableModel.Rows(), table.Row([]string{
			msg.WithProxy,
			fmt.Sprintf("%d", position),
			msg.Method,
			"...",
			msg.Host,
			msg.PathAndQuery,
			"...",
		})))
	case TableRowUpdate:
		position := m.GetRowPosition(msg.ID)
		if position != 0 {
			contentType := msg.ContentType

			target := m.TableModel.Rows()[position - 1]
			target[0] = msg.WithProxy
			target[3] = fmt.Sprintf("%d", msg.Status)

			if contentType == "" {
				target[6] = contentType
			} else {
				explode := strings.Split(contentType, ";")
				target[6] = explode[0]
			}
		}
	}

	m.TableModel, cmd = m.TableModel.Update(msg)
	return m, cmd
}

func (m *TableModel) View() string {
	return m.TableModel.View()
}
