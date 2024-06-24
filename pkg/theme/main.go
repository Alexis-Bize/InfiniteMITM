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

package theme

import (
	"github.com/charmbracelet/huh"
	"github.com/charmbracelet/lipgloss"
)

var (
	ColorNormalFg     = lipgloss.AdaptiveColor{Light: "233", Dark: "255"}
	ColorSunsetOrange = lipgloss.AdaptiveColor{Light: "#FF5F5D", Dark: "#FF6F6D"}
	ColorSunsetYellow = lipgloss.AdaptiveColor{Light: "#FFDD57", Dark: "#FFE066"}
	ColorTwilightBlue = lipgloss.AdaptiveColor{Light: "#1C2541", Dark: "#3A506B"}
	ColorLight = lipgloss.AdaptiveColor{Light: "#f5f5f5", Dark: "#dedede"}
	ColorDark = lipgloss.AdaptiveColor{Light: "#020202", Dark: "#212121"}
	ColorGrey = lipgloss.AdaptiveColor{Light: "#3c3c3c", Dark: "#6c6c6c"}
	ColorLightYellow = lipgloss.Color("229")
	ColorNeonBlue = lipgloss.Color("57")

	ColorSuccess      = lipgloss.AdaptiveColor{Light: "#43BF6D", Dark: "#73F59F"}
	ColorError        = lipgloss.AdaptiveColor{Light: "#bf4343", Dark: "#f57373"}
	ColorWarn         = lipgloss.AdaptiveColor{Light: "#bf9543", Dark: "#f5be73"}
)

func ThemeMITM() *huh.Theme {
	t := huh.ThemeBase()

	t.Focused.Base = t.Focused.Base.Border(lipgloss.Border{}).Padding(1, 2)
	t.Focused.Title = t.Focused.Title.Foreground(ColorSunsetOrange).Bold(true).MarginBottom(1)
	t.Focused.NoteTitle = t.Focused.NoteTitle.Foreground(ColorSunsetOrange).Bold(true).MarginBottom(1)
	t.Focused.Directory = t.Focused.Directory.Foreground(ColorTwilightBlue)
	t.Focused.Description = t.Focused.Description.Foreground(lipgloss.AdaptiveColor{Light: "236", Dark: "243"})
	t.Focused.ErrorIndicator = t.Focused.ErrorIndicator.Foreground(ColorError)
	t.Focused.ErrorMessage = t.Focused.ErrorMessage.Foreground(ColorError)
	t.Focused.SelectSelector = t.Focused.SelectSelector.Foreground(ColorSunsetYellow)
	t.Focused.NextIndicator = t.Focused.NextIndicator.Foreground(ColorSunsetYellow)
	t.Focused.PrevIndicator = t.Focused.PrevIndicator.Foreground(ColorSunsetYellow)
	t.Focused.Option = t.Focused.Option.Foreground(ColorNormalFg)
	t.Focused.MultiSelectSelector = t.Focused.MultiSelectSelector.Foreground(ColorSunsetYellow)
	t.Focused.SelectedOption = t.Focused.SelectedOption.Foreground(ColorSuccess)
	t.Focused.SelectedPrefix = lipgloss.NewStyle().Foreground(lipgloss.AdaptiveColor{Light: "#02CF92", Dark: "#02A877"}).SetString("✓ ")
	t.Focused.UnselectedPrefix = lipgloss.NewStyle().Foreground(lipgloss.AdaptiveColor{Light: "236", Dark: "243"}).SetString("• ")
	t.Focused.UnselectedOption = t.Focused.UnselectedOption.Foreground(ColorNormalFg)
	t.Focused.FocusedButton = t.Focused.FocusedButton.Foreground(ColorWarn).Background(ColorSunsetOrange)
	t.Focused.Next = t.Focused.FocusedButton
	t.Focused.BlurredButton = t.Focused.BlurredButton.Foreground(ColorNormalFg).Background(lipgloss.AdaptiveColor{Light: "252", Dark: "237"})

	t.Focused.TextInput.Cursor = t.Focused.TextInput.Cursor.Foreground(ColorSuccess)
	t.Focused.TextInput.Placeholder = t.Focused.TextInput.Placeholder.Foreground(lipgloss.AdaptiveColor{Light: "248", Dark: "238"})
	t.Focused.TextInput.Prompt = t.Focused.TextInput.Prompt.Foreground(ColorSunsetYellow)

	t.Blurred = t.Focused
	t.Blurred.Base = t.Focused.Base.BorderStyle(lipgloss.HiddenBorder())
	t.Blurred.NextIndicator = lipgloss.NewStyle()
	t.Blurred.PrevIndicator = lipgloss.NewStyle()

	return t
}

func StatusCodeToColorStyle(statusCode int) lipgloss.Style {
	successColor := lipgloss.NewStyle().Foreground(ColorSuccess)
	errorColor := lipgloss.NewStyle().Foreground(ColorError)
	otherColor := lipgloss.NewStyle().Foreground(ColorGrey)
	warnColor := lipgloss.NewStyle().Foreground(ColorWarn)

	if statusCode >= 200 && statusCode < 300 {
		return successColor
	} else if statusCode >= 400 && statusCode < 500 {
		return warnColor
	} else if statusCode >= 500 {
		return errorColor
	}

	return otherColor
}
