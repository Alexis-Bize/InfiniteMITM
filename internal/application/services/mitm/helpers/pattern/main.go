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

package MITMApplicationMITMServicePatternHelpers

import (
	"regexp"
	"strings"
)

type Patterns struct {
	XUID    string
	GUID    string
	MVAR    string
	SANDBOX string
}

var MatchParameters = Patterns{
	XUID:    ":xuid",
	GUID:    ":guid",
	MVAR:    ":map-mvar",
	SANDBOX: ":sandbox",
}

var MatchPatterns = Patterns{
	XUID:    `xuid\((\d+)\)`,
	GUID:    `([0-9a-fA-F]{8}-[0-9a-fA-F]{4}-[0-9a-fA-F]{4}-[0-9a-fA-F]{4}-[0-9a-fA-F]{12})`,
	MVAR:    `([a-z0-9-_]+\.mvar)`,
	SANDBOX: `(retail|test|beta|beta-test)`,
}

func Create(value string) *regexp.Regexp {
	return regexp.MustCompile(Replace(value))
}

func Replace(value string) string {
	value = strings.Replace(value, MatchParameters.XUID, MatchPatterns.XUID, -1)
	value = strings.Replace(value, MatchParameters.GUID, MatchPatterns.GUID, -1)
	value = strings.Replace(value, MatchParameters.MVAR, MatchPatterns.MVAR, -1)
	value = strings.Replace(value, MatchParameters.SANDBOX, MatchPatterns.SANDBOX, -1)
	return value
}

func Process(reqp string, resp string) string {
	re1 := Replace(reqp)
	re2 := resp

	// xuid
	XUIDExtractor := regexp.MustCompile(MatchPatterns.XUID)
	XUIDMatches := XUIDExtractor.FindAllStringSubmatch(re1, -1)
	var XUIDValues []string
	for _, match := range XUIDMatches {
		if len(match) > 0 {
			XUIDValues = append(XUIDValues, match[0])
		}
	}

	if len(XUIDValues) > 0 {
		re2 = strings.Replace(re2, MatchParameters.XUID, XUIDValues[0], -1)
	}

	// guid
	GUIDExtractor := regexp.MustCompile(MatchPatterns.GUID)
	GUIDMatches := GUIDExtractor.FindAllStringSubmatch(re1, -1)
	var GUIDValues []string
	for _, match := range GUIDMatches {
		if len(match) > 0 {
			GUIDValues = append(GUIDValues, match[0])
		}
	}

	if len(GUIDValues) > 0 {
		re2 = strings.Replace(re2, MatchParameters.GUID, GUIDValues[0], -1)
	}

	// sandbox
	SANDBOXExtractor := regexp.MustCompile(MatchPatterns.SANDBOX)
	SANDBOXMatches := SANDBOXExtractor.FindAllStringSubmatch(re1, -1)
	var SANDBOXValues []string
	for _, match := range SANDBOXMatches {
		if len(match) > 0 {
			SANDBOXValues = append(SANDBOXValues, match[0])
		}
	}

	if len(SANDBOXValues) > 0 {
		re2 = strings.Replace(re2, MatchParameters.SANDBOX, SANDBOXValues[0], -1)
	}

	return re2
}
