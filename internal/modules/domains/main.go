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

package MITMDomainsModule

type Domain = string

type Domains struct {
	Root      Domain
	Blobs     Domain
	Authoring Domain
	Discovery Domain
	HaloStats Domain
	Settings  Domain
	GameCMS   Domain
	Economy   Domain
	Lobby     Domain
}

const (
	Root      Domain  = "svc.halowaypoint.com:443"
	Blobs     Domain  = "blobs-infiniteugc.svc.halowaypoint.com:443"
	Authoring Domain  = "authoring-infiniteugc.svc.halowaypoint.com:443"
	Discovery Domain  = "discovery-infiniteugc.svc.halowaypoint.com:443"
	HaloStats Domain  = "halostats.svc.halowaypoint.com:443"
	Settings  Domain  = "settings.svc.halowaypoint.com:443"
	GameCMS   Domain  = "gamecms-hacs.svc.halowaypoint.com:443"
	Economy   Domain  = "economy.svc.halowaypoint.com:443"
	Lobby     Domain  = "lobby-hi.svc.halowaypoint.com:443"
)

var HaloWaypointSVCDomains = Domains{
	Root:      Root,
	Blobs:     Blobs,
	Authoring: Authoring,
	Discovery: Discovery,
	HaloStats: HaloStats,
	Settings:  Settings,
	GameCMS:   GameCMS,
	Economy:   Economy,
	Lobby:     Lobby,
}
