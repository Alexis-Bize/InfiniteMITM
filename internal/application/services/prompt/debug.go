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

package MITMApplicationPromptService

import (
	"fmt"
	MITMApplicationMITMServicePatternHelpers "infinite-mitm/internal/application/services/mitm/helpers/pattern"
	"regexp"
	"strconv"
)

func DebugPattern() error {

	e := "/audiences/:sandbox/players/:xuid/active/:guid/:guid"
	f := "/audiences/retail/players/xuid(1234)/active/3de2648f-5488-4fdf-be70-ce4a76164df4/22c1153b-7289-4532-9eef-276c46381ab4"
	g := "/audiences/retail/players/xuid(1234)/active/$3/22c1153b-7289-4532-9eef-276c46381ab4"
	fmt.Println(e)
	fmt.Println(f)
	fmt.Println(g)
	fmt.Println("")

	e = MITMApplicationMITMServicePatternHelpers.Replace(e)
	fmt.Println(e)
	fmt.Println("")

	rx1 := regexp.MustCompile(e)
	m1 := rx1.FindAllStringSubmatch(f, -1)

	var matches []string
	if len(m1) > 0 {
		for _, match := range m1[0][1:] {
			if len(match) > 0 {
				matches = append(matches, match)
				fmt.Println(match)
			}
		}
	}

	re := regexp.MustCompile(`\$\d+`)
	result := re.ReplaceAllStringFunc(g, func(m string) string {
        // Extract the number from the match and convert it to an index
        index, err := strconv.Atoi(m[1:])
        if err != nil || index < 1 || index > len(matches) {
            // If the index is out of range or conversion fails, return the original match
            return m
        }
        // Return the replacement value
        return matches[index-1]
    })

	fmt.Println("")
	fmt.Println(result)

	return nil
}
