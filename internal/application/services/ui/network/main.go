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

package MITMApplicationUIServiceNetworkTable

import (
	"fmt"
	events "infinite-mitm/internal/application/events"
	errors "infinite-mitm/pkg/modules/errors"
	"net/http"
	"net/url"
	"os"
	"strconv"
	"strings"

	"github.com/charmbracelet/bubbles/table"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/gookit/event"
	"golang.org/x/term"
)

type NetworkData struct {
	Requests  []*events.ProxyRequestEventData
	Responses []*events.ProxyResponseEventData
}

type model struct {
	table             table.Model
	details           string
	selectedID        int
	currentView       string
}

type tableRowPush struct {
	ID           string
	WithProxy    string
	Method       string
	Host         string
	PathAndQuery string
}

type tableRowUpdate struct {
	ID          string
	WithProxy   string
	Status      int
	ContentType string
}

var modelInstance *model
var program *tea.Program
var networkData = &NetworkData{
	Requests:  make([]*events.ProxyRequestEventData, 0),
	Responses: make([]*events.ProxyResponseEventData, 0),
}

var rowPositionIDMap = make(map[string]string)

var (
	subtle    = lipgloss.AdaptiveColor{Light: "#D9DCCF", Dark: "#383838"}
	highlight = lipgloss.AdaptiveColor{Light: "#874BFD", Dark: "#7D56F4"}
	special   = lipgloss.AdaptiveColor{Light: "#43BF6D", Dark: "#73F59F"}

	divider = lipgloss.NewStyle().
		SetString("•").
		Padding(0, 1).
		Foreground(subtle).
		String()

	colorSuccess = lipgloss.NewStyle().Foreground(lipgloss.AdaptiveColor{Light: "#43BF6D", Dark: "#73F59F"}).Render
	colorWarn = lipgloss.NewStyle().Foreground(lipgloss.AdaptiveColor{Light: "#bf9543", Dark: "#f5be73"}).Render
	colorError =lipgloss.NewStyle().Foreground(lipgloss.AdaptiveColor{Light: "#bf4343", Dark: "#f57373"}).Render

	docStyle = lipgloss.NewStyle().Padding(1, 2, 1, 2)
)

func (m *model) Init() tea.Cmd {
	return nil
}

func Create() {
	p := createProgram()
	initEvents()
	if _, err := p.Run(); err != nil {
		errors.Create(errors.ErrFatalException, err.Error())
		os.Exit(1)
	}
}

func createProgram() *tea.Program {
	modelInstance = &model{
		createNetworkView(),
		"",
		-1,
		"network",
	}

	program = tea.NewProgram(modelInstance, tea.WithAltScreen())
	return program
}

func createNetworkView() table.Model {
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

	return t
}

func initEvents() {
	event.On(events.ProxyRequestSent, event.ListenerFunc(func(e event.Event) error {
		str := fmt.Sprintf("%s", e.Data()["details"])
		data := events.ParseRequestEventData(str)

		networkData.Requests = append(networkData.Requests, &data)
		pushRow(data)

		return nil
	}), event.Normal)

	event.On(events.ProxyResponseReceived, event.ListenerFunc(func(e event.Event) error {
		str := fmt.Sprintf("%s", e.Data()["details"])
		data := events.ParseResponseEventData(str)

		networkData.Responses = append(networkData.Responses, &data)
		updateRow(data)

		return nil
	}), event.Normal)
}

func pushRow(data events.ProxyRequestEventData) {
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

	program.Send(tableRowPush(tableRowPush{
		ID: data.ID,
		WithProxy: withProxy,
		Method: data.Method,
		Host: parse.Hostname(),
		PathAndQuery: path,
	}))
}

func updateRow(data events.ProxyResponseEventData) {
	withProxy := ""
	if data.Proxified {
		withProxy = "✔"
	}

	program.Send(tableRowUpdate(tableRowUpdate{
		ID: data.ID,
		WithProxy: withProxy,
		Status: data.Status,
		ContentType: data.Headers["Content-Type"],
	}))
}

func getNextRowPosition() int {
	rows := modelInstance.table.Rows()
	position := len(rows) + 1
	return position
}

func (m *model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd
	switch msg := msg.(type) {
	case tableRowPush:
		position := fmt.Sprintf("%d", getNextRowPosition())
		rowPositionIDMap[msg.ID] = position
		m.table.SetRows(append(m.table.Rows(), table.Row([]string{
			msg.WithProxy,
			position,
			msg.Method,
			"...",
			msg.Host,
			msg.PathAndQuery,
			"...",
		})))
	case tableRowUpdate:
		position := rowPositionIDMap[msg.ID]
		if position != "" {
			s, _ := strconv.Atoi(position)
			contentType := msg.ContentType

			target := m.table.Rows()[s - 1]
			target[0] = msg.WithProxy
			target[3] = fmt.Sprintf("%d", msg.Status)

			if contentType == "" {
				target[6] = contentType
			} else {
				explode := strings.Split(contentType, ";")
				target[6] = explode[0]
			}
		}
	case tea.KeyMsg:
		switch msg.String() {
		case "esc":
			if m.table.Focused() {
				m.table.Blur()
			} else {
				m.table.Focus()
			}
		case "q":
			if m.currentView == "details" {
				m.currentView = "network"
			}
		case "ctrl+c":
			return m, tea.Quit
		case "enter":
			if len(m.table.SelectedRow()) != 0 {
				doc := strings.Builder{}

				method := lipgloss.NewStyle().Bold(true).Render(m.table.SelectedRow()[2])
				statusCode, _ := strconv.Atoi(m.table.SelectedRow()[3])
				statusText := http.StatusText(statusCode)
				if statusText == "" {
					statusText = "ongoing"
				}

				requestUrl := fmt.Sprintf("https://%s%s", m.table.SelectedRow()[4], m.table.SelectedRow()[5])

				if statusCode >= 200 && statusCode < 300 {
					requestUrl = colorSuccess(requestUrl)
					statusText = colorSuccess(statusText)
				} else if statusCode >= 400 && statusCode < 500 {
					requestUrl = colorWarn(requestUrl)
					statusText = colorWarn(statusText)
				} else if statusCode >= 500 {
					requestUrl = colorError(requestUrl)
					statusText = colorError(statusText)
				}

				doc.WriteString(method + " [" + statusText + "] " + requestUrl)

				m.details = lipgloss.JoinHorizontal(
					lipgloss.Top,
					docStyle.Render(doc.String()),
				)

				m.currentView = "details"
			}
		}
	}

	m.table, cmd = m.table.Update(msg)
	return m, cmd
}

func (m *model) View() string {
	physicalWidth, _, _ := term.GetSize(int(os.Stdout.Fd()))

	if physicalWidth > 0 {
		docStyle = docStyle.MaxWidth(physicalWidth)
	}

	tableStyle := lipgloss.
		NewStyle().
		BorderStyle(lipgloss.NormalBorder()).
		BorderForeground(lipgloss.Color("240"))

	elem := ""

	if m.currentView == "details" {
		elem = m.details
	} else {
		elem = tableStyle.Render(m.table.View())
	}

	render := lipgloss.JoinVertical(
		lipgloss.Top,
		elem,
	)

	return docStyle.Render(render)
}
