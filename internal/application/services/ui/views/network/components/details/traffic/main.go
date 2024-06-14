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
	"infinite-mitm/pkg/request"
	"infinite-mitm/pkg/utilities"
	"net/url"
	"path"
	"sort"
	"strings"

	"github.com/charmbracelet/bubbles/viewport"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type activeViewType string
type TrafficData struct {
	Dummy     bool
	Headers   map[string]string
	Body      []byte
	URL       string
}

type TrafficModel struct {
	width int

	headersModel viewport.Model
	bodyModel    viewport.Model

	data        *TrafficData
	activeView  activeViewType
	copyPressed bool
	focused     bool
}

const BodyMaxViewLength = 200 * 1024

const (
	HeadersViewKey activeViewType = "headers"
	BodyViewKey    activeViewType = "body"
)

const (
	CopyCommand = "ctrl+p"
	SaveCommand = "ctrl+o"
)

var (
	emptyHeadersString = "Empty headers"
	emptyBodyString    = "Empty body"
	dummyString        = "Not available (skipped)"
	waitingString      = "Waiting..."
	copiedString       = "✓ Copied"
	bodyTooLargeString = "Content too large to be displayed"

	switchHintString   = "Enter ↵: Switch between headers and body"
	scrollHintString   = "↑/↓: Scroll"

	copyHeadersString  = fmt.Sprintf("%s: Copy headers to clipboard", CopyCommand)
	copyHexString      = fmt.Sprintf("%s: Copy hex to clipboard", CopyCommand)
	saveBodyString     = fmt.Sprintf("%s: Save", SaveCommand)
)

var (
	baseContentStyle = lipgloss.NewStyle().Padding(1, 2)
	contentStyle = baseContentStyle

	emptyContentStyle = lipgloss.NewStyle().
		Padding(1, 2).
		Foreground(theme.ColorGrey).
		Italic(true)

	viewportActionsStyle = lipgloss.NewStyle().
		Padding(0, 1).MarginRight(1).
		Foreground(theme.ColorLight).
		Background(theme.ColorGrey)
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
		data:         &TrafficData{},
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

	if m.activeView == HeadersViewKey {
		m.headersModel.SetYOffset(0)
	} else if m.activeView == BodyViewKey {
		m.bodyModel.SetYOffset(0)
	}

	m.SetContent(map[string]string{}, nil)
}

func (m *TrafficModel) SetCopyPress(pressed bool) {
	m.copyPressed = pressed
}

func (m *TrafficModel) SetWidth(width int) {
	m.width = width
	m.headersModel.Width = width - 20
	m.bodyModel.Width = width - 20
}

func (m *TrafficModel) SwitchActiveView() {
	if m.activeView == HeadersViewKey {
		m.SetActiveView(BodyViewKey)
	} else if m.activeView == BodyViewKey {
		m.SetActiveView(HeadersViewKey)
	}
}

func (m *TrafficModel) SetActiveView(key activeViewType) {
	m.activeView = key
	m.SetCopyPress(false)

	if key == HeadersViewKey {
		m.SetContent(m.data.Headers, m.data.Body)
		m.Focus()
	} else if key == BodyViewKey {
		m.SetContent(m.data.Headers, m.data.Body)
		m.Focus()
	}
}

func (m *TrafficModel) SetTrafficData(data *TrafficData) {
	if data.Headers[request.ContentTypeHeaderKey] == "" {
		delete(data.Headers, request.ContentTypeHeaderKey)
	}

	m.data = data

	if m.focused {
		m.SetContent(m.data.Headers, m.data.Body)
	}
}

func (m *TrafficModel) SetContent(headers map[string]string, body []byte) {
	var headersString []string
	for key, value := range headers {
		if strings.TrimSpace(value) != "" {
			headersString = append(
				headersString,
				lipgloss.NewStyle().Bold(true).Render(key) + ": " + utilities.WrapText(value, m.headersModel.Width),
			)
		}
	}

	if m.activeView == HeadersViewKey {
		if len(headersString) != 0 {
			sort.Strings(headersString)
			m.headersModel.SetContent(strings.Join(headersString, "\n"))
		} else {
			m.headersModel.SetContent("")
		}

		m.headersModel.GotoTop()
	} else if m.activeView == BodyViewKey {
		if len(body) <= BodyMaxViewLength {
			m.bodyModel.SetContent(helpers.FormatHexView(body, m.bodyModel.Width))
		} else {
			m.bodyModel.SetContent("")
		}

		m.bodyModel.GotoTop()
	}
}

