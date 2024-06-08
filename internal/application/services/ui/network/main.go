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
	errors "infinite-mitm/pkg/modules/errors"
	"log"
	"net/http"
	"os"
	"os/exec"
	"strconv"
	"strings"

	"github.com/charmbracelet/bubbles/table"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/gookit/event"
)

type model struct {
	networkTable table.Model
	statusBar    StatusBarModel
	detailsView  string
	selectedID   int
	currentView  string
}

const DefaultWidth = 200

var (
	modelInstance *model
	program *tea.Program
)

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

	docStyle = lipgloss.NewStyle().Padding(1, 2, 1, 2).MaxWidth(GetTerminalWidth())
)

func Create() *errors.MITMError {
	modelInstance = &model{createNetworkTable(), createStatusBar(), "", -1, "network"}
	program = tea.NewProgram(modelInstance, tea.WithAltScreen())

	initEvents()
	if _, err := program.Run(); err != nil {
		log.Fatalln(errors.Create(errors.ErrFatalException, err.Error()))
		return nil
	}

	return nil
}

func initEvents() {
	event.On(events.ProxyRequestSent, event.ListenerFunc(func(e event.Event) error {
		if modelInstance == nil {
			return nil
		}

		str := fmt.Sprintf("%s", e.Data()["details"])
		data := events.ParseRequestEventData(str)
		networkData.Requests = append(networkData.Requests, &data)
		pushNetworkTableRow(data)

		return nil
	}), event.Normal)

	event.On(events.ProxyResponseReceived, event.ListenerFunc(func(e event.Event) error {
		if modelInstance == nil {
			return nil
		}

		str := fmt.Sprintf("%s", e.Data()["details"])
		data := events.ParseResponseEventData(str)
		networkData.Responses = append(networkData.Responses, &data)
		updateNetworkTableRow(data)

		return nil
	}), event.Normal)

	event.On(events.ProxyStatusMessage, event.ListenerFunc(func(e event.Event) error {
		if modelInstance == nil {
			return nil
		}

		str := fmt.Sprintf("%s", e.Data()["details"])
		updateStatusBarInfo(str)

		return nil
	}), event.Normal)
}

func GetTerminalWidth() int {
	cmd := exec.Command("tput", "cols")
	cmd.Stdin = os.Stdin
	cmd.Stderr = os.Stderr
	out, err := cmd.Output()
	if err != nil {
		return DefaultWidth
	}

	width, err := strconv.Atoi(strings.TrimSpace(string(out)))
	if err != nil {
		return DefaultWidth
	}

	return width
}

func (m *model) Init() tea.Cmd {
	return nil
}

func (m *model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd
	switch msg := msg.(type) {
	case StatusBarInfoUpdate:
		m.statusBar.info.content = msg.Message
	case TableRowPush:
		position := fmt.Sprintf("%d", getNextRowPosition())
		rowPositionIDMap[msg.ID] = position
		m.networkTable.SetRows(append(m.networkTable.Rows(), table.Row([]string{
			msg.WithProxy,
			position,
			msg.Method,
			"...",
			msg.Host,
			msg.PathAndQuery,
			"...",
		})))
	case TableRowUpdate:
		position := rowPositionIDMap[msg.ID]
		if position != "" {
			s, _ := strconv.Atoi(position)
			contentType := msg.ContentType

			target := m.networkTable.Rows()[s - 1]
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
		case "q":
			if m.currentView == "details" {
				m.currentView = "network"
			}
		case "ctrl+c":
			return m, tea.Quit
		case "enter":
			if len(m.networkTable.SelectedRow()) != 0 {
				doc := strings.Builder{}

				method := lipgloss.NewStyle().Bold(true).Render(m.networkTable.SelectedRow()[2])
				statusCode, _ := strconv.Atoi(m.networkTable.SelectedRow()[3])
				statusText := http.StatusText(statusCode)
				if statusText == "" {
					statusText = "ongoing"
				}

				requestUrl := fmt.Sprintf("https://%s%s", m.networkTable.SelectedRow()[4], m.networkTable.SelectedRow()[5])

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

				m.detailsView = lipgloss.JoinHorizontal(
					lipgloss.Top,
					docStyle.Render(doc.String()),
				)

				m.currentView = "details"
			}
		}
	}

	m.networkTable, cmd = m.networkTable.Update(msg)
	return m, cmd
}

func (m *model) View() string {
	statusBarView := lipgloss.NewStyle().PaddingTop(1).Render(m.statusBar.View())
	activeView := ""

	if m.currentView == "details" {
		activeView = lipgloss.NewStyle().Render(m.detailsView)
	} else {
		activeView = lipgloss.NewStyle().Render(m.networkTable.View())
	}

	return docStyle.Render(lipgloss.JoinVertical(
		lipgloss.Top,
		activeView,
		statusBarView,
	))
}
