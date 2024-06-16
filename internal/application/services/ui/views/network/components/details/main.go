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

package MITMApplicationUIServiceNetworkDetailsComponent

import (
	traffic "infinite-mitm/internal/application/services/ui/views/network/components/details/traffic"
	"infinite-mitm/pkg/request"
	"infinite-mitm/pkg/theme"
	utilities "infinite-mitm/pkg/utilities"
	"net/http"
	"strings"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)


type activeTabType string
type DetailsModel struct {
	width int
	height int

	requestTrafficModel  traffic.TrafficModel
	responseTrafficModel traffic.TrafficModel

	trafficID     string
	requestMethod string
	requestUrl    string

	responseStatusCode  int

	activeTab activeTabType
	focused   bool
}

type RequestTraffic struct {
	ID      string
	Headers map[string] string
	Body    []byte
}

type ResponseTraffic struct {
	ID      string
	Headers map[string] string
	Body    []byte
}

type ResponseStatus struct {
	ID     string
	Status int
}

const (
	RequestTabKey  activeTabType = "request"
	ResponseTabKey activeTabType = "response"
)

var (
	tabStyle = lipgloss.NewStyle().Padding(0, 1)
	tabsGroupStyle = lipgloss.NewStyle()

	activeTabStyle = lipgloss.NewStyle().
		Padding(0, 1).
		Foreground(theme.ColorLightYellow).
		Background(theme.ColorNeonBlue).Bold(true)
)

var (
	requestString      = "Request"
	responseString     = "Response"
	ongoingString      = "Ongoing"
	windowHeightString = "Please increase the window height."
)


func NewDetailsModel(id string, method string, url string, width int, height int) DetailsModel {
	m := DetailsModel{
		trafficID: id,
		requestMethod: method,
		requestUrl: url,
		requestTrafficModel: traffic.NewTrafficDetailsModel(width, 15),
		responseTrafficModel: traffic.NewTrafficDetailsModel(width, 15),
		activeTab: RequestTabKey,
		focused: false,
		width: width,
		height: height,
	}

	return m
}

func (m *DetailsModel) Focus() {
	m.focused = true
	m.activeTab = RequestTabKey
	m.requestTrafficModel.Focus()
	m.responseTrafficModel.Blur()

	m.requestTrafficModel.SetActiveView(traffic.HeadersViewKey)
	m.responseTrafficModel.SetActiveView(traffic.HeadersViewKey)
}

func (m *DetailsModel) Blur() {
	m.focused = false
	m.activeTab = RequestTabKey
	m.requestTrafficModel.Blur()
	m.responseTrafficModel.Blur()

	m.requestTrafficModel.SetActiveView(traffic.HeadersViewKey)
	m.responseTrafficModel.SetActiveView(traffic.HeadersViewKey)
}

func (m *DetailsModel) SetHeight(height int) {
	m.height = height
}

func (m *DetailsModel) SetActiveTab(key activeTabType) {
	m.activeTab = key
	m.requestTrafficModel.SetActiveView(traffic.HeadersViewKey)
	m.responseTrafficModel.SetActiveView(traffic.HeadersViewKey)

	if key == RequestTabKey {
		m.responseTrafficModel.Blur()
	} else if key == ResponseTabKey {
		m.requestTrafficModel.Blur()
	}
}

func (m *DetailsModel) SwitchActiveTab() {
	if m.activeTab == RequestTabKey {
		m.SetActiveTab(ResponseTabKey)
	} else if m.activeTab == ResponseTabKey {
		m.SetActiveTab(RequestTabKey)
	}
}

func (m *DetailsModel) SetRequestTrafficData(data *traffic.TrafficData) {
	data.URL = m.requestUrl
	m.requestTrafficModel.SetTrafficData(data)
}

func (m *DetailsModel) SetResponseTrafficData(data *traffic.TrafficData) {
	data.URL = m.requestUrl
	m.responseTrafficModel.SetTrafficData(data)
}

