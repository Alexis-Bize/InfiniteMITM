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

package MITMApplicationUIServiceNetworkPartialTrafficDetails

import (
	"encoding/hex"
	"fmt"
	helpers "infinite-mitm/internal/application/services/ui/helpers"
	theme "infinite-mitm/internal/application/services/ui/theme"
	utilities "infinite-mitm/pkg/modules/utilities"
	"sort"
	"strings"

	"github.com/charmbracelet/bubbles/viewport"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type activeViewportType string
type TrafficData struct {
	Headers   map[string]string
	Body      []byte
}

type TrafficModel struct {
	Width int

	HeadersViewportModel viewport.Model
	BodyViewportModel    viewport.Model

	data            TrafficData
	activeViewport  activeViewportType
	copyPressed     bool
	focused bool
}

const (
	HeadersViewport activeViewportType = "headers"
	BodyViewport    activeViewportType = "body"
)

var (
	CopyCommand = "ctrl+p"
	SaveCommand = "ctrl+o"
)

func NewTrafficDetailsModel(width int, height int) TrafficModel {
	hvp := viewport.New(width, height)
	bvp := viewport.New(width, height)

	hvp.SetContent("...")
	bvp.SetContent("...")

	m := TrafficModel{
		HeadersViewportModel: hvp,
		BodyViewportModel:    bvp,
		activeViewport:  HeadersViewport,
		data:            TrafficData{},
		Width:           width,
		copyPressed:     false,
	}

	return m
}

func (m *TrafficModel) Focus() {
	m.focused = true
	m.activeViewport = HeadersViewport
}

func (m *TrafficModel) Blur() {
	m.focused = false
	m.activeViewport = HeadersViewport
}

func (m TrafficModel) Init() tea.Cmd {
	return nil
}

func (m *TrafficModel) SetWidth(width int) {
	m.Width = width
}

func (m *TrafficModel) SetTrafficData(data TrafficData) {
	m.data = TrafficData(data)
	m.SetContent(m.data.Headers, m.data.Body)
}

func (m *TrafficModel) SetContent(headers map[string]string, body []byte) {
	var headersString []string
	for key, value := range headers {
		headersString = append(
			headersString,
			lipgloss.NewStyle().Bold(true).Render(key) + ": " + utilities.WrapText(value, m.Width),
		)
	}

	sort.Strings(headersString)
	m.HeadersViewportModel.SetContent(strings.Join(headersString, "\n"))
	m.BodyViewportModel.SetContent(helpers.FormatHexView(body, m.Width))
}

func (m *TrafficModel) CopyToClipboard() {
	if utilities.IsEmpty(m.data) {
		return
	}

	m.copyPressed = true

	if m.activeViewport == HeadersViewport {
		var headersString []string
		for key, value := range m.data.Headers {
			headersString = append(headersString, key + ": " + utilities.WrapText(value, m.Width))
		}

		sort.Strings(headersString)
		helpers.CopyToClipboard(strings.Join(headersString, "\n"))
	} else if m.activeViewport == BodyViewport {
		helpers.CopyToClipboard(hex.EncodeToString(m.data.Body))
	}
}

func (m *TrafficModel) SaveToDisk() {
	if !utilities.IsEmpty(m.data) && m.activeViewport == BodyViewport {
		helpers.SaveToDisk(m.data.Body, m.data.Headers["Content-Type"])
	}
}

func (m *TrafficModel) ResetCopyPress() {
	m.copyPressed = false
}

func (m TrafficModel) Update(msg tea.Msg) (TrafficModel, tea.Cmd) {
	var cmd tea.Cmd
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.SetWidth(msg.Width)
		if !utilities.IsEmpty(m.data) {
			m.SetContent(m.data.Headers, m.data.Body)
		}
	case TrafficData:
		m.SetTrafficData(TrafficData(msg))
	case tea.KeyMsg:
		if m.focused {
			switch msg.String() {
			case CopyCommand:
				m.CopyToClipboard()
			case SaveCommand:
				m.SaveToDisk()
			case "left", "right":
				m.ResetCopyPress()
				if m.activeViewport == HeadersViewport {
					m.activeViewport = BodyViewport
				} else {
					m.activeViewport = HeadersViewport
				}
			}
		}
	}

	cmds := []tea.Cmd{cmd}

	if m.activeViewport == HeadersViewport {
		m.HeadersViewportModel, cmd = m.HeadersViewportModel.Update(msg)
		cmds = append(cmds, cmd)
	}

	if m.activeViewport == BodyViewport {
		m.BodyViewportModel, cmd = m.BodyViewportModel.Update(msg)
		cmds = append(cmds, cmd)
	}

	return m, tea.Batch(cmds...)
}

func (m TrafficModel) View() string {
	var sectionTitle string
	var viewportActions string

	contentStyle := lipgloss.NewStyle().Padding(1, 2)
	emptyContentStyle := lipgloss.NewStyle().Padding(1, 2).Foreground(theme.ColorGrey).Italic(true)

	sectionStyle := lipgloss.NewStyle().Bold(true)
	sectionHintStyle := lipgloss.NewStyle().Foreground(theme.ColorGrey).MarginLeft(1)
	sectionHint := "(use the mouse wheel or arrow keys ⭥ to scroll)"

	switchHintStyle := lipgloss.NewStyle()
	switchHint := "Use arrow keys ↔ to switch between headers and body"

	viewportActionsStyle := lipgloss.NewStyle().Padding(0, 1).MarginRight(1).Foreground(theme.ColorLight).Background(theme.ColorGrey)
	viewportActionsList := []string{"Waiting..."}

	copiedText := "✓ Copied"
	content := "..."

	if m.activeViewport == HeadersViewport {
		content = m.HeadersViewportModel.View()
		sectionTitle = "Headers"
	} else if m.activeViewport == BodyViewport {
		content = m.BodyViewportModel.View()
		sectionTitle = "Body"
	}

	if !utilities.IsEmpty(m.data) {
		viewportActionsList = []string{}

		switch m.activeViewport {
		case HeadersViewport:
			if len(m.data.Headers) != 0 {
				if m.copyPressed {
					viewportActionsList = append(viewportActionsList, copiedText)
				} else {
					viewportActionsList = append(viewportActionsList, fmt.Sprintf("Copy headers to clipboard (%s)", CopyCommand))
				}
			} else {
				content = "Empty headers"
				contentStyle = emptyContentStyle
			}
		case BodyViewport:
			if len(m.data.Body) != 0 {
				if m.copyPressed {
					viewportActionsList = append(viewportActionsList, copiedText)
				} else {
					viewportActionsList = append(viewportActionsList, fmt.Sprintf("Copy hex to clipboard (%s)", CopyCommand))
				}

				viewportActionsList = append(viewportActionsList, fmt.Sprintf("Save (%s)", SaveCommand))
			} else {
				content = "Empty body"
				contentStyle = emptyContentStyle
			}
		}
	}

	for _, k := range viewportActionsList {
		viewportActions += viewportActionsStyle.Render(k)
	}

	return lipgloss.NewStyle().MaxWidth(m.Width).Render(lipgloss.JoinVertical(
		lipgloss.Top,
		lipgloss.JoinHorizontal(
			lipgloss.Top,
			sectionStyle.Render(sectionTitle),
			sectionHintStyle.Render(sectionHint),
		),
		switchHintStyle.Render(switchHint),
		contentStyle.Render(content),
		viewportActions,
	))
}
