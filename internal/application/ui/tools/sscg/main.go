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
	"infinite-mitm/configs"
	sscg "infinite-mitm/internal/application/tools/sscg"
	"infinite-mitm/pkg/errors"
	"infinite-mitm/pkg/sysutilities"
	"infinite-mitm/pkg/theme"
	"log"
	"strings"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type pingMap map[string]int
type model struct {
	cursor  int
	ready   bool
	pings   pingMap
	regions []string
}

type setPingForRegionMsg struct {
    region string
    ping   int
}

var program *tea.Program
var servers []sscg.Server

var (
	serversHeadlineString = "Servers:"
	pingFailString        = "FAIL"
)

var (
	serversHeadlineStyle = lipgloss.NewStyle().Bold(true)

	pingStyle = lipgloss.NewStyle().Foreground(theme.ColorNormalFg).Bold(true)
	pingGoodStyle = lipgloss.NewStyle().Foreground(theme.ColorSuccess).Bold(true)
	pingAverageStyle = lipgloss.NewStyle().Foreground(theme.ColorWarn).Bold(true)
	pingBadStyle = lipgloss.NewStyle().Foreground(theme.ColorError).Bold(true)
)

func Create() {
	m := model{pings: make(pingMap)}
	program = tea.NewProgram(m)
	if _, err := program.Run(); err != nil {
		log.Fatalln(errors.Create(errors.ErrFatalException, err.Error()))
	}
}

func (m model) Init() tea.Cmd {
	servers = *sscg.GetServers()
	return tea.SetWindowTitle(fmt.Sprintf("%s: %s", configs.GetConfig().Name, "Servers"))
}

func (m *model) pingServers() {
	if m.ready {
		return
	}

	go func() {
		for _, v := range servers {
			pingTime, _ := sscg.GetPingTime(v.ServerURL)
			program.Send(setPingForRegionMsg{region: v.Region, ping: pingTime})
		}
	}()

	m.ready = true
}

func (m *model) setPingForRegion(region string, ping int) tea.Cmd {
	m.pings[region] = ping
	return tea.Batch()
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		if !m.ready {
			m.pingServers()
		}
	case setPingForRegionMsg:
		return m, m.setPingForRegion(msg.region, msg.ping)
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c", "q", "esc":
			sysutilities.ClearTerminal()
			return m, tea.Quit
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
		/*
		if m.ready {
			if m.cursor == i {
				s.WriteString("[✓] ")
			} else {
				s.WriteString("[ ] ")
			}
		}
		*/

		prefix := "├──"
		if i == total - 1 {
			prefix = "└──"
		}

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

		s.WriteString(fmt.Sprintf("%s %s: %s", prefix, servers[i].Region, pingColor.Render(pingStatus)))
		s.WriteString("\n")
	}

	return lipgloss.JoinVertical(
		lipgloss.Top,
		lipgloss.NewStyle().
			Padding(1, 2).
			Render(s.String()),
	)
}
