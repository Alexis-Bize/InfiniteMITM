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
	"infinite-mitm/configs"
	ui "infinite-mitm/internal/application/services/ui"
	"strings"

	"github.com/charmbracelet/lipgloss"
)

type StatusBarModel struct {
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
		Background(lipgloss.AdaptiveColor{Light: "#D9DCCF", Dark: "#353533"})

	statusStyle = lipgloss.NewStyle().
		Inherit(statusBarStyle).
		Foreground(lipgloss.Color("#FFFDF5")).
		Background(lipgloss.Color(ui.SunsetOrange.Light)).
		Padding(0, 1).
		MarginRight(1)

	statusInfoStyle = lipgloss.NewStyle().Inherit(statusBarStyle)
)

func createStatusBar() StatusBarModel {
	m := StatusBarModel{
		name: ItemProperties{
			content: configs.GetConfig().Name,
			style: statusStyle,
		},
		info: ItemProperties{
			content: "",
			style: statusInfoStyle,
		},
	}

	return m
}

func getStatusBarWidth() int {
	w := GetTerminalWidth() - statusBarStyle.GetHorizontalPadding()
	return w
}

func updateStatusBarInfo(message string) {
	program.Send(StatusBarInfoUpdate(StatusBarInfoUpdate{
		Message: message,
	}))
}

func (m *StatusBarModel) View() string {
	doc := strings.Builder{}

	bar := lipgloss.JoinHorizontal(lipgloss.Top,
		m.name.style.Render(m.name.content),
		m.info.style.Render(m.info.content),
	)

	doc.WriteString(statusBarStyle.Width(getStatusBarWidth()).Render(bar))
	return doc.String()
}
