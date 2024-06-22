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

package MITMApplicationSSCGToolUI

import (
	"fmt"
	sscg "infinite-mitm/internal/application/tools/sscg"
	"infinite-mitm/pkg/errors"
	"infinite-mitm/pkg/resources"
	"infinite-mitm/pkg/sysutilities"
	"infinite-mitm/pkg/theme"
	"infinite-mitm/pkg/utilities"
	"log"
	"strings"
	"sync"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type pings map[string]int
type model struct {
	cursor  int
	regions []string
	pings   pings

	pingStarted  bool
	pingFinished bool
}

type setPingForRegionMsg struct {
	region string
	ping   int
}

var program *tea.Program
var servers []sscg.Server

var (
	serversHeadlineString = "Select Servers:"
	pingFailString        = "FAIL"
)

var (
	containerStyle = lipgloss.NewStyle().Padding(1, 2)
	serversHeadlineStyle = lipgloss.NewStyle().Bold(true)
	selectedServerStyle = lipgloss.NewStyle().
		Foreground(theme.ColorLightYellow).
		Background(theme.ColorNeonBlue)

	pingStyle = lipgloss.NewStyle().Foreground(theme.ColorNormalFg).Bold(true)
	pingGoodStyle = lipgloss.NewStyle().Foreground(theme.ColorSuccess).Bold(true)
	pingAverageStyle = lipgloss.NewStyle().Foreground(theme.ColorWarn).Bold(true)
	pingBadStyle = lipgloss.NewStyle().Foreground(theme.ColorError).Bold(true)
)

func Create() {
	m := model{pings: make(pings)}
	program = tea.NewProgram(m)
	if _, err := program.Run(); err != nil {
		log.Fatalln(errors.Create(errors.ErrFatalException, err.Error()))
	}
}

func (m model) Init() tea.Cmd {
	servers = *sscg.GetServers()
	return nil
}

func (m *model) pingServers() {
	if m.pingStarted {
		return
	}

	go func() {
		var wg sync.WaitGroup

		for _, v := range servers {
			wg.Add(1)

			go func(v sscg.Server) {
				defer wg.Done()
				pingTime, _ := sscg.GetPingTime(v.ServerURL)
				program.Send(setPingForRegionMsg{region: v.Region, ping: pingTime})
			}(v)
		}

		wg.Wait()
	}()

	m.pingStarted = true
}

func (m *model) setPingForRegion(region string, ping int) tea.Cmd {
	m.pings[region] = ping

	if len(m.pings) == len(servers) {
		m.pingFinished = true
	}

	return tea.Batch()
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		if !m.pingStarted {
			m.pingServers()
		}
	case setPingForRegionMsg:
		return m, m.setPingForRegion(msg.region, msg.ping)
	case tea.KeyMsg:
		switch msg.String() {
		case "enter":
			target := servers[m.cursor].Region
			if utilities.Contains(m.regions, target) {
				index := -1
				for i, v := range m.regions {
					if v == target {
						index = i
						break
					}
				}

				if index != -1 {
					m.regions = append(m.regions[:index], m.regions[index+1:]...)
				}
			} else {
				m.regions = append(m.regions, target)
			}
		case "ctrl+c", "q", "esc":
			sysutilities.ClearTerminal()
			return m, tea.Quit
		case "down", "j":
			if m.pingFinished {
				m.cursor++
				if m.cursor >= len(servers) {
					m.cursor = 0
				}
			}
		case "up", "k":
			if m.pingFinished {
				m.cursor--
				if m.cursor < 0 {
					m.cursor = len(servers) - 1
				}
			}
		}
	}

	return m, nil
}

func (m model) View() string {
	s := strings.Builder{}
	s.WriteString(serversHeadlineStyle.Render(serversHeadlineString))
	s.WriteString("\n")
	total := len(servers)

	for i := 0; i < total; i++ {
		content := "[!]"

		if m.pingFinished {
			if utilities.Contains(m.regions, servers[i].Region) {
				content = "[✔]"
			} else {
				content = "[ ]"
			}

			if m.cursor == i {
				content = selectedServerStyle.Render(content)
			}
		}

		s.WriteString(content + " ")

		region := servers[i].Region
		pingColor := pingStyle
		pingStatus := "..."

		if ping, exists := m.pings[region]; exists {
			if ping == -1 {
				pingColor = pingBadStyle
				pingStatus = pingFailString
			} else {
				if ping >= 0 && ping <= 75 {
					pingColor = pingGoodStyle
				} else if ping >= 76 && ping <= 120 {
					pingColor = pingAverageStyle
				} else if ping >= 121 {
					pingColor = pingBadStyle
				}

				pingStatus = fmt.Sprintf("%d ms", ping)
			}
		}

		s.WriteString(fmt.Sprintf("%s: %s", servers[i].Region, pingColor.Render(pingStatus)))
		s.WriteString("\n")
	}

	s.WriteString("\n")
	s.WriteString(lipgloss.JoinHorizontal(
		lipgloss.Top,
		lipgloss.NewStyle().
			Padding(0, 1).MarginRight(1).
			Foreground(theme.ColorLight).
			Background(theme.ColorGrey).
			Render("q: Quit"),
		lipgloss.NewStyle().
			Padding(0, 1).MarginRight(1).
			Foreground(theme.ColorLight).
			Background(theme.ColorGrey).
			Render(fmt.Sprintf("ctlr+s: Save to %s", resources.GetRootPath())),
	))

	return containerStyle.Render(
		lipgloss.JoinVertical(
			lipgloss.Top,
			s.String(),
		),
	)
}
