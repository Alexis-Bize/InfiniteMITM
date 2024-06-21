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

package utilities

import (
	"encoding/hex"
	"fmt"
	"reflect"
	"strings"

	"github.com/charmbracelet/lipgloss"
)

func Contains(slice []string, item string) bool {
	defer func() {
		_ = recover();
	}()

	for _, v := range slice {
		if v == item {
			return true
		}
	}

	return false
}

func InterfaceToMap(i interface{}) map[string]string {
	defer func() {
		_ = recover();
	}()

	v := reflect.ValueOf(i)
	output := make(map[string]string)

	if i == nil {
		return output
	}

	for _, key := range v.MapKeys() {
		ckey := fmt.Sprintf("%v", key.Interface())
		cvalue := fmt.Sprintf("%v", v.MapIndex(key).Interface())
		output[ckey] = cvalue
	}

	return output
}

func WrapText(text string, width int) string {
	if width <= 0 {
		return text
	}

	var sb strings.Builder
	for len(text) > width {
		cut := width
		for cut > 0 && text[cut] != ' ' && text[cut] != '\n' {
			cut--
		}

		if cut == 0 {
			cut = width
		}

		sb.WriteString(text[:cut] + "\n")
		text = text[cut:]
		for len(text) > 0 && (text[0] == ' ' || text[0] == '\n') {
			text = text[1:]
		}
	}

	sb.WriteString(text)
	return sb.String()
}



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

