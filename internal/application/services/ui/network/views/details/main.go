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
	"strings"

	tea "github.com/charmbracelet/bubbletea"
)


type DetailsModel struct {
	activeTab activeTabType
	activeSection activeSectionType
}

type activeTabType string
type activeSectionType string

const (
	RequestTab  activeTabType = "request"
	ResponseTab activeTabType = "response"
)

func NewDetailsModel() DetailsModel {
	m := DetailsModel{
		activeTab: RequestTab,
	}

	return m
}

func (m *DetailsModel) Init() tea.Cmd {
	return nil
}

func (m *DetailsModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd
	return m, cmd
}

func (m *DetailsModel) View() string {
	doc := strings.Builder{}
	return doc.String()
}