func (m *TrafficModel) CopyToClipboard() {
	if len(m.data.Headers) == 0 {
		return
	}

	m.SetCopyPress(true)

	if m.activeView == HeadersViewKey {
		var headersString []string
		for key, value := range m.data.Headers {
			if strings.TrimSpace(value) != "" {
				headersString = append(headersString, key + ":" + value)
			}
		}

		sort.Strings(headersString)
		helpers.CopyToClipboard(strings.Join(headersString, "\n"))
	} else if m.activeView == BodyViewKey {
		helpers.CopyToClipboard(hex.EncodeToString(m.data.Body))
	}
}

func (m TrafficModel) SaveToDisk() {
	if len(m.data.Body) != 0 && m.activeView == BodyViewKey {
		ct := m.data.Headers[request.ContentTypeHeaderKey]
		// dirty hack for old blobs PNG assets
		if ct == "application/octet-stream" {
			if string(m.data.Body[1:4]) == "PNG" {
				ct = "image/png"
			}
		}

		filename := "body"
		parse, err := url.Parse(m.data.URL)
		if err == nil {
			filename = strings.Split(path.Base(parse.Path), ".")[0]
		}

		helpers.SaveToDisk(m.data.Body, filename, ct)
	}
}

func (m TrafficModel) Update(msg tea.Msg) (TrafficModel, tea.Cmd) {
	var cmd tea.Cmd
	cmds := []tea.Cmd{}

	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		if m.focused {
			m.SetContent(m.data.Headers, m.data.Body)
		}
	case TrafficData:
		m.SetTrafficData(&msg)
	case tea.KeyMsg:
		if !m.focused {
			return m, tea.Batch(cmds...)
		}

		switch msg.String() {
		case CopyCommand:
			m.CopyToClipboard()
			return m, tea.Batch(cmds...)
		case SaveCommand:
			m.SaveToDisk()
			return m, tea.Batch(cmds...)
		case "enter":
			m.SwitchActiveView()
			return m, tea.Batch(cmds...)
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

	viewportActionsList := []string{"q: Go back"}
	content := waitingString

	headersLength := len(m.data.Headers)
	bodyLength := len(m.data.Body)
	hasData := headersLength != 0 || bodyLength != 0

	if !hasData || m.data.Dummy {
		contentStyle = emptyContentStyle
		if m.data.Dummy {
			content = dummyString
		}
	} else {
		if m.copyPressed {
			viewportActionsList = append(viewportActionsList, copiedString)
		}

		switch m.activeView {
		case HeadersViewKey:
			if headersLength == 0 {
				content = emptyHeadersString
				contentStyle = emptyContentStyle
			} else {
				content = m.headersModel.View()
				if !m.copyPressed {
					viewportActionsList = append(viewportActionsList, copyHeadersString)
				}
			}
		case BodyViewKey:
			if bodyLength == 0 {
				content = emptyBodyString
				contentStyle = emptyContentStyle
			} else {
				if bodyLength < BodyMaxViewLength {
					content = m.bodyModel.View()
					if !m.copyPressed {
						viewportActionsList = append(viewportActionsList, copyHexString)
					}

				} else {
					content = bodyTooLargeString
					contentStyle = emptyContentStyle
				}

				viewportActionsList = append(viewportActionsList, saveBodyString)
			}
		}
	}

	for _, k := range viewportActionsList {
		viewportActions += viewportActionsStyle.Render(k)
	}

	return lipgloss.NewStyle().
		MaxWidth(m.width).
		Render(lipgloss.JoinVertical(
			lipgloss.Top,
			lipgloss.JoinHorizontal(
				lipgloss.Top,
				lipgloss.NewStyle().
					Bold(true).
					Render(sectionTitle),
				lipgloss.NewStyle().
					MarginLeft(2).
					Foreground(theme.ColorGrey).
					Render(switchHintString),
			),
			contentStyle.Render(content),
			lipgloss.NewStyle().
				Foreground(theme.ColorGrey).
				MarginBottom(1).
				Render(scrollHintString),
			viewportActions,
	))
}
