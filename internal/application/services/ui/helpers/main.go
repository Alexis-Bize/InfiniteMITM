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

package MITMApplicationUIServiceUIHelpers

import (
	"encoding/hex"
	"fmt"
	"infinite-mitm/configs"
	"infinite-mitm/pkg/resources"
	"os"
	"path/filepath"
	"strings"

	"github.com/atotto/clipboard"
	"github.com/charmbracelet/lipgloss"
	"github.com/ncruces/zenity"
)

func FormatHexView(data []byte, width int) string {
	bytesPerLine := (width - 10)
	if bytesPerLine <= 10 {
		bytesPerLine = 10
	}

	bytesPerLine = width / 4

	var output string

	for i := 0; i < len(data); i += bytesPerLine {
		end := i + bytesPerLine
		if end > len(data) {
			end = len(data)
		}

		hexPart := hex.EncodeToString(data[i:end])
		hexPart = fmt.Sprintf("%-*s", bytesPerLine*2, hexPart)
		i2hex := fmt.Sprintf("%08x", i)
		asciiPart := toASCII(data[i:end], 32, 126)
		line := fmt.Sprintf("%s  %s  %s", lipgloss.NewStyle().Bold(true).Render(i2hex), hexPart, asciiPart)
		output += line

		if end != len(data) {
			output += "\n"
		}
	}

	return output
}

func toASCII(data []byte, minByte byte, maxByte byte) string {
	var dot rune = '.'
	var ascii []rune

	for _, b := range data {
		if b < minByte || b > maxByte {
			ascii = append(ascii, dot)
		} else {
			ascii = append(ascii, rune(b))
		}
	}

	return string(ascii)
}

func CopyToClipboard(data string) {
	clipboard.WriteAll(data)
}

func SaveToDisk(data []byte, filename string, contentType string) {
	defer func() {
		_ = recover();
	}()

	contentType = strings.Split(contentType, ";")[0]
	var mimeExtensions = map[string]string{
		"application/json":         "json",
		"application/xml":          "xml",
		"text/html":                "html",
		"text/plain":               "txt",
		"image/jpeg":               "jpeg",
		"image/jpg":                "jpg",
		"image/png":                "png",
		"image/gif":                "gif",
		"application/octet-stream": "bin",
	}

	extension := mimeExtensions[contentType]
	if extension == "" {
		extension = "bin"
	}

	outputDir := filepath.Join(resources.GetDownloadsDirPath())
	if _, err := os.Stat(outputDir); os.IsNotExist(err) {
		outputDir = configs.GetConfig().Extra.ProjectDir;
	}

	filename = fmt.Sprintf("%s.%s", filename, extension)
	filePath, err := zenity.SelectFileSave(
		zenity.Title("Save body content"),
		zenity.Filename(filepath.Join(outputDir, filename)),
		zenity.ConfirmOverwrite(),
	)

	if err == nil {
		os.WriteFile(filePath, data, 0644)
	}
}
