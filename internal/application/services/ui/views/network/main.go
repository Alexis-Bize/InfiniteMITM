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
	"infinite-mitm/configs"
	events "infinite-mitm/internal/application/events"
	details "infinite-mitm/internal/application/services/ui/views/network/components/details"
	traffic "infinite-mitm/internal/application/services/ui/views/network/components/details/traffic"
	status "infinite-mitm/internal/application/services/ui/views/network/components/status"
	table "infinite-mitm/internal/application/services/ui/views/network/components/table"
	errors "infinite-mitm/pkg/errors"
	"infinite-mitm/pkg/request"
	utilities "infinite-mitm/pkg/utilities"
	"log"
	"net/url"
	"runtime"
	"strconv"
	"strings"
	"sync"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/gookit/event"
)


type activeKeyType string
type model struct {
	width int
	height int

	networkTableModel   table.TableModel
	networkDetailsModel details.DetailsModel
	statusBarModel      status.StatusBarModel

	activeElement  activeKeyType
}

type networkDataType struct {
	Requests  map[string]*events.ProxyRequestEventData
	Responses map[string]*events.ProxyResponseEventData
}

const (
	NetworkElementKey activeKeyType = "network"
	DetailsElementKey activeKeyType = "details"
)

var program *tea.Program
var programPushRowMutex, programUpdateRowMutex, programUpdateStatusMutex sync.Mutex
var networkData = &networkDataType{
	Requests:  make(map[string]*events.ProxyRequestEventData),
	Responses: make(map[string]*events.ProxyResponseEventData),
}

func Create() {
	width, height := utilities.GetTerminalSize()
	statusBarModel := status.NewStatusBarModel(width)
	modelsHeight := height - statusBarModel.Size.Height

	m := &model{
		width,
		height,
		table.NewNetworkModel(width, modelsHeight),
		details.NewDetailsModel("", "", "", width, modelsHeight),
		statusBarModel,
		NetworkElementKey,
	}

	m.networkTableModel.Focus()

	go func () {
		event.On(events.ProxyRequestSent, event.ListenerFunc(func(e event.Event) error {
			str := fmt.Sprintf("%s", e.Data()["details"])
			data := events.ParseRequestEventData(str)
			pushNetworkData(data)

			return nil
		}), event.Normal)

		event.On(events.ProxyResponseReceived, event.ListenerFunc(func(e event.Event) error {
			str := fmt.Sprintf("%s", e.Data()["details"])
			data := events.ParseResponseEventData(str)
			updateNetworkData(data)

			return nil
		}), event.Normal)

		event.On(events.ProxyStatusMessage, event.ListenerFunc(func(e event.Event) error {
			str := fmt.Sprintf("%s", e.Data()["details"])
			updateStatusBar(str)

			return nil
		}), event.Normal)
	}()

	program = tea.NewProgram(m)
	if _, err := program.Run(); err != nil {
		log.Fatalln(errors.Create(errors.ErrFatalException, err.Error()))
	}
}

func pushNetworkData(data events.ProxyRequestEventData) {
	programPushRowMutex.Lock()
	defer programPushRowMutex.Unlock()

	networkData.Requests[data.ID] = &data

	prefix := ""
	if data.SmartCached {
		prefix = "⧖"
	} else if data.Proxified {
		prefix = "⭑"
	}

	hostname, path := explodeURL(data.URL)
	if hostname == "" || path == "" {
		return
	}

	program.Send(table.TableRowPush(table.TableRowPush{
		ID: data.ID,
		Prefix: prefix,
		Method: data.Method,
		Host: hostname,
		PathAndQuery: path,
	}))

	program.Send(details.RequestTraffic(details.RequestTraffic{
		ID: data.ID,
		Headers: data.Headers,
		Body: data.Body,
	}))
}

func updateNetworkData(data events.ProxyResponseEventData) {
	programUpdateRowMutex.Lock()
	defer programUpdateRowMutex.Unlock()

	networkData.Responses[data.ID] = &data

	prefix := ""
	if data.SmartCached {
		prefix = "⚡"

		var smartCacheHeaderValue string
		for k, v := range data.Headers {
			if strings.EqualFold(k, request.CacheHeaderKey) {
				smartCacheHeaderValue = v
				break
			}
		}

		if smartCacheHeaderValue != request.CacheHeaderHitValue {
			if data.Status == 200 {
				prefix = "⧗"
			} else if data.Status != 302 && (data.Status < 200 || data.Status >= 300) {
				prefix = "✘"
			}
		}
	} else if data.Proxified {
		prefix = "⭑"
	}

	hostname, path := explodeURL(data.URL)
	if hostname == "" || path == "" {
		return
	}

	program.Send(table.TableRowUpdate(table.TableRowUpdate{
		ID: data.ID,
		Prefix: prefix,
		Method: data.Method,
		Host: hostname,
		PathAndQuery: path,
		Status: data.Status,
		ContentType: data.Headers[request.ContentTypeHeaderKey],
	}))

	program.Send(details.ResponseTraffic(details.ResponseTraffic{
		ID: data.ID,
		Headers: data.Headers,
		Body: data.Body,
	}))

	program.Send(details.ResponseStatus(details.ResponseStatus{
		ID: data.ID,
		Status: data.Status,
	}))
}

