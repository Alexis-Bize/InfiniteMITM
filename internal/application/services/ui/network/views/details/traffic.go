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
	"encoding/hex"
	"fmt"
	"infinite-mitm/configs"
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

	"github.com/atotto/clipboard"
	"github.com/charmbracelet/bubbles/viewport"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/ncruces/zenity"
)

type activeViewportType string
type TrafficData struct {
	ID        string
	URL       string
	Status    int
	Method    string
	Headers   map[string]string
	Body      []byte
	Proxified  bool
}

type TrafficModel struct {
	activeViewport  activeViewportType
	headersViewport viewport.Model
	bodyViewport    viewport.Model
	data            TrafficData
	width int
}

const (
	HeadersViewportKey activeViewportType = "headers"
	BodyViewportKey    activeViewportType = "body"
)

var (
	CopyCommand = "ctrl+p"
	SaveCommand = "ctrl+o"
)

var TrafficDetailsModel *TrafficModel
var requestDetailsProgram *tea.Program

func NewTrafficDetailsModel(width int, height int) *TrafficModel {
	hvp := viewport.New(width, height)
	bvp := viewport.New(width, height)

	hvp.SetContent("...")
	bvp.SetContent("...")

	m := &TrafficModel{
		activeViewport:  HeadersViewportKey,
		headersViewport: hvp,
		bodyViewport:    bvp,
		data:            TrafficData{},
		width:           width,
	}

	TrafficDetailsModel = m
	return TrafficDetailsModel
}

func DEBUG__RunTrafficDetails() {
	requestDetailsProgram = tea.NewProgram(
		NewTrafficDetailsModel(utilities.GetTerminalWidth() - 20, 10),
		tea.WithAltScreen(),
	)

	go func() {
		time.Sleep(1 * time.Second)
		data, _ := os.ReadFile(path.Join(configs.GetConfig().Extra.ProjectDir, "traffic", "test.jpg"))
		requestDetailsProgram.Send(TrafficData{
			ID: "1234",
			URL: "https://example.com/foo",
			Method: "GET",
			Headers: map[string]string{
				"Accept":              "text/html,application/xhtml+xml,application/xml;q=0.9,image/webp,*/*;q=0.8",
				"Accept-Encoding":     "gzip, deflate, br",
				"Accept-Language":     "en-US,en;q=0.5",
				"Cache-Control":       "no-cache",
				"Content-Type":        "image/jpeg",
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

func (m *TrafficModel) Init() tea.Cmd {
	return nil
}

func (m *TrafficModel) SetTrafficData(data TrafficData) {
	m.data = data
	m.SetContent(m.data.Headers, m.data.Body)
}

func (m *TrafficModel) SetContent(headers map[string]string, body []byte) {
	var headersString []string
	for key, value := range headers {
		headersString = append(
			headersString,
			lipgloss.NewStyle().Bold(true).Render(key) + ": " + utilities.WrapText(value, m.width),
		)
	}

	sort.Strings(headersString)
	m.headersViewport.SetContent(strings.Join(headersString, "\n"))
	m.bodyViewport.SetContent(helpers.FormatHexView(body, m.width))
}

func (m *TrafficModel) CopyToClipboard() {
	if utilities.IsEmpty(m.data) {
		return
	}

	if m.activeViewport == HeadersViewportKey {
		var headersString []string
		for key, value := range m.data.Headers {
			headersString = append(
				headersString,
				key + ": " + utilities.WrapText(value, m.width),
			)
		}

		sort.Strings(headersString)
		clipboard.WriteAll(strings.Join(headersString, "\n"))
	} else if m.activeViewport == BodyViewportKey {
		clipboard.WriteAll(hex.EncodeToString(m.data.Body))
	}
}

func (m *TrafficModel) SaveToDrive() {
	if utilities.IsEmpty(m.data) {
		return
	}

	if m.activeViewport == BodyViewportKey {
		var extension string

		switch m.data.Headers["Content-Type"] {
		case "application/json":
				extension = "json"
		case "image/jpg":
		case "image/jpeg":
				extension = "jpg"
		case "image/png":
				extension = "png"
		default:
				extension = "bin"
		}

		filename := fmt.Sprintf("body.%s", extension)
		filePath, err := zenity.SelectFileSave(
			zenity.Title("Save File"),
			zenity.Filename(path.Join(configs.GetConfig().Extra.ProjectDir, "traffic", filename)),
			zenity.ConfirmOverwrite(),
		)

		if err == nil {
			os.WriteFile(filePath, m.data.Body, 0644)
		}
	}
}

func (m *TrafficModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		utilities.ClearTerminal()

		m.width = msg.Width
		m.headersViewport.Width = m.width
		m.bodyViewport.Width = m.width

		if !utilities.IsEmpty(m.data) {
			m.SetContent(m.data.Headers, m.data.Body)
		}
	case TrafficData:
		m.SetTrafficData(TrafficData(msg))
	case tea.KeyMsg:
		switch msg.String() {
		case "q", "ctlr+c":
			return m, tea.Quit
		case CopyCommand:
			m.CopyToClipboard()
		case SaveCommand:
			m.SaveToDrive()
		case "tab":
			if m.activeViewport == HeadersViewportKey {
				m.activeViewport = BodyViewportKey
			} else {
				m.activeViewport = HeadersViewportKey
			}
		}
	}

	cmds := []tea.Cmd{cmd}

	if m.activeViewport == HeadersViewportKey {
		m.headersViewport, cmd = m.headersViewport.Update(msg)
		cmds = append(cmds, cmd)
	}

	if m.activeViewport == BodyViewportKey {
		m.bodyViewport, cmd = m.bodyViewport.Update(msg)
		cmds = append(cmds, cmd)
	}

	return m, tea.Batch(cmds...)
}

func (m *TrafficModel) View() string {
	var content string
	var sectionTitle string
	var viewportActions string

	viewportActionsList := []string{"Loading..."}
	sectionHint := "(use mouse whell or arrows keys to scroll)"

	if m.activeViewport == HeadersViewportKey {
		content = m.headersViewport.View()
		sectionTitle = "Headers"
	} else if m.activeViewport == BodyViewportKey {
		content = m.bodyViewport.View()
		sectionTitle = "Body"
	}

	if !utilities.IsEmpty(m.data) {
		viewportActionsList = []string{}

		if m.activeViewport == HeadersViewportKey {
			viewportActionsList = append(viewportActionsList, fmt.Sprintf("Copy headers to clipboard (%s)", CopyCommand))
		} else if m.activeViewport == BodyViewportKey {
			viewportActionsList = append(viewportActionsList, fmt.Sprintf("Copy hex to clipboard (%s)", CopyCommand))
			viewportActionsList = append(viewportActionsList, fmt.Sprintf("Save (%s)", SaveCommand))
		}
	}

	contentStyle := lipgloss.NewStyle().Padding(1, 2)
	sectionStyle := lipgloss.NewStyle().Bold(true)
	sectionHintStyle := lipgloss.NewStyle().Bold(false).Foreground(ui.Grey)
	viewportActionsStyle := lipgloss.NewStyle().Padding(0, 1).MarginRight(1).Foreground(ui.Light).Background(ui.Grey)

	if content == "" {
		content = "..."
	}

	for _, k := range viewportActionsList {
		viewportActions += viewportActionsStyle.Render(k)
	}

	return lipgloss.NewStyle().MaxWidth(m.width).Render(lipgloss.JoinVertical(
		lipgloss.Top,
		sectionStyle.Render(sectionTitle + " " + sectionHintStyle.Render(sectionHint)),
		contentStyle.Render(content),
		viewportActions,
	))
}
