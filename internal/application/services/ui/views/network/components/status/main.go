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

package MITMApplicationUIServiceNetworkStatusBarComponent

import (
	"infinite-mitm/configs"
	"infinite-mitm/pkg/theme"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type ViewSize struct {
	Width  int
	Height int
}

type StatusBarModel struct {
	Size ViewSize

	width int

	name ItemProperties
	info ItemProperties
}

type StatusBarInfoUpdate struct {
	Message string
}

type ItemProperties struct {
	content string
	style   lipgloss.Style
}

var (
	statusBarStyle = lipgloss.NewStyle().
		Foreground(lipgloss.AdaptiveColor{Light: "#343433", Dark: "#C1C6B2"}).
		Background(lipgloss.AdaptiveColor{Light: "#D9DCCF", Dark: "#353533"}).
		Margin(1, 2)

	statusStyle = lipgloss.NewStyle().
		Inherit(statusBarStyle).
		Foreground(lipgloss.Color("#FFFDF5")).
		Background(lipgloss.Color(theme.ColorSunsetOrange.Light)).
		Padding(0, 1)

	statusInfoStyle = lipgloss.NewStyle().Inherit(statusBarStyle).MarginLeft(2)
)

func NewStatusBarModel(width int) StatusBarModel {
	m := StatusBarModel{
		name: ItemProperties{
			content: configs.GetConfig().Name,
			style: statusStyle,
		},
		info: ItemProperties{
			content: "Loading...",
			style: statusInfoStyle,
		},
	}

	vwidth, vheight := lipgloss.Size(m.View())
	m.Size.Width = vwidth
	m.Size.Height = vheight

	return m
}

func (m *StatusBarModel) SetWidth(width int) {
	m.width = width
}

func (m *StatusBarModel) SetInfoContent(text string) {
	m.info.content = text
}

func (m StatusBarModel) Update(msg tea.Msg) (StatusBarModel, tea.Cmd) {
	var cmd tea.Cmd
	switch msg := msg.(type) {
	case StatusBarInfoUpdate:
		m.SetInfoContent(msg.Message)
	}

	return m, cmd
}

func (m *StatusBarModel) View() string {
	bar := lipgloss.JoinHorizontal(lipgloss.Top,
		m.name.style.Render(m.name.content),
		m.info.style.Render(m.info.content),
	)

	return statusBarStyle.
		Width(m.width).
		MaxWidth(m.width).
		Margin(1, 2).
		Render(bar)
}
