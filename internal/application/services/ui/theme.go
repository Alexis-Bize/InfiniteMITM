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

package MITMApplicationUIService

import (
	"github.com/charmbracelet/huh"
	"github.com/charmbracelet/lipgloss"
)

var (
	NormalFg     = lipgloss.AdaptiveColor{Light: "233", Dark: "255"}
	SunsetOrange = lipgloss.AdaptiveColor{Light: "#FF5F5D", Dark: "#FF6F6D"}
	SunsetYellow = lipgloss.AdaptiveColor{Light: "#FFDD57", Dark: "#FFE066"}
	TwilightBlue = lipgloss.AdaptiveColor{Light: "#1C2541", Dark: "#3A506B"}
	Light = lipgloss.AdaptiveColor{Light: "#f5f5f5", Dark: "#dedede"}
	Dark = lipgloss.AdaptiveColor{Light: "#020202", Dark: "#212121"}
	Grey = lipgloss.AdaptiveColor{Light: "#3c3c3c", Dark: "#6c6c6c"}

	Success      = lipgloss.AdaptiveColor{Light: "#43BF6D", Dark: "#73F59F"}
	Error        = lipgloss.AdaptiveColor{Light: "#bf4343", Dark: "#f57373"}
	Warn         = lipgloss.AdaptiveColor{Light: "#bf9543", Dark: "#f5be73"}
)

func ThemeMITM() *huh.Theme {
	t := huh.ThemeBase()

	t.Focused.Base = t.Focused.Base.BorderForeground(lipgloss.Color("238"))
	t.Focused.Title = t.Focused.Title.Foreground(SunsetOrange).Bold(true)
	t.Focused.NoteTitle = t.Focused.NoteTitle.Foreground(SunsetOrange).Bold(true).MarginBottom(1)
	t.Focused.Directory = t.Focused.Directory.Foreground(TwilightBlue)
	t.Focused.Description = t.Focused.Description.Foreground(lipgloss.AdaptiveColor{Light: "236", Dark: "243"})
	t.Focused.ErrorIndicator = t.Focused.ErrorIndicator.Foreground(Error)
	t.Focused.ErrorMessage = t.Focused.ErrorMessage.Foreground(Error)
	t.Focused.SelectSelector = t.Focused.SelectSelector.Foreground(SunsetYellow)
	t.Focused.NextIndicator = t.Focused.NextIndicator.Foreground(SunsetYellow)
	t.Focused.PrevIndicator = t.Focused.PrevIndicator.Foreground(SunsetYellow)
	t.Focused.Option = t.Focused.Option.Foreground(NormalFg)
	t.Focused.MultiSelectSelector = t.Focused.MultiSelectSelector.Foreground(SunsetYellow)
	t.Focused.SelectedOption = t.Focused.SelectedOption.Foreground(Success)
	t.Focused.SelectedPrefix = lipgloss.NewStyle().Foreground(lipgloss.AdaptiveColor{Light: "#02CF92", Dark: "#02A877"}).SetString("✓ ")
	t.Focused.UnselectedPrefix = lipgloss.NewStyle().Foreground(lipgloss.AdaptiveColor{Light: "236", Dark: "243"}).SetString("• ")
	t.Focused.UnselectedOption = t.Focused.UnselectedOption.Foreground(NormalFg)
	t.Focused.FocusedButton = t.Focused.FocusedButton.Foreground(Warn).Background(SunsetOrange)
	t.Focused.Next = t.Focused.FocusedButton
	t.Focused.BlurredButton = t.Focused.BlurredButton.Foreground(NormalFg).Background(lipgloss.AdaptiveColor{Light: "252", Dark: "237"})

	t.Focused.TextInput.Cursor = t.Focused.TextInput.Cursor.Foreground(Success)
	t.Focused.TextInput.Placeholder = t.Focused.TextInput.Placeholder.Foreground(lipgloss.AdaptiveColor{Light: "248", Dark: "238"})
	t.Focused.TextInput.Prompt = t.Focused.TextInput.Prompt.Foreground(SunsetYellow)

	t.Blurred = t.Focused
	t.Blurred.Base = t.Focused.Base.BorderStyle(lipgloss.HiddenBorder())
	t.Blurred.NextIndicator = lipgloss.NewStyle()
	t.Blurred.PrevIndicator = lipgloss.NewStyle()

	return t
}
