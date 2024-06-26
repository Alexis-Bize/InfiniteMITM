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

package MITMApplicationNetworkUI

import (
	eventsService "infinite-mitm/internal/application/services/events"
	details "infinite-mitm/internal/application/ui/network/components/details"
	traffic "infinite-mitm/internal/application/ui/network/components/details/traffic"
	status "infinite-mitm/internal/application/ui/network/components/status"
	table "infinite-mitm/internal/application/ui/network/components/table"
	"infinite-mitm/pkg/errors"
	"infinite-mitm/pkg/request"
	"infinite-mitm/pkg/sysutilities"
	"log"
	"net/url"
	"strconv"
	"strings"
	"sync"
	"time"

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
	Requests  map[string]*eventsService.ProxyRequestEventData
	Responses map[string]*eventsService.ProxyResponseEventData
}

const (
	NetworkElementKey activeKeyType = "network"
	DetailsElementKey activeKeyType = "details"
)

const (
	QuitCommand  = "q"
	StopCommand  = "ctrl+c"
	EnterCommand = "enter"
)
var (
	once    sync.Once
	program *tea.Program
)

var (
	networkDataMutex = &sync.Mutex{}
	networkData = networkDataType{
		Requests:  make(map[string]*eventsService.ProxyRequestEventData),
		Responses: make(map[string]*eventsService.ProxyResponseEventData),
	}
)

func Create() {
	width, height := sysutilities.GetTerminalSize()
	statusBarModel := status.NewStatusBarModel(width)
	modelsHeight := height - statusBarModel.Height

	m := model{
		width,
		height,
		table.NewNetworkModel(modelsHeight),
		details.NewDetailsModel("", "", "", width, modelsHeight),
		statusBarModel,
		NetworkElementKey,
	}

	m.networkTableModel.Focus()

	once.Do(func() {
		event.On(eventsService.ProxyRequestSent, event.ListenerFunc(func(e event.Event) error {
			details := e.Data()["details"].(string)
			data := eventsService.ParseRequestEventData(details)
			go pushNetworkData(data)

			return nil
		}), event.Normal)

		event.On(eventsService.ProxyResponseReceived, event.ListenerFunc(func(e event.Event) error {
			details := e.Data()["details"].(string)
			data := eventsService.ParseResponseEventData(details)
			go updateNetworkData(data)

			return nil
		}), event.Normal)

		event.On(eventsService.ProxyStatusMessage, event.ListenerFunc(func(e event.Event) error {
			details := e.Data()["details"].(string)
			go updateStatusBar(details)

			return nil
		}), event.Normal)
	})

	program = tea.NewProgram(m)
	if _, err := program.Run(); err != nil {
		log.Fatalln(errors.Create(errors.ErrFatalException, err.Error()))
	}
}

func (m *model) setActiveElement(key activeKeyType) {
	m.activeElement = key

	if key == NetworkElementKey {
		m.networkTableModel.Focus()
		m.networkDetailsModel.Blur()
	} else if key == DetailsElementKey {
		m.networkDetailsModel.Focus()
		m.networkTableModel.Blur()
	}
}

func (m model) Init() tea.Cmd {
	return nil
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd
	cmds := []tea.Cmd{}

	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height

		m.networkTableModel.SetHeight(m.height - m.statusBarModel.Height)
		m.statusBarModel.SetWidth(m.width)
		m.networkDetailsModel.SetWidth(m.width)
		m.networkDetailsModel.SetHeight(m.height - m.statusBarModel.Height)

		sysutilities.ClearTerminal()
	case tea.KeyMsg:
		switch msg.String() {
		case QuitCommand:
			if m.activeElement == NetworkElementKey {
				sysutilities.ClearTerminal()
				return m, tea.Quit
			}

			m.setActiveElement(NetworkElementKey)
			return m, tea.Batch(cmds...)
		case StopCommand:
			sysutilities.ClearTerminal()
			return m, tea.Quit
		case EnterCommand:
			if m.activeElement != NetworkElementKey {
				break
			}

			rowData := m.networkTableModel.GetSelectedRowData()
			if len(rowData) == 0 {
				return m, tea.Batch(cmds...)
			}

			selectedRowPositionID := rowData[1]
			rowPositionMap := m.networkTableModel.GetRowPositionMap()

			for v, k := range rowPositionMap {
				if strconv.Itoa(k) != selectedRowPositionID {
					continue
				}

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

				m.setActiveElement(DetailsElementKey)
				break
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

func (m model) View() string {
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

func pushNetworkData(data eventsService.ProxyRequestEventData) {
	networkDataMutex.Lock()
	networkData.Requests[data.ID] = &data
	networkDataMutex.Unlock()

	prefix := ""
	if data.SmartCached {
		prefix = "°"
	} else if data.Proxified {
		prefix = "≈"
	}

	hostname, path := explodeURL(data.URL)
	if hostname == "" || path == "" {
		return
	}

	time.Sleep(50 * time.Millisecond)
	program.Send(details.RequestTraffic(details.RequestTraffic{
		ID: data.ID,
		Headers: data.Headers,
		Body: data.Body,
	}))

	time.Sleep(50 * time.Millisecond)
	program.Send(table.TableRowMsg(table.TableRowMsg{
		ID: data.ID,
		Prefix: prefix,
		Method: data.Method,
		Host: hostname,
		PathAndQuery: path,
	}))
}

func updateNetworkData(data eventsService.ProxyResponseEventData) {
	networkDataMutex.Lock()
	networkData.Responses[data.ID] = &data
	networkDataMutex.Unlock()

	prefix := ""
	if data.SmartCached {
		prefix = "λ"

		var smartCacheHeaderValue string
		for k, v := range data.Headers {
			if strings.EqualFold(k, request.MITMCacheHeaderKey) {
				smartCacheHeaderValue = v
				break
			}
		}

		if smartCacheHeaderValue != request.MITMCacheHeaderHitValue {
			if data.Status == 200 {
				prefix = "∙"
			} else if data.Status != 302 && (data.Status < 200 || data.Status >= 300) {
				prefix = ""
			}
		}
	} else if data.Proxified {
		prefix = "≈"
	}

	hostname, path := explodeURL(data.URL)
	if hostname == "" || path == "" {
		return
	}

	time.Sleep(50 * time.Millisecond)
	program.Send(details.ResponseStatus(details.ResponseStatus{
		ID: data.ID,
		Status: data.Status,
	}))

	time.Sleep(50 * time.Millisecond)
	program.Send(details.ResponseTraffic(details.ResponseTraffic{
		ID: data.ID,
		Headers: data.Headers,
		Body: data.Body,
	}))

	time.Sleep(50 * time.Millisecond)
	program.Send(table.TableRowMsg(table.TableRowMsg{
		ID: data.ID,
		Prefix: prefix,
		Method: data.Method,
		Host: hostname,
		PathAndQuery: path,
		Status: data.Status,
		ContentType: data.Headers[request.ContentTypeHeaderKey],
	}))
}

func updateStatusBar(message string) {
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
