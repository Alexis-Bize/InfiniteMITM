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

package MITMApplicationUIServiceNetworkTrafficDetailsComponent

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

type activeViewType string
type TrafficData struct {
	Headers   map[string]string
	Body      []byte
}

type TrafficModel struct {
	width int

	headersModel viewport.Model
	bodyModel    viewport.Model

	data        TrafficData
	activeView  activeViewType
	copyPressed bool
	focused     bool
}

const (
	HeadersViewKey activeViewType = "headers"
	BodyViewKey    activeViewType = "body"
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
		headersModel: hvp,
		bodyModel:    bvp,
		activeView:   HeadersViewKey,
		data:         TrafficData{},
		width:        width,
		copyPressed:  false,
		focused:      false,
	}

	return m
}

func (m *TrafficModel) Focus() {
	m.focused = true
	m.SetCopyPress(false)
}

func (m *TrafficModel) Blur() {
	m.focused = false
	m.SetCopyPress(false)
	m.headersModel.SetYOffset(0)
	m.bodyModel.SetYOffset(0)
}

func (m *TrafficModel) SetCopyPress(pressed bool) {
	m.copyPressed = pressed
}

func (m *TrafficModel) SetWidth(width int) {
	m.width = width
	m.headersModel.Width = width - 20
	m.bodyModel.Width = width - 20
}

func (m *TrafficModel) SetActiveView(key activeViewType) {
	m.activeView = key
	m.headersModel.SetYOffset(0)
	m.bodyModel.SetYOffset(0)
	m.SetCopyPress(false)
}

func (m *TrafficModel) SwitchActiveView() {
	if m.activeView == HeadersViewKey {
		m.SetActiveView(BodyViewKey)
	} else if m.activeView == BodyViewKey {
		m.SetActiveView(HeadersViewKey)
	}
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
			lipgloss.NewStyle().Bold(true).Render(key) + ": " + utilities.WrapText(value, m.headersModel.Width),
		)
	}

	sort.Strings(headersString)
	m.headersModel.SetContent(strings.Join(headersString, "\n"))
	m.bodyModel.SetContent(helpers.FormatHexView(body, m.bodyModel.Width))
}

func (m *TrafficModel) CopyToClipboard() {
	if len(m.data.Headers) == 0 {
		return
	}

	m.SetCopyPress(true)

	if m.activeView == HeadersViewKey {
		var headersString []string
		for key, value := range m.data.Headers {
			headersString = append(headersString, key + ":" + value)
		}

		sort.Strings(headersString)
		helpers.CopyToClipboard(strings.Join(headersString, "\n"))
	} else if m.activeView == BodyViewKey {
		helpers.CopyToClipboard(hex.EncodeToString(m.data.Body))
	}
}

func (m TrafficModel) SaveToDisk() {
	if len(m.data.Body) != 0 && m.activeView == BodyViewKey {
		helpers.SaveToDisk(m.data.Body, m.data.Headers["Content-Type"])
	}
}

func (m TrafficModel) Update(msg tea.Msg) (TrafficModel, tea.Cmd) {
	var cmd tea.Cmd
	cmds := []tea.Cmd{}

	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
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
			case "enter":
				m.SwitchActiveView()
			}
		}
	}

	if m.activeView == HeadersViewKey {
		m.headersModel, cmd = m.headersModel.Update(msg)
		cmds = append(cmds, cmd)
	}

	if m.activeView == BodyViewKey {
		m.bodyModel, cmd = m.bodyModel.Update(msg)
		cmds = append(cmds, cmd)
	}

	return m, tea.Batch(cmds...)
}

func (m TrafficModel) View() string {
	var sectionTitle string
	var viewportActions string

	if m.activeView == HeadersViewKey {
		sectionTitle = "Headers"
	} else if m.activeView == BodyViewKey {
		sectionTitle = "Body"
	}

	baseContentStyle := lipgloss.NewStyle().Padding(1, 2)
	contentStyle := baseContentStyle
	emptyContentStyle := lipgloss.NewStyle().Padding(1, 2).Foreground(theme.ColorGrey).Italic(true)
	viewportActionsStyle := lipgloss.NewStyle().Padding(0, 1).MarginRight(1).Foreground(theme.ColorLight).Background(theme.ColorGrey)

	viewportActionsList := []string{"q: Go back"}
	content := "Waiting..."

	hasData := !utilities.IsEmpty(m.data)
	if !hasData {
		contentStyle = emptyContentStyle
	}

	if hasData {
		copiedText := "✓ Copied"

		if m.activeView == HeadersViewKey {
			headersLength := len(m.data.Headers)
			if headersLength == 0 {
				content = "Empty headers"
				contentStyle = emptyContentStyle
			}

			if headersLength != 0 {
				content = m.headersModel.View()
				if m.copyPressed {
					viewportActionsList = append(viewportActionsList, copiedText)
				} else {
					viewportActionsList = append(viewportActionsList, fmt.Sprintf("%s: Copy headers to clipboard", CopyCommand))
				}
			}
		} else if m.activeView == BodyViewKey {
			bodyLength := len(m.data.Body)
			if bodyLength == 0 {
				content = "Empty body"
				contentStyle = emptyContentStyle
			}

			if bodyLength != 0 {
				content = m.bodyModel.View()
				if m.copyPressed {
					viewportActionsList = append(viewportActionsList, copiedText)
				} else {
					viewportActionsList = append(viewportActionsList, fmt.Sprintf("%s: Copy hex to clipboard", CopyCommand))
				}

				viewportActionsList = append(viewportActionsList, fmt.Sprintf("%s: Save", SaveCommand))
			}
		}
	}

	for _, k := range viewportActionsList {
		viewportActions += viewportActionsStyle.Render(k)
	}

	return lipgloss.NewStyle().MaxWidth(m.width).Render(lipgloss.JoinVertical(
		lipgloss.Top,
		lipgloss.JoinHorizontal(
			lipgloss.Top,
			lipgloss.NewStyle().Bold(true).Render(sectionTitle),
			lipgloss.NewStyle().MarginLeft(2).Foreground(theme.ColorGrey).Render("Enter ↵: Switch between headers and body"),
		),
		contentStyle.Render(content),
		lipgloss.NewStyle().Foreground(theme.ColorGrey).MarginBottom(1).Render("↑/↓: Scroll"),
		viewportActions,
	))
}
