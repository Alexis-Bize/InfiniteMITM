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

package MITMApplicationUIServiceNetworkUI

import (
	"fmt"
	events "infinite-mitm/internal/application/events"
	details "infinite-mitm/internal/application/services/ui/views/network/partials/details"
	table "infinite-mitm/internal/application/services/ui/views/network/partials/table"
	errors "infinite-mitm/pkg/modules/errors"
	utilities "infinite-mitm/pkg/modules/utilities"
	"log"
	"net/url"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/gookit/event"
)

type model struct {
	Width int

	tableModel   table.TableModel
	detailsModel details.DetailsModel

	currentView  string
}

type NetworkData struct {
	Requests  []*events.ProxyRequestEventData
	Responses []*events.ProxyResponseEventData
}

var (
	modelInstance *model
	program *tea.Program
)

var networkData = &NetworkData{
	Requests:  make([]*events.ProxyRequestEventData, 0),
	Responses: make([]*events.ProxyResponseEventData, 0),
}

var (
	docStyle = lipgloss.NewStyle().Padding(1, 2, 1, 2)
)

func Create() {
	width := utilities.GetTerminalWidth()
	modelInstance = &model{
		width,
		table.NewNetworkModel(width),
		details.NewDetailsModel("", "", "", width),
		"network",
	}

	event.On(events.ProxyRequestSent, event.ListenerFunc(func(e event.Event) error {
		if modelInstance == nil {
			return nil
		}

		str := fmt.Sprintf("%s", e.Data()["details"])
		data := events.ParseRequestEventData(str)
		networkData.Requests = append(networkData.Requests, &data)
		PushNetworkTableRow(data)

		return nil
	}), event.Normal)

	event.On(events.ProxyResponseReceived, event.ListenerFunc(func(e event.Event) error {
		if modelInstance == nil {
			return nil
		}

		str := fmt.Sprintf("%s", e.Data()["details"])
		data := events.ParseResponseEventData(str)
		networkData.Responses = append(networkData.Responses, &data)
		UpdateNetworkTableRow(data)

		return nil
	}), event.Normal)

	event.On(events.ProxyStatusMessage, event.ListenerFunc(func(e event.Event) error {
		if modelInstance == nil {
			return nil
		}

		// str := fmt.Sprintf("%s", e.Data()["details"])

		return nil
	}), event.Normal)

	program = tea.NewProgram(modelInstance, tea.WithAltScreen())
	if _, err := program.Run(); err != nil {
		log.Fatalln(errors.Create(errors.ErrFatalException, err.Error()))
	}
}

func PushNetworkTableRow(data events.ProxyRequestEventData) {
	parse, _ := url.Parse(data.URL)
	withProxy := ""
	if data.Proxified {
		withProxy = "✔"
	}

	path := parse.Path
	if path == "" {
		path = "/"
	}

	query := parse.RawQuery
	if query != "" {
		path += "?" + query
	}

	program.Send(table.TableRowPush(table.TableRowPush{
		ID: data.ID,
		WithProxy: withProxy,
		Method: data.Method,
		Host: parse.Hostname(),
		PathAndQuery: path,
	}))
}

func UpdateNetworkTableRow(data events.ProxyResponseEventData) {
	withProxy := ""
	if data.Proxified {
		withProxy = "✔"
	}

	program.Send(table.TableRowUpdate(table.TableRowUpdate{
		ID: data.ID,
		WithProxy: withProxy,
		Status: data.Status,
		ContentType: data.Headers["Content-Type"],
	}))
}

func (m *model) Init() tea.Cmd {
	return nil
}

func (m *model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "q":
			if m.currentView == "details" {
				m.tableModel.Table.Focus()
				m.detailsModel.Blur()
				m.currentView = "network"
			} else {
				return m, tea.Quit
			}
		case "ctrl+c":
			return m, tea.Quit
		case "enter":
			if m.currentView == "network" && m.tableModel.Table.SelectedRow() != nil {
				m.tableModel.Table.Blur()
				m.detailsModel.Focus()
				m.currentView = "details"
			}
		}
	}

	cmds := []tea.Cmd{cmd}

	m.tableModel, cmd = m.tableModel.Update(msg)
	cmds = append(cmds, cmd)

	m.detailsModel, cmd = m.detailsModel.Update(msg)
	cmds = append(cmds, cmd)

	return m, tea.Batch(cmds...)
}

func (m *model) View() string {
	var activeView string

	if m.currentView == "network" {
		activeView = lipgloss.NewStyle().Render(m.tableModel.View())
	} else if m.currentView == "details" {
		activeView = lipgloss.NewStyle().Render(m.detailsModel.View())
	}

	return docStyle.Render(lipgloss.JoinVertical(
		lipgloss.Top,
		activeView,
	))
}
