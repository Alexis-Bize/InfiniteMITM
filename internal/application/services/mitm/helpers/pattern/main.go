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

package InfiniteMITMApplicationMITMServicePatternHelpers

import (
	"infinite-mitm/configs"
	"regexp"
	"strings"
)

type Patterns struct {
	XUID string
	GUID string
	MVAR string
}

var MatchPatterns = Patterns{
	XUID: `xuid\((\d+)\)`,
	GUID: `([0-9a-fA-F]{8}-[0-9a-fA-F]{4}-[0-9a-fA-F]{4}-[0-9a-fA-F]{4}-[0-9a-fA-F]{12})`,
	MVAR: `([a-z0-9-_]+\.mvar)`,
}

func Create(value string) *regexp.Regexp {
	return regexp.MustCompile(Replace(value))
}

func Replace(value string) string {
	value = strings.Replace(value, ":guid", MatchPatterns.GUID, -1)
	value = strings.Replace(value, ":map-mvar", MatchPatterns.MVAR, -1)
	value = strings.Replace(value, ":xuid", MatchPatterns.XUID, -1)
	value = strings.Replace(value, ":infinite-mitm-version", configs.GetConfig().Version, -1)
	return value
}
