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

package MITMApplicationUIServiceTestTable

import (
	"encoding/json"
	"fmt"
	events "infinite-mitm/internal/application/events"
	errors "infinite-mitm/pkg/modules/errors"
	"net/http"
	"net/url"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/charmbracelet/bubbles/table"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/gookit/event"
	"golang.org/x/term"
)

type RequestData struct {
	ID        string `json:"id"`
	URL       string `json:"url"`
	Method    string `json:"method"`
	Headers   map[string]string `json:"headers"`
	Body      []byte `json:"body"`
	Proxified bool `json:"proxified"`
}

type ResponseData struct {
	ID     string `json:"id"`
	Status int `json:"status"`
	Headers map[string]string `json:"headers"`
	Body []byte `json:"body"`
}

type NetworkData struct {
	Request  RequestData
	Response ResponseData
}

type model struct {
	table             table.Model
	details           string
	selectedID        int
	currentView       string
}

type tableRowPush table.Row
type tableRowUpdate table.Row

var program *tea.Program
var modelInstance *model
var networkData = make(map[string]NetworkData)

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

func Start() {
	initEvents()
	p := createProgram()
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

	program = tea.NewProgram(modelInstance)
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
		if modelInstance == nil || program == nil {
			return nil
		}

		s := fmt.Sprintf("%s", e.Data()["data"])
		var data RequestData
		json.Unmarshal([]byte(s), &data)

		networkData[data.ID] = NetworkData{
			Request: data,
		}

		position := len(modelInstance.table.Rows()) + 1

		proxified := ""
		if data.Proxified {
			proxified = "✔"
		}

		parse, _ := url.Parse(data.URL)
		program.Send(tableRowPush(table.Row{
			proxified,
			fmt.Sprintf("%d", position),
			data.Method,
			"...",
			parse.Host,
			parse.Path,
			"...",
		}))

		return nil
	}), event.Normal)

	go func () {
		i := 0

		for {
			i++

			data1 := RequestData{
				ID: fmt.Sprint("%d", i),
				Method: "GET",
				URL: "https://example.com/foo",
				Headers: map[string]string{
					"Accept": "application/json",
				},
				Body: nil,
				Proxified: true,
			}

			mashal, _ := json.Marshal(data1)
			event.MustFire(events.ProxyRequestSent, event.M{"data": mashal})
			time.Sleep(5 * time.Second)

			data2 := ResponseData{
				Status: 200,
				Headers: map[string]string{
					"Content-Type": "application/json",
				},
				Body: nil,
			}

			mashal, _ = json.Marshal(data2)
			event.MustFire(events.ProxyResponseReceived, event.M{"data": mashal})
		}
	}()
}

func (m *model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd
	switch msg := msg.(type) {
	case tableRowPush:
		m.table.SetRows(append(m.table.Rows(), table.Row(msg)))
	case tableRowUpdate:
		// TODO
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
