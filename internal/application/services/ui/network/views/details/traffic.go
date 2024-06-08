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

package MITMApplicationUIServiceNetworkDetailsView

import (
	"infinite-mitm/configs"
	events "infinite-mitm/internal/application/events"
	ui "infinite-mitm/internal/application/services/ui"
	helpers "infinite-mitm/internal/application/services/ui/network/helpers"
	errors "infinite-mitm/pkg/modules/errors"
	utilities "infinite-mitm/pkg/modules/utilities"
	"log"
	"os"
	"path"
	"sort"
	"strings"
	"time"

	"github.com/charmbracelet/bubbles/viewport"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type activeViewportType string
type updateRequest events.ProxyRequestEventData
type updateResponse events.ProxyResponseEventData

type TrafficModel struct {
	activeViewport  activeViewportType
	headersViewport viewport.Model
	bodyViewport    viewport.Model
	req             events.ProxyRequestEventData
	resp            events.ProxyResponseEventData

	maxWidth int
}

const (
	HeadersViewport activeViewportType = "headers"
	BodyViewport    activeViewportType = "body"
)

var TrafficDetailsModel *TrafficModel
var requestDetailsProgram *tea.Program

func NewTrafficDetailsModel(height int) *TrafficModel {
	maxWidth := utilities.GetTerminalWidth() - 30
	hvp := viewport.New(maxWidth, height)
	bvp := viewport.New(maxWidth, height)

	hvp.SetContent("...")
	bvp.SetContent("...")

	m := &TrafficModel{
		activeViewport: "headers",
		headersViewport: hvp,
		bodyViewport:    bvp,
		req:             events.ProxyRequestEventData{},
		resp:            events.ProxyResponseEventData{},
		maxWidth:        maxWidth,
	}

	TrafficDetailsModel = m
	return TrafficDetailsModel
}

func DEBUG__RunTrafficDetails() {
	requestDetailsProgram = tea.NewProgram(NewTrafficDetailsModel(25), tea.WithAltScreen())

	go func() {
		time.Sleep(1 * time.Second)
		data, _ := os.ReadFile(path.Join(configs.GetConfig().Extra.ProjectDir, "traffic", "test.jpg"))
		requestDetailsProgram.Send(updateRequest{
			ID: "1234",
			URL: "https://example.com/foo",
			Method: "GET",
			Headers: map[string]string{
				"Accept":              "text/html,application/xhtml+xml,application/xml;q=0.9,image/webp,*/*;q=0.8",
				"Accept-Encoding":     "gzip, deflate, br",
				"Accept-Language":     "en-US,en;q=0.5",
				"Cache-Control":       "no-cache",
				"Connection":          "keep-alive",
				"Host":                "example.com",
				"Referer":             "https://www.example.com/",
				"User-Agent":          "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/58.0.3029.110 Safari/537.3",
				"X-Requested-With":    "XMLHttpRequest",
			},
			Body: data,
			Proxified: true,
		})
	}()

	if _, err := requestDetailsProgram.Run(); err != nil {
		log.Fatalln(errors.Create(errors.ErrFatalException, err.Error()))
	}
}

func (m *TrafficModel) SetRequest(req events.ProxyRequestEventData) {
	m.req = req
	m.SetContent(m.req.Headers, m.req.Body)
}

func (m *TrafficModel) SetResponse(resp events.ProxyResponseEventData) {
	m.resp = resp
	m.SetContent(m.resp.Headers, m.resp.Body)
}

func (m *TrafficModel) SetContent(headers map[string]string, body []byte) {
	var headersString []string
	for key, value := range headers {
		headersString = append(headersString, lipgloss.NewStyle().Bold(true).Render(key) + ": " + utilities.WrapText(value, m.maxWidth))
	}

	sort.Strings(headersString)
	m.headersViewport.SetContent(strings.Join(headersString, "\n"))
	m.bodyViewport.SetContent(helpers.FormatHexView(body, m.maxWidth))
}

func (m *TrafficModel) Init() tea.Cmd {
	return nil
}

func (m *TrafficModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.maxWidth = utilities.GetTerminalWidth() - 20
		if m.activeViewport == "headers" && !utilities.IsEmpty(m.req) {
			m.headersViewport.Width = m.maxWidth
			m.SetContent(m.req.Headers, m.req.Body)
		} else if m.activeViewport == "body" && !utilities.IsEmpty(m.resp) {
			m.bodyViewport.Width = m.maxWidth
			m.SetContent(m.resp.Headers, m.resp.Body)
		}
	case updateRequest:
		m.SetRequest(events.ProxyRequestEventData(msg))
	case updateResponse:
		m.SetResponse(events.ProxyResponseEventData(msg))
	case tea.KeyMsg:
		switch msg.String() {
		case "q", "ctrl+c":
			return m, tea.Quit
		case "tab":
			if m.activeViewport == "headers" {
				m.activeViewport = "body"
			} else {
				m.activeViewport = "headers"
			}
		}

	}

	var cmds []tea.Cmd

	m.headersViewport, cmd = m.headersViewport.Update(msg)
	cmds = append(cmds, cmd)

	m.bodyViewport, cmd = m.bodyViewport.Update(msg)
	cmds = append(cmds, cmd)

	return m, tea.Batch(cmds...)
}

func (m *TrafficModel) View() string {
	var content string
	var sectionTitle string

	sectionHint := "(use mouse whell or arrows keys to scroll)"

	if m.activeViewport == HeadersViewport {
		content = m.headersViewport.View()
		sectionTitle = "Headers"
	} else if m.activeViewport == BodyViewport {
		content = m.bodyViewport.View()
		sectionTitle = "Body"
	}

	contentStyle := lipgloss.NewStyle().Padding(1, 2)
	sectionStyle := lipgloss.NewStyle().Bold(true)
	sectionHintStyle := lipgloss.NewStyle().Bold(false).Foreground(ui.Grey)

	if content == "" {
		content = "..."
	}

	return lipgloss.NewStyle().Render(lipgloss.JoinVertical(
		lipgloss.Top,
		sectionStyle.Render(sectionTitle + " " + sectionHintStyle.Render(sectionHint)),
		contentStyle.Render(content),
	))
}
