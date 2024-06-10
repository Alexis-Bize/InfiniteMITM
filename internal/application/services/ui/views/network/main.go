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
	traffic "infinite-mitm/internal/application/services/ui/views/network/partials/details/traffic"
	table "infinite-mitm/internal/application/services/ui/views/network/partials/table"
	errors "infinite-mitm/pkg/modules/errors"
	utilities "infinite-mitm/pkg/modules/utilities"
	"log"
	"net/url"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/gookit/event"
)


type activeModelType string
type model struct {
	Width int

	networkTableModel   table.TableModel
	networkDetailsModel details.DetailsModel

	activeModel  activeModelType
}

type NetworkData struct {
	Requests  []*events.ProxyRequestEventData
	Responses []*events.ProxyResponseEventData
}

const (
	Network activeModelType = "network"
	Details activeModelType = "details"
)

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
		Network,
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

func (m *model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "q":
			if m.activeModel == Details {
				m.networkTableModel.TableModel.Focus()
				m.networkDetailsModel.Blur()
				m.activeModel = Network
			} else {
				return m, tea.Quit
			}
		case "ctrl+c":
			return m, tea.Quit
		case "enter":
			if m.activeModel == Network && m.networkTableModel.TableModel.SelectedRow() != nil {
				go func() {
					for v, k := range m.networkTableModel.RowPositionIDMap {
						if fmt.Sprintf("%d", k) == m.networkTableModel.TableModel.SelectedRow()[1] {
							m.networkDetailsModel.SetID(v)

							for _, k := range networkData.Requests {
								if k.ID == v {
									m.networkDetailsModel.SetRequestInfo(k.URL, k.Method)
									trafficData := traffic.TrafficData{Headers: k.Headers, Body: k.Body}
									m.networkDetailsModel.RequestTrafficModel.SetTrafficData(trafficData)
									break
								}
							}

							for _, k := range networkData.Responses {
								if k.ID == v {
									m.networkDetailsModel.SetResponseInfo(k.Status)
									trafficData := traffic.TrafficData{Headers: k.Headers, Body: k.Body}
									m.networkDetailsModel.ResponseTrafficModel.SetTrafficData(trafficData)
									break
								}
							}

							break
						}
					}

					m.networkTableModel.TableModel.Blur()
					m.networkDetailsModel.Focus()
					m.activeModel = Details
				}()
			}
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

	if m.activeModel == Network {
		activeModel = lipgloss.NewStyle().Render(m.networkTableModel.View())
	} else if m.activeModel == Details {
		activeModel = lipgloss.NewStyle().Render(m.networkDetailsModel.View())
	}

	return docStyle.Render(lipgloss.JoinVertical(
		lipgloss.Top,
		activeModel,
	))
}