func updateStatusBar(message string) {
	programUpdateStatusMutex.Lock()
	defer programUpdateStatusMutex.Unlock()

	program.Send(status.StatusBarInfoUpdate(status.StatusBarInfoUpdate{
		Message: message,
	}))
}

func explodeURL(value string) (string, string) {
	parse, err := url.Parse(value)
	if err != nil {
		return "", ""
	}

	path := parse.Path
	if path == "" {
		path = "/"
	}

	path = strings.Replace(path, " ", "%20", -1)
	query := parse.RawQuery
	if query != "" {
		path += "?" + query
	}

	return parse.Hostname(), path
}

func (m *model) Init() tea.Cmd {
	return tea.SetWindowTitle(configs.GetConfig().Name)
}

func (m *model) SetActiveElement(key activeKeyType) {
	m.activeElement = key

	if key == NetworkElementKey {
		m.networkTableModel.Focus()
		m.networkDetailsModel.Blur()
	} else if key == DetailsElementKey {
		m.networkDetailsModel.Focus()
		m.networkTableModel.Blur()
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
	cmds := []tea.Cmd{}

	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height

		m.networkTableModel.SetHeight(
			m.height - m.statusBarModel.Size.Height,
			m.statusBarModel.Size.Height,
		)

		m.statusBarModel.SetWidth(m.width)
		m.networkDetailsModel.SetWidth(m.width)
		m.networkDetailsModel.SetHeight(m.height - m.statusBarModel.Size.Height)

		if runtime.GOOS != "window" {
			utilities.ClearTerminal()
		}
	case tea.KeyMsg:
		switch msg.String() {
		case "q":
			m.SetActiveElement(NetworkElementKey)
			return m, tea.Batch(cmds...)
		case "ctrl+c":
			return m, tea.Quit
		case "enter":
			if m.activeElement != NetworkElementKey {
				break
			}

			rowData := m.networkTableModel.GetSelectedRowData()
			if len(rowData) == 0 {
				return m, tea.Batch(cmds...)
			}

			selectedRowPositionID := rowData[1]

			for v, k := range *m.networkTableModel.GetRowPositionMap() {
				if strconv.Itoa(k) == selectedRowPositionID {
					m.networkDetailsModel.SetID(v)
					emptyRequestData := true

					if req, exists := networkData.Requests[v]; exists {
						m.networkDetailsModel.SetRequestInfo(req.URL, req.Method)
						trafficData := traffic.TrafficData{Headers: req.Headers, Body: req.Body}
						m.networkDetailsModel.SetRequestTrafficData(&trafficData)
						emptyRequestData = false
					}

					if resp, exists := networkData.Responses[v]; exists {
						m.networkDetailsModel.SetResponseStatusCode(resp.Status)
						trafficData := traffic.TrafficData{Headers: resp.Headers, Body: resp.Body}
						m.networkDetailsModel.SetResponseTrafficData(&trafficData)

						if emptyRequestData {
							m.networkDetailsModel.SetRequestInfo(resp.URL, resp.Method)
							m.networkDetailsModel.SetRequestTrafficData(&traffic.TrafficData{Dummy: true})
						}
					} else {
						m.networkDetailsModel.SetResponseStatusCode(0)
						m.networkDetailsModel.SetResponseTrafficData(&traffic.TrafficData{})
					}

					m.SetActiveElement(DetailsElementKey)
					break
				}
			}

			return m, tea.Batch(cmds...)
		}
	}

	m.networkTableModel, cmd = m.networkTableModel.Update(msg)
	cmds = append(cmds, cmd)

	if m.activeElement == DetailsElementKey {
		m.networkDetailsModel, cmd = m.networkDetailsModel.Update(msg)
		cmds = append(cmds, cmd)
	}

	m.statusBarModel, cmd = m.statusBarModel.Update(msg)
	cmds = append(cmds, cmd)

	return m, tea.Batch(cmds...)
}

func (m *model) View() string {
	var activeModel string
	statusBar := m.statusBarModel.View()

	if m.activeElement == NetworkElementKey {
		activeModel = m.networkTableModel.View()
	} else if m.activeElement == DetailsElementKey {
		activeModel = m.networkDetailsModel.View()
	}

	return lipgloss.NewStyle().
		Render(lipgloss.JoinVertical(
			lipgloss.Top,
			activeModel,
			statusBar,
		))
}