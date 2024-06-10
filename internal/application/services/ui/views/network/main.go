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
	details "infinite-mitm/internal/application/services/ui/views/network/components/details"
	traffic "infinite-mitm/internal/application/services/ui/views/network/components/details/traffic"
	table "infinite-mitm/internal/application/services/ui/views/network/components/table"
	errors "infinite-mitm/pkg/modules/errors"
	utilities "infinite-mitm/pkg/modules/utilities"
	"log"
	"net/url"
	"strconv"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/gookit/event"
)


type activeKeyType string
type model struct {
	Width int

	networkTableModel   table.TableModel
	networkDetailsModel details.DetailsModel

	activeElement  activeKeyType
}

type networkData struct {
	Requests  map[string]*events.ProxyRequestEventData
	Responses map[string]*events.ProxyResponseEventData
}

const (
	NetworkElementKey activeKeyType = "network"
	DetailsElementKey activeKeyType = "details"
)

var program *tea.Program
var activeNetworkData = &networkData{
	Requests:  make(map[string]*events.ProxyRequestEventData),
	Responses: make(map[string]*events.ProxyResponseEventData),
}

var (
	docStyle = lipgloss.NewStyle().Padding(1, 2, 1, 2)
)

func Create() {
	width := utilities.GetTerminalWidth()
	m := &model{width, table.NewNetworkModel(width), details.NewDetailsModel("", "", "", width), NetworkElementKey}

	event.On(events.ProxyRequestSent, event.ListenerFunc(func(e event.Event) error {
		str := fmt.Sprintf("%s", e.Data()["details"])
		data := events.ParseRequestEventData(str)
		activeNetworkData.Requests[data.ID] = &data
		PushNetworkTableRow(data)

		return nil
	}), event.Normal)

	event.On(events.ProxyResponseReceived, event.ListenerFunc(func(e event.Event) error {
		str := fmt.Sprintf("%s", e.Data()["details"])
		data := events.ParseResponseEventData(str)
		activeNetworkData.Responses[data.ID] = &data
		UpdateNetworkTableRow(data)

		return nil
	}), event.Normal)

	event.On(events.ProxyStatusMessage, event.ListenerFunc(func(e event.Event) error {
		// str := fmt.Sprintf("%s", e.Data()["details"])

		return nil
	}), event.Normal)

	program = tea.NewProgram(m, tea.WithAltScreen())
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

	program.Send(details.RequestTraffic(details.RequestTraffic{
		ID: data.ID,
		Headers: data.Headers,
		Body: data.Body,
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

	program.Send(details.ResponseTraffic(details.ResponseTraffic{
		ID: data.ID,
		Headers: data.Headers,
		Body: data.Body,
	}))
}

func (m *model) Init() tea.Cmd {
	return nil
}

func (m *model) SetActiveElement(key activeKeyType) {
	m.activeElement = key

	if m.activeElement == NetworkElementKey {
		m.networkTableModel.TableModel.Focus()
		m.networkDetailsModel.Blur()
	} else if m.activeElement == DetailsElementKey {
		m.networkDetailsModel.Focus()
		m.networkTableModel.TableModel.Blur()
	}
}

func (m *model) SwitchActiveElement() {
	if m.activeElement == NetworkElementKey {
		m.SetActiveElement(DetailsElementKey)
	} else if m.activeElement == DetailsElementKey {
		m.SetActiveElement(NetworkElementKey)
	}
}

func (m *model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		utilities.ClearTerminal()
	case tea.KeyMsg:
		switch msg.String() {
		case "q", "ctrl+c":
			if m.activeElement == NetworkElementKey {
				return m, tea.Quit
			}

			m.SwitchActiveElement()
		case "enter":
			if len(m.networkTableModel.TableModel.SelectedRow()) < 2 {
				break
			} else if m.activeElement != NetworkElementKey {
				break
			}

			selectedRowID := m.networkTableModel.TableModel.SelectedRow()[1]

			for v, k := range *m.networkTableModel.RowPositionIDMap {
				if strconv.Itoa(k) == selectedRowID {
					m.networkDetailsModel.SetID(v)

					if req, exists := activeNetworkData.Requests[v]; exists {
						m.networkDetailsModel.SetRequestInfo(req.URL, req.Method)
						trafficData := traffic.TrafficData{Headers: req.Headers, Body: req.Body}
						m.networkDetailsModel.SetRequestTrafficData(trafficData)
					}

					if resp, exists := activeNetworkData.Responses[v]; exists {
						m.networkDetailsModel.SetResponseInfo(resp.Status)
						trafficData := traffic.TrafficData{Headers: resp.Headers, Body: resp.Body}
						m.networkDetailsModel.SetResponseTrafficData(trafficData)
					}

					break
				}
			}

			m.SetActiveElement(DetailsElementKey)
		}
	}

	cmds := []tea.Cmd{cmd}

	m.networkTableModel, cmd = m.networkTableModel.Update(msg)
	cmds = append(cmds, cmd)

	m.networkDetailsModel, cmd = m.networkDetailsModel.Update(msg)
	cmds = append(cmds, cmd)

	return m, tea.Batch(cmds...)
}

func (m *model) View() string {
	var activeModel string

	if m.activeElement == NetworkElementKey {
		activeModel = lipgloss.NewStyle().Render(m.networkTableModel.View())
	} else if m.activeElement == DetailsElementKey {
		activeModel = lipgloss.NewStyle().Render(m.networkDetailsModel.View())
	}

	return docStyle.Render(lipgloss.JoinVertical(
		lipgloss.Top,
		activeModel,
	))
}
