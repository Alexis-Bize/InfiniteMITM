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

package MITMApplicationSelectServersToolUI

import (
	"fmt"
	selectServersTool "infinite-mitm/internal/application/tools/select-servers"
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

type model struct {
	cursor  int

	selectedRegions []string
	pingResults     []selectServersTool.PingResult

	pingStarted   bool
	pingFinished  bool
	savingServers bool
	serverSaved   bool
}

type setPingForServerMsg selectServersTool.PingResult

const (
	MoveUpCommand   = "up"
	MoveDownCommand = "down"
	SaveCommand     = "ctrl+s"
	EnterCommand    = "enter"
	QuitCommand     = "q"
	ExitCommand     = "ctrl+c"
)

var (
	GoodPingThreshold = []int{0, 60}
	AveragePingThreshold = []int{61, 120}
	BadPingThreshold = []int{121}
)

var program *tea.Program
var servers selectServersTool.QOSServers

var (
	serversHeadlineString = "Select Servers"
	pingFailString        = "FAIL"
	waitingString         = "Please wait..."
	savingString          = "Saving..."
	savedString           = "✓ Saved; restart Halo Infinite to apply changes"

	toggleHintString      = "Enter ↵: Select/Unselect"
	moveHintString        = "↑/↓: Move"
	resetHintString       = "Saving an empty selection will restore the default servers list."

	quitString            = fmt.Sprintf("%s: Quit", QuitCommand)
	saveString            = fmt.Sprintf("%s: Save to %s", SaveCommand, resources.GetRootPath())
)

var (
	containerStyle = lipgloss.NewStyle().Padding(1, 2)
	serversHeadlineStyle = lipgloss.NewStyle().
		Foreground(theme.ColorSunsetOrange).
		Bold(true).
		BorderBottom(true).
		BorderBottomForeground(theme.ColorSunsetOrange).
		BorderStyle(lipgloss.ThickBorder())

	selectedRegionStyle = lipgloss.NewStyle().
		Foreground(theme.ColorLightYellow).
		Background(theme.ColorNeonBlue)
	actionStyle = lipgloss.NewStyle().
		Padding(0, 1).
		MarginRight(1).
		Foreground(theme.ColorLight).
		Background(theme.ColorGrey)
	hintStyle = lipgloss.NewStyle().
		MarginRight(1).
		Foreground(theme.ColorGrey)

	pingStyle = lipgloss.NewStyle().Foreground(theme.ColorNormalFg).Bold(true)
	pingGoodStyle = lipgloss.NewStyle().Foreground(theme.ColorSuccess).Bold(true)
	pingAverageStyle = lipgloss.NewStyle().Foreground(theme.ColorWarn).Bold(true)
	pingBadStyle = lipgloss.NewStyle().Foreground(theme.ColorError).Bold(true)
)

func Create() {
	m := model{pingResults: make([]selectServersTool.PingResult, 0)}
	program = tea.NewProgram(m)
	if _, err := program.Run(); err != nil {
		log.Fatalln(errors.Create(errors.ErrFatalException, err.Error()))
	}
}

func (m model) Init() tea.Cmd {
	servers = selectServersTool.GetQOSServers()
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

			go func(v selectServersTool.QOSServer) {
				defer wg.Done()
				rtt := -1
				stats, _ := selectServersTool.GetPingTime(v.ServerURL)
				if len(stats.Rtts) > 0 {
					rtt = int(stats.AvgRtt.Milliseconds())
				}

				program.Send(setPingForServerMsg{
					ServerURL: v.ServerURL,
					Region: v.Region,
					Ping: rtt,
				})
			}(v)
		}

		wg.Wait()
	}()

	m.pingStarted = true
}

func (m *model) setPingForRegion(result selectServersTool.PingResult) tea.Cmd {
	m.pingResults = append(m.pingResults, result)
	if len(m.pingResults) == len(servers) {
		m.pingFinished = true
	}

	return tea.Batch()
}

func (m *model) toggleRegion() {
	if m.savingServers {
		return
	}

	m.serverSaved = false
	target := servers[m.cursor].Region
	if !utilities.Contains(m.selectedRegions, target) {
		m.selectedRegions = append(m.selectedRegions, target)
		return
	}

	for i, v := range m.selectedRegions {
		if v == target {
			m.selectedRegions = append(m.selectedRegions[:i], m.selectedRegions[i+1:]...)
			break
		}
	}
}

