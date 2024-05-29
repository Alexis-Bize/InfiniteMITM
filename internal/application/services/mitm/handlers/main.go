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

package InfiniteMITMApplicationMITMServiceHandlers

import (
	"encoding/json"
	InfiniteMITMDomainsModule "infinite-mitm/internal/modules"
	HaloWaypointLibRequestModule "infinite-mitm/pkg/libs/halowaypoint/modules/request"
	HaloWaypointLibRequestModuleDomains "infinite-mitm/pkg/libs/halowaypoint/modules/request/domains"
	ErrorsModule "infinite-mitm/pkg/modules/errors"
	UtilitiesRequestModule "infinite-mitm/pkg/modules/utilities/request"
	"io"
	"log"
	"net/http"
	"regexp"
	"sync"

	"github.com/elazarl/goproxy"
)

type HandlerStruct struct {
	Match goproxy.ReqConditionFunc
	Fn    func(req *http.Request, ctx *goproxy.ProxyCtx) (*http.Request, *http.Response)
}

var (
	favoriteFilmsJSONCache     = HaloWaypointLibRequestModuleDomains.FavoriteFilmsResult{}
	spectateFilmEstimatedIndex = -1
	favoriteFilmsSyncMutex     sync.Mutex
	spectateFilmSyncMutex      sync.Mutex
)

func HandleHaloWaypointRequests() HandlerStruct {
	target := regexp.MustCompile(`(?i)` + InfiniteMITMDomainsModule.HaloWaypointSVCDomains.Root)

	return HandlerStruct{
		Match: goproxy.UrlMatches(target),
		Fn: func(req *http.Request, ctx *goproxy.ProxyCtx) (*http.Request, *http.Response) {
			if req.Method != http.MethodOptions {
				log.Printf("[%s] %s", req.Method, req.URL.String())
			}

			return req, nil
		},
	}
}

func CacheBookmarkedFilms() HandlerStruct {
	target := regexp.MustCompile(`(?i)` + InfiniteMITMDomainsModule.HaloWaypointSVCDomains.Authoring + `/hi/players/xuid\(\d+\)/favorites/films`)

	return HandlerStruct{
		Match: goproxy.UrlMatches(target),
		Fn: func(req *http.Request, ctx *goproxy.ProxyCtx) (*http.Request, *http.Response) {
			favoriteFilmsSyncMutex.Lock()
			defer favoriteFilmsSyncMutex.Unlock()

			if req.Method == http.MethodOptions {
				return req, nil
			}

			replay, err := UtilitiesRequestModule.ReplayRequestWithJSONAccept(req)
			if err != nil {
				return req, nil
			}
			defer replay.Body.Close()

			err = HaloWaypointLibRequestModule.ValidateResponseStatusCode(replay)
			if err != nil {
				return req, nil
			}

			body, err := io.ReadAll(replay.Body)
			if err != nil {
				ErrorsModule.Log(ErrorsModule.ErrIOReadException, err.Error())
				return req, nil
			}

			var unmarshal HaloWaypointLibRequestModuleDomains.FavoriteFilmsResult
			if err := json.Unmarshal(body, &unmarshal); err != nil {
				ErrorsModule.Log(ErrorsModule.ErrJSONUnmarshalException, err.Error())
				return req, nil
			}

			spectateFilmEstimatedIndex = -1
			favoriteFilmsJSONCache = unmarshal
			return req, nil
		},
	}
}

func DirtyFixInvalidMatchSpectateID() HandlerStruct {
	target := regexp.MustCompile(`(?i)` + InfiniteMITMDomainsModule.HaloWaypointSVCDomains.Discovery + `/hi/films/matches/00000000-0000-0000-0000-000000000000/spectate`)

	return HandlerStruct{
		Match: goproxy.UrlMatches(target),
		Fn: func(req *http.Request, ctx *goproxy.ProxyCtx) (*http.Request, *http.Response) {
			spectateFilmSyncMutex.Lock()
			defer spectateFilmSyncMutex.Unlock()

			if req.Method == http.MethodOptions {
				return req, nil
			}

			totalFilms := len(favoriteFilmsJSONCache.Results)
			if totalFilms == 0 {
				return req, nil
			}

			spectateFilmEstimatedIndex++
			exists := spectateFilmEstimatedIndex >= 0 && spectateFilmEstimatedIndex < totalFilms

			if !exists {
				// will fallback to the 1st item
				spectateFilmEstimatedIndex = 0
			}

			userAgent := UtilitiesRequestModule.ExtractHeaderValue(req, "user-agent")
			spartanToken := UtilitiesRequestModule.ExtractHeaderValue(req, "x-343-authorization-spartan")
			telemetryID := UtilitiesRequestModule.ExtractHeaderValue(req, "343-telemetry-session-id")
			flightID := UtilitiesRequestModule.ExtractHeaderValue(req, "343-clearance")

			film, err := HaloWaypointLibRequestModuleDomains.GetFilmByAssetID(HaloWaypointLibRequestModule.RequestAttributes{
				UserAgent:    userAgent,
				SpartanToken: spartanToken,
				ExtraHeaders: map[string]string{
					"343-Telemetry-Session-Id": telemetryID,
					"343-Clearance":            flightID,
				},
			}, favoriteFilmsJSONCache.Results[spectateFilmEstimatedIndex].AssetID)

			if err != nil {
				return req, nil
			}

			matchID := film.CustomData.MatchID
			req.URL.Path = "/hi/films/matches/"+matchID+"/spectate"

			return req, nil
		},
	}
}
