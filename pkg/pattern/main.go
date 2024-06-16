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

package pattern

import (
	"infinite-mitm/configs"
	"infinite-mitm/pkg/domains"
	"regexp"
	"strconv"
	"strings"
)

type Patterns struct {
	SANDBOX, XUID, GUID, MVAR, CGUI, EGV, TITLE string
	BLOBS_SVC, STATS_SVC, GAMECMS_SVC, ECONOMY_SVC, SETTINGS_SVC, AUTHORING_SVC, DISCOVERY_SVC string
	XML, BOND, JSON string
	ALL, STOP string
}

var MatchParameters = Patterns{
	EGV:     ":egv-bin",
	XUID:    ":xuid",
	GUID:    ":guid",
	MVAR:    ":map-mvar",
	CGUI:    ":cgui-bin",
	TITLE:   ":title",
	SANDBOX: ":sandbox",

	BLOBS_SVC:      ":blobs-svc",
	STATS_SVC:      ":stats-svc",
	GAMECMS_SVC:    ":gamecms-svc",
	ECONOMY_SVC:    ":economy-svc",
	SETTINGS_SVC:   ":settings-svc",
	AUTHORING_SVC:  ":authoring-svc",
	DISCOVERY_SVC:  ":discovery-svc",

	XML:  ":ct-xml",
	BOND: ":ct-bond",
	JSON: ":ct-json",

	ALL:  ":\\*",
	STOP: ":\\$",
}

var MatchPatterns = Patterns{
	EGV:     `([a-z0-9-_]+\.bin)`,
	GUID:    `([0-9a-fA-F]{8}-[0-9a-fA-F]{4}-[0-9a-fA-F]{4}-[0-9a-fA-F]{4}-[0-9a-fA-F]{12})`,
	XUID:    `xuid\((\d+)\)`,
	MVAR:    `([a-z0-9-_]+\.mvar)`,
	CGUI:    `([a-z]+_CustomGamesUIMarkup_[a-z]+\.bin)`,
	SANDBOX: `(retail|test|beta|beta-test)`,
	TITLE:   `(hi|hipnk|higrn|hired|hipur|hiorg|hiblu|hi343)`,

	BLOBS_SVC:      domains.DomainToBaseUrl(domains.Blobs),
	STATS_SVC:      domains.DomainToBaseUrl(domains.HaloStats),
	GAMECMS_SVC:    domains.DomainToBaseUrl(domains.GameCMS),
	ECONOMY_SVC:    domains.DomainToBaseUrl(domains.Economy),
	SETTINGS_SVC:   domains.DomainToBaseUrl(domains.Settings),
	AUTHORING_SVC:  domains.DomainToBaseUrl(domains.Authoring),
	DISCOVERY_SVC:  domains.DomainToBaseUrl(domains.Discovery),

	XML:  "application/xml",
	BOND: "application/x-bond-compact-binary",
	JSON: "application/json",

	ALL:  `(.*)`,
	STOP: `$`,
}

func Create(value string) *regexp.Regexp {
	return regexp.MustCompile(ReplaceParameters(value))
}

func Match(re *regexp.Regexp, value string) []string {
	var matches []string

	submatches := re.FindAllStringSubmatch(value, -1)
	if len(submatches) > 0 {
		for _, match := range submatches[0][1:] {
			if len(match) > 0 {
				matches = append(matches, match)
			}
		}
	}

	return matches
}

func ReplaceParameters(value string) string {
	replacer := strings.NewReplacer(
		MatchParameters.SANDBOX, MatchPatterns.SANDBOX,
		MatchParameters.XUID, MatchPatterns.XUID,
		MatchParameters.GUID, MatchPatterns.GUID,
		MatchParameters.MVAR, MatchPatterns.MVAR,
		MatchParameters.CGUI, MatchPatterns.CGUI,
		MatchParameters.EGV, MatchPatterns.EGV,
		MatchParameters.TITLE, MatchPatterns.TITLE,

		MatchParameters.BLOBS_SVC, MatchPatterns.BLOBS_SVC,
		MatchParameters.STATS_SVC, MatchPatterns.STATS_SVC,
		MatchParameters.GAMECMS_SVC, MatchPatterns.GAMECMS_SVC,
		MatchParameters.ECONOMY_SVC, MatchPatterns.ECONOMY_SVC,
		MatchParameters.SETTINGS_SVC, MatchPatterns.SETTINGS_SVC,
		MatchParameters.AUTHORING_SVC, MatchPatterns.AUTHORING_SVC,
		MatchParameters.DISCOVERY_SVC, MatchPatterns.DISCOVERY_SVC,

		MatchParameters.XML, MatchPatterns.XML,
		MatchParameters.BOND, MatchPatterns.BOND,
		MatchParameters.JSON, MatchPatterns.JSON,

		MatchParameters.ALL, MatchPatterns.ALL,
		MatchParameters.STOP, MatchPatterns.STOP,

		"*", MatchPatterns.ALL,
		"$", MatchPatterns.STOP,

		":mitm-dir", configs.GetConfig().Extra.ProjectDir,
		":mitm-version", configs.GetConfig().Version,
	)

	return replacer.Replace(value)
}

func ReplaceMatches(target string, matches []string) string {
	re := regexp.MustCompile(`\$\d+`)
	return re.ReplaceAllStringFunc(target, func(m string) string {
		index, err := strconv.Atoi(m[1:])
		if err != nil || index < 1 || index > len(matches) {
			return ""
		}

		return matches[index-1]
	})
}