func (m *model) saveServers() {
    if m.savingServers {
        return
    }

    m.savingServers = true
    m.serverSaved = false

    var updatedServers = selectServersTool.QOSServers{}
    
    if len(m.selectedRegions) != 0 {
        var bestSelectedServerURL string
        lowestPing := int(^uint(0) >> 1)

        for _, result := range m.pingResults {
            if utilities.Contains(m.selectedRegions, result.Region) && 
               result.Ping != selectServersTool.PingErrorValue && 
               result.Ping < lowestPing {
                lowestPing = result.Ping
                bestSelectedServerURL = result.ServerURL
            }
        }

        for _, server := range servers {
            newServer := selectServersTool.QOSServer{
                Region: server.Region,
            }
            
            if utilities.Contains(m.selectedRegions, server.Region) {
                newServer.ServerURL = server.ServerURL
            } else if bestSelectedServerURL != "" {
                newServer.ServerURL = bestSelectedServerURL
            } else {
                newServer.ServerURL = server.ServerURL
            }
            
            updatedServers = append(updatedServers, newServer)
        }
    }

    if len(updatedServers) == 0 {
        updatedServers = servers
    }

    selectServersTool.WriteLocalQOSServers(updatedServers)
    selectServersTool.WriteMITMFile()

    m.savingServers = false
    m.serverSaved = true
}

func (m *model) moveUp() {
	if m.pingFinished && !m.savingServers {
		m.cursor--
		if m.cursor < 0 {
			m.cursor = len(servers) - 1
		}
	}
}

func (m *model) moveDown() {
	if m.pingFinished && !m.savingServers {
		m.cursor++
		if m.cursor >= len(servers) {
			m.cursor = 0
		}
	}
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		if !m.pingStarted {
			m.pingServers()
		}
	case setPingForServerMsg:
		return m, m.setPingForRegion(selectServersTool.PingResult(msg))
	case tea.KeyMsg:
		switch msg.String() {
		case QuitCommand, ExitCommand:
			sysutilities.ClearTerminal()
			return m, tea.Quit
		case MoveUpCommand:
			m.moveUp()
		case MoveDownCommand:
			m.moveDown()
		case EnterCommand:
			m.toggleRegion()
		case SaveCommand:
			m.saveServers()
		}
	}

	return m, nil
}

func (m model) View() string {
	s := strings.Builder{}
	s.WriteString(serversHeadlineStyle.Render(serversHeadlineString))
	s.WriteString("\n")
	s.WriteString(hintStyle.Render(resetHintString))
	s.WriteString("\n")

	var hints = []string{}
	if m.pingFinished {
		hints = append(hints, hintStyle.Render(moveHintString), hintStyle.Render(toggleHintString))
	} else {
		hints = append(hints, hintStyle.Render(waitingString))
	}

	s.WriteString(lipgloss.JoinHorizontal(
		lipgloss.Top,
		strings.Join(hints, hintStyle.Render("|")),
	))

	s.WriteString("\n\n")

	total := len(servers)
	for i := 0; i < total; i++ {
		content := "[ ]"

		if m.pingFinished {
			if utilities.Contains(m.selectedRegions, servers[i].Region) {
				content = "[x]"
			}

			if m.cursor == i {
				content = selectedRegionStyle.Render(content)
			}
		}

		s.WriteString(content + " ")

		region := servers[i].Region
		pingColor := pingStyle
		pingStatus := "..."

		for _, v := range m.pingResults {
			if v.Region != region {
				continue
			}

			if v.Ping == selectServersTool.PingErrorValue {
				pingColor = pingBadStyle
				pingStatus = pingFailString
			} else {
				if v.Ping >= GoodPingThreshold[0] && v.Ping <= GoodPingThreshold[1] {
					pingColor = pingGoodStyle
				} else if v.Ping >= AveragePingThreshold[0] && v.Ping <= AveragePingThreshold[1] {
					pingColor = pingAverageStyle
				} else if v.Ping >= BadPingThreshold[0] {
					pingColor = pingBadStyle
				}

				pingStatus = fmt.Sprintf("%d ms", v.Ping)
			}
		}

		s.WriteString(fmt.Sprintf("%s: %s", servers[i].Region, pingColor.Render(pingStatus)))
		s.WriteString("\n")
	}

	s.WriteString("\n")

	extraContent := waitingString
	if m.pingFinished {
		if m.savingServers {
			extraContent = savingString
		} else if m.serverSaved {
			extraContent = savedString
		} else {
			extraContent = saveString
		}
	}

	s.WriteString(lipgloss.JoinHorizontal(
		lipgloss.Top,
		actionStyle.Render(quitString),
		actionStyle.Render(extraContent),
	))

	return containerStyle.Render(
		lipgloss.JoinVertical(
			lipgloss.Top,
			s.String(),
		),
	)
}