func (m *DetailsModel) SetID(activeID string) {
	m.trafficID = activeID
}

func (m *DetailsModel) SetRequestInfo(activeURL string, activeMethod string) {
	m.requestUrl = activeURL
	m.requestMethod = activeMethod
}

func (m *DetailsModel) SetResponseStatusCode(statusCode int) {
	m.responseStatusCode = statusCode
}

func (m *DetailsModel) SetWidth(width int) {
	m.width = width
	m.requestTrafficModel.SetWidth(width)
	m.responseTrafficModel.SetWidth(width)
}

func (m DetailsModel) Update(msg tea.Msg) (DetailsModel, tea.Cmd) {
	var cmd tea.Cmd
	cmds := []tea.Cmd{}

	switch msg := msg.(type) {
	case RequestTraffic:
		if m.focused && msg.ID == m.trafficID {
			m.SetRequestTrafficData(&traffic.TrafficData{Headers: msg.Headers, Body: msg.Body})
		}

		return m, tea.Batch(cmds...)
	case ResponseTraffic:
		if m.focused && msg.ID == m.trafficID {
			m.SetResponseTrafficData(&traffic.TrafficData{Headers: msg.Headers, Body: msg.Body})
		}

		return m, tea.Batch(cmds...)
	case ResponseStatus:
		if m.focused && msg.ID == m.trafficID {
			m.SetResponseStatusCode(msg.Status)
		}

		return m, tea.Batch(cmds...)
	case tea.KeyMsg:
		if !m.focused {
			return m, tea.Batch(cmds...)
		}

		switch msg.String() {
		case "tab":
			m.SwitchActiveTab()
		}
	}

	m.requestTrafficModel, cmd = m.requestTrafficModel.Update(msg)
	cmds = append(cmds, cmd)

	m.responseTrafficModel, cmd = m.responseTrafficModel.Update(msg)
	cmds = append(cmds, cmd)

	return m, tea.Batch(cmds...)
}

func (m DetailsModel) View() string {
	var render string
	var trafficView string
	var tabs []string

	if m.height > 15 {
		url := request.StripPort(m.requestUrl)
		method := lipgloss.NewStyle().Bold(true).Render(m.requestMethod)
		statusCode := m.responseStatusCode
		statusText := http.StatusText(statusCode)
		if statusText == "" {
			statusText = ongoingString
		}

		urlContent := theme.StatusCodeToColorStyle(statusCode).Render(url)
		statusTextContent := theme.StatusCodeToColorStyle(statusCode).Render(statusText)
		details := method + " " + "[" + statusTextContent + "]" + " " + utilities.WrapText(urlContent, m.width - 2)

		if m.activeTab == RequestTabKey {
			tabs = append(tabs, activeTabStyle.Render(requestString), tabStyle.Render(responseString))
			trafficView = m.requestTrafficModel.View()
		} else if m.activeTab == ResponseTabKey {
			tabs = append(tabs, tabStyle.Render(requestString), activeTabStyle.Render(responseString))
			trafficView = m.responseTrafficModel.View()
		}

		render = lipgloss.NewStyle().
			Render(lipgloss.JoinVertical(
				lipgloss.Top,
				lipgloss.NewStyle().
					MarginBottom(1).
					Render(details),
				lipgloss.JoinHorizontal(
					lipgloss.Top,
					tabsGroupStyle.Render(strings.Join(tabs, " ")),
					lipgloss.NewStyle().
						MarginLeft(2).
						Foreground(theme.ColorGrey).
						Render("Tab ↹: Switch"),
				),
				lipgloss.NewStyle().
					Padding(1, 2).
					Render(trafficView),
			))
	}

	if render == "" {
		render = lipgloss.NewStyle().
			Foreground(theme.ColorGrey).
			Render(windowHeightString)
	}

	return lipgloss.NewStyle().
		Height(m.height).
		MaxHeight(m.height).
		Margin(1, 2).
		Render(render)
}
