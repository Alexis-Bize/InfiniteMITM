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
	"fmt"
	traffic "infinite-mitm/internal/application/services/ui/views/network/components/details/traffic"
	"infinite-mitm/pkg/request"
	"infinite-mitm/pkg/sysutilities"
	"infinite-mitm/pkg/theme"
	"infinite-mitm/pkg/utilities"
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

	activeTab   activeTabType
	copyPressed bool
	focused     bool
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

const (
	CopyUrlCommand    = "ctrl+u"
	SwitchViewCommand = "tab"
)

var (
	tabStyle = lipgloss.NewStyle().Padding(0, 1)
	tabsGroupStyle = lipgloss.NewStyle()

	activeTabStyle = lipgloss.NewStyle().
		Padding(0, 1).
		Foreground(theme.ColorLightYellow).
		Background(theme.ColorNeonBlue).Bold(true)

	actionsStyle = lipgloss.NewStyle().
		Padding(0, 1).MarginRight(1).
		Foreground(theme.ColorLight).
		Background(theme.ColorGrey)
)

var (
	requestString      = "Request"
	responseString     = "Response"
	ongoingString      = "Ongoing"
	windowHeightString = "[ Please increase the window height ]"
	copiedString       = "✓ Copied"

	copyUrlString      = fmt.Sprintf("%s: Copy URL to clipboard", CopyUrlCommand)
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
	m.setCopyPress(false)

	m.requestTrafficModel.SetActiveView(traffic.HeadersViewKey)
	m.requestTrafficModel.Focus()

	m.responseTrafficModel.SetActiveView(traffic.HeadersViewKey)
	m.responseTrafficModel.Blur()
}

func (m *DetailsModel) Blur() {
	m.focused = false
	m.activeTab = RequestTabKey
	m.setCopyPress(false)

	m.requestTrafficModel.SetActiveView(traffic.HeadersViewKey)
	m.requestTrafficModel.Blur()

	m.responseTrafficModel.SetActiveView(traffic.HeadersViewKey)
	m.responseTrafficModel.Blur()
}

func (m *DetailsModel) SetHeight(height int) {
	m.height = height
}

func (m *DetailsModel) SetRequestTrafficData(data *traffic.TrafficData) {
	data.URL = m.requestUrl
	m.requestTrafficModel.SetTrafficData(data)
}

func (m *DetailsModel) SetResponseTrafficData(data *traffic.TrafficData) {
	data.URL = m.requestUrl
	m.responseTrafficModel.SetTrafficData(data)
}

func (m *DetailsModel) SetID(trafficID string) {
	m.trafficID = trafficID
}

func (m *DetailsModel) SetRequestInfo(requestUrl string, requestMethod string) {
	m.requestUrl = requestUrl
	m.requestMethod = requestMethod
}

func (m *DetailsModel) SetResponseStatusCode(statusCode int) {
	m.responseStatusCode = statusCode
}

func (m *DetailsModel) SetWidth(width int) {
	m.width = width
	m.requestTrafficModel.SetWidth(width)
	m.responseTrafficModel.SetWidth(width)
}

func (m *DetailsModel) setCopyPress(pressed bool) {
	m.copyPressed = pressed
}

func (m *DetailsModel) setActiveTab(key activeTabType) {
	m.activeTab = key
	m.requestTrafficModel.SetActiveView(traffic.HeadersViewKey)
	m.responseTrafficModel.SetActiveView(traffic.HeadersViewKey)

	if key == RequestTabKey {
		m.requestTrafficModel.Focus()
		m.responseTrafficModel.Blur()
	} else if key == ResponseTabKey {
		m.responseTrafficModel.Focus()
		m.requestTrafficModel.Blur()
	}
}

func (m *DetailsModel) switchActiveTab() {
	if m.activeTab == RequestTabKey {
		m.setActiveTab(ResponseTabKey)
	} else if m.activeTab == ResponseTabKey {
		m.setActiveTab(RequestTabKey)
	}
}

func (m *DetailsModel) copyToClipboard() {
	m.setCopyPress(true)
	sysutilities.CopyToClipboard(request.StripPort(m.requestUrl))
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
		case SwitchViewCommand:
			m.switchActiveTab()
		case CopyUrlCommand:
			m.copyToClipboard()
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
		copyElement := copyUrlString
		if m.copyPressed {
			copyElement = copiedString
		}

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
				lipgloss.NewStyle().
					MarginBottom(1).
					Render(actionsStyle.Render(copyElement)),
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
