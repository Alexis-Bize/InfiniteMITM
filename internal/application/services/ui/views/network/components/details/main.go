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
	theme "infinite-mitm/internal/application/services/ui/theme"
	traffic "infinite-mitm/internal/application/services/ui/views/network/components/details/traffic"
	"net/http"
	"strings"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)


type activeTabType string
type DetailsModel struct {
	Width int

	requestTrafficModel  traffic.TrafficModel
	responseTrafficModel traffic.TrafficModel

	trafficID     string
	requestMethod string
	requestUrl    string
	responseCode  int

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

const (
	RequestTabKey  activeTabType = "request"
	ResponseTabKey activeTabType = "response"
)

func NewDetailsModel(id string, method string, url string, width int) DetailsModel {
	m := DetailsModel{
		trafficID: id,
		requestMethod: method,
		requestUrl: url,
		requestTrafficModel: traffic.NewTrafficDetailsModel(width, 10),
		responseTrafficModel: traffic.NewTrafficDetailsModel(width, 10),
		activeTab: RequestTabKey,
	}

	return m
}

func (m *DetailsModel) Focus() {
	m.focused = true
	m.activeTab = RequestTabKey
	m.requestTrafficModel.Focus()
	m.responseTrafficModel.Blur()
}

func (m *DetailsModel) Blur() {
	m.focused = false
	m.activeTab = RequestTabKey
	m.requestTrafficModel.Blur()
	m.responseTrafficModel.Blur()
}

func (m *DetailsModel) SetActiveTab(key activeTabType) {
	m.activeTab = key
	m.requestTrafficModel.SetActiveView(traffic.HeadersViewKey)
	m.responseTrafficModel.SetActiveView(traffic.HeadersViewKey)

	if m.activeTab == RequestTabKey {
		m.requestTrafficModel.Focus()
		m.responseTrafficModel.Blur()
	} else if m.activeTab == ResponseTabKey {
		m.responseTrafficModel.Focus()
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

func (m *DetailsModel) SetRequestTrafficData(data traffic.TrafficData) {
	m.requestTrafficModel.SetTrafficData(traffic.TrafficData{
		Headers: data.Headers,
		Body: data.Body,
	})
}

func (m *DetailsModel) SetResponseTrafficData(data traffic.TrafficData) {
	m.responseTrafficModel.SetTrafficData(traffic.TrafficData{
		Headers: data.Headers,
		Body: data.Body,
	})
}

func (m *DetailsModel) SetID(activeID string) {
	m.trafficID = activeID
}

func (m *DetailsModel) SetRequestInfo(activeURL string, activeMethod string) {
	m.requestUrl = activeURL
	m.requestMethod = activeMethod
}

func (m *DetailsModel) SetResponseInfo(activeResponseCode int) {
	m.responseCode = activeResponseCode
}

func (m *DetailsModel) SetWidth(width int) {
	m.Width = width
}

func (m DetailsModel) Update(msg tea.Msg) (DetailsModel, tea.Cmd) {
	var cmd tea.Cmd
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.SetWidth(msg.Width)
	case RequestTraffic:
		if msg.ID == m.trafficID {
			m.SetRequestTrafficData(traffic.TrafficData{Headers: msg.Headers, Body: msg.Body})
		}
	case ResponseTraffic:
		if msg.ID == m.trafficID {
			m.SetResponseTrafficData(traffic.TrafficData{Headers: msg.Headers, Body: msg.Body})
		}
	case tea.KeyMsg:
		if m.focused {
			switch msg.String() {
			case "tab":
				m.SwitchActiveTab()
			}
		}
	}

	cmds := []tea.Cmd{cmd}

	m.requestTrafficModel, cmd = m.requestTrafficModel.Update(msg)
	cmds = append(cmds, cmd)

	m.responseTrafficModel, cmd = m.responseTrafficModel.Update(msg)
	cmds = append(cmds, cmd)

	return m, tea.Batch(cmds...)
}

func (m DetailsModel) View() string {
	var trafficView string
	var tabs []string

	detailsStyle := lipgloss.NewStyle().MarginBottom(1)
	switchHintStyle := lipgloss.NewStyle().MarginLeft(1)

	tabStyle := lipgloss.NewStyle().MarginRight(2)
	activeTabStyle := lipgloss.NewStyle().MarginRight(2).Foreground(theme.ColorSunsetOrange).Bold(true).Underline(true)
	tabsGroupStyle := lipgloss.NewStyle().Padding(0, 1)

	trafficViewStyle := lipgloss.NewStyle().Padding(1, 2)

	methodStyle := lipgloss.NewStyle().Bold(true)

	url := m.requestUrl
	method := methodStyle.Render(m.requestMethod)
	statusText := http.StatusText(m.responseCode)
	if statusText == "" {
		statusText = "Ongoing"
	}

	if m.responseCode >= 200 && m.responseCode < 300 {
		successColor := lipgloss.NewStyle().Foreground(theme.ColorSuccess)
		url = successColor.Render(url)
		statusText = successColor.Render(statusText)
	} else if m.responseCode >= 400 && m.responseCode < 500 {
		warnColor := lipgloss.NewStyle().Foreground(theme.ColorWarn)
		url = warnColor.Render(url)
		statusText =  warnColor.Render(statusText)
	} else if m.responseCode >= 500 {
		errorColor := lipgloss.NewStyle().Foreground(theme.ColorError)
		url = errorColor.Render(url)
		statusText = errorColor.Render(statusText)
	} else {
		otherColor := lipgloss.NewStyle().Foreground(theme.ColorGrey)
		url = otherColor.Render(url)
		statusText = otherColor.Render(statusText)
	}

	details := method + " " + "[" + statusText + "]" + " " + url
	switchHint := "Use Tab ↹ to switch between request and response"

	if m.activeTab == RequestTabKey {
		tabs = append(tabs, activeTabStyle.Render("Request"), tabStyle.Render("Response"))
		trafficView = m.requestTrafficModel.View()
	} else if m.activeTab == ResponseTabKey {
		tabs = append(tabs, tabStyle.Render("Request"), activeTabStyle.Render("Response"))
		trafficView = m.responseTrafficModel.View()
	}

	return lipgloss.NewStyle().MaxWidth(m.Width).Render(lipgloss.JoinVertical(
		lipgloss.Top,
		detailsStyle.Render(details),
		tabsGroupStyle.Render(strings.Join(tabs, " ")),
		switchHintStyle.Render(switchHint),
		trafficViewStyle.Render(trafficView),
	))
}
