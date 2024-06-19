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
	"infinite-mitm/pkg/resources"
	"regexp"
	"sort"
	"strconv"
	"strings"
)

type Patterns struct {
	SANDBOX, XUID, GUID, MVAR, CGUI, EGV, TITLE string
	BLOBS_SVC, STATS_SVC, GAMECMS_SVC, ECONOMY_SVC, SETTINGS_SVC, AUTHORING_SVC, DISCOVERY_SVC string
	XML, BOND, JSON string
	ALL, END, OR string
}

const MITMDirPrefix = ":mitm-dir"

var (
	mitmDirs     map[string]string
	mitmDirsKeys []string
	replacer     *strings.Replacer
)

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

	ALL: ":all",
	END: ":end",
	OR:  ":or",
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

	ALL: `(.*)`,
	END: `$`,
}

func init() {
	replacer = strings.NewReplacer(
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
		MatchParameters.END, MatchPatterns.END,

		":mitm-version", configs.GetConfig().Version,
	)

	mitmDirs = make(map[string]string)
	mitmDirs[MITMDirPrefix] = configs.GetConfig().Extra.ProjectDir
	for k, v := range resources.GetDirPaths() {
		if k != "smartcache" && k != "downloads" {
			key := strings.Join([]string{MITMDirPrefix, k}, "-")
			mitmDirs[key] = v
		}
	}

	keys := make([]string, 0, len(mitmDirs))
	for k := range mitmDirs {
		keys = append(keys, k)
	}

	sort.Slice(keys, func(i, j int) bool {
		return len(keys[i]) > len(keys[j])
	})

	mitmDirsKeys = keys
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

func Create(domain domains.DomainType, path string) *regexp.Regexp {
	if path == "" || !strings.HasPrefix(path, "/") {
		path = "/" + path
	}

	path = strings.Replace(path, "*", ":all", -1)
	path = strings.Replace(path, "$", ":end", -1)
	path = ReplaceParameters(escapeExceptParentheses(path))

	re := regexp.MustCompile(`(?i)` + regexp.QuoteMeta(domain) + path)
	return re
}

func ReplaceParameters(value string) string {
	for _, k := range mitmDirsKeys {
		value = strings.ReplaceAll(value, k, mitmDirs[k])
	}

	return replacer.Replace(value)
}

func ReplaceMatches(target string, matches []string) string {
	re := regexp.MustCompile(`\$\d+`)
	return re.ReplaceAllStringFunc(target, func(m string) string {
		index, err := strconv.Atoi(m[1:])
		if err != nil || index < 1 || index > len(matches) {
			return ""
		}

		return matches[index - 1]
	})
}

func escapeExceptParentheses(str string) string {
	var result strings.Builder

	re := regexp.MustCompile(`\([^)]*\)`)
	lastEnd := 0

	for _, loc := range re.FindAllStringIndex(str, -1) {
		result.WriteString(regexp.QuoteMeta(str[lastEnd:loc[0]]))
		result.WriteString(str[loc[0]:loc[1]])
		lastEnd = loc[1]
	}

	result.WriteString(regexp.QuoteMeta(str[lastEnd:]))
	return result.String()
}
