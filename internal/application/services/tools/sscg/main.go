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

package MITMApplicationSSCGTool

import (
	"bytes"
	"encoding/json"
	"fmt"
	embed "infinite-mitm/internal/application/embed"
	"infinite-mitm/pkg/errors"
	"infinite-mitm/pkg/sysutilities"
	"infinite-mitm/pkg/utilities"
	"log"
	"os"
	"os/exec"
	"runtime"
	"strconv"
	"strings"
	"time"

	"github.com/atotto/clipboard"
)

type Server struct {
	Region    string `json:"region"`
	ServerURL string `json:"serverUrl"`
}

func Run() {
	qosservers, err := embed.FileSystem.ReadFile("assets/resources/shared/qosservers.json"); if err != nil {
		log.Fatalln(errors.Create(errors.ErrFatalException, err.Error()))
	}

	var servers []Server
	if err := json.Unmarshal([]byte(qosservers), &servers); err != nil {
		fmt.Println("Failed to parse servers:", err)
		os.Exit(1)
	}

	fmt.Println("\033[1mServers:\033[0m")
	responseTimes := make(map[string]int)

	for _, server := range servers {
		pingTime := getPingTime(server.ServerURL)
		fmt.Printf("%s: %s\n", server.Region, pingTime)

		if pingTime != "FAIL" {
			pingTimeInt, _ := strconv.Atoi(strings.Split(pingTime, " ")[0])
			responseTimes[server.ServerURL] = pingTimeInt
		}
	}

	minPing := 999999
	maxPing := 0
	minServerURL := ""
	maxServerURL := ""

	for url, ping := range responseTimes {
		if ping < minPing {
			minPing = ping
			minServerURL = url
		}

		if ping > maxPing {
			maxPing = ping
			maxServerURL = url
		}
	}

	fmt.Print("\n> Enter the \033[1mregions\033[0m you want to \033[1mprioritize\033[0m (separated by space): ")
	var selectedRegions []string
	fmt.Scanln(&selectedRegions)

	var jsonBuffer bytes.Buffer
	jsonBuffer.WriteString("[")

	for i, server := range servers {
		var assignedServerURL = server.ServerURL
		if utilities.Contains(selectedRegions, server.Region) {
			assignedServerURL = minServerURL
		} else {
			assignedServerURL = maxServerURL
		}

		server.ServerURL = assignedServerURL
		serverJSON, _ := json.Marshal(server)
		jsonBuffer.Write(serverJSON)

		if i < len(servers)-1 {
			jsonBuffer.WriteString(",")
		}
	}

	jsonBuffer.WriteString("]")
	finalJSON := jsonBuffer.String()

	fmt.Println()
	copyToClipboard(finalJSON)
	fmt.Println("\033[1m✔ The configuration has been copied to your clipboard.\033[0m")

	time.Sleep(2 * time.Second)
	sysutilities.ClearTerminal()
}

func getPingTime(serverURL string) string {
	cmd := exec.Command("ping", "-c", "1", serverURL)
	output, err := cmd.CombinedOutput()
	if err != nil {
		return "FAIL"
	}

	outputStr := string(output)
	if strings.Contains(outputStr, "time=") {
		timeStr := strings.Split(strings.Split(outputStr, "time=")[1], " ")[0]
		return timeStr + " ms"
	}

	return "FAIL"
}

func copyToClipboard(text string) {
	switch runtime.GOOS {
	case "darwin":
		clipboard.WriteAll(text)
	case "windows":
		cmd := exec.Command("powershell", "-Command", "Set-Clipboard", text)
		cmd.Run()
	default:
		fmt.Println("Unsupported OS:", runtime.GOOS)
		os.Exit(1)
	}
}
