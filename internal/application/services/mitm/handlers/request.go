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
	"net/http"
	"regexp"
	"sync"

	"github.com/elazarl/goproxy"
)

type RequestHandlerStruct struct {
	Match goproxy.ReqConditionFunc
	Fn    func(req *http.Request, ctx *goproxy.ProxyCtx) (*http.Request, *http.Response)
}

type ResponseHandlerStruct struct {
	Match goproxy.ReqConditionFunc
	Fn    func(resp *http.Response, ctx *goproxy.ProxyCtx) *http.Response
}

var (
	cacheUserFavoriteFilmsMutex sync.Mutex
	userFavoriteFilms = HaloWaypointLibRequestModuleDomains.FavoriteFilmsResult{}
)

func HandleHaloWaypointRequests() RequestHandlerStruct {
	target := regexp.MustCompile(`(?i)` + InfiniteMITMDomainsModule.HaloWaypointSVCDomains.Root)

	return RequestHandlerStruct{
		Match: goproxy.UrlMatches(target),
		Fn: func(req *http.Request, ctx *goproxy.ProxyCtx) (*http.Request, *http.Response) {
			return req, nil
		},
	}
}

func Dirty__CacheUserFavoriteFilms() RequestHandlerStruct {
	target := regexp.MustCompile(`(?i)` + InfiniteMITMDomainsModule.HaloWaypointSVCDomains.Authoring + `/hi/players/xuid\(\d+\)/favorites/films`)

	return RequestHandlerStruct{
		Match: goproxy.UrlMatches(target),
		Fn: func(req *http.Request, ctx *goproxy.ProxyCtx) (*http.Request, *http.Response) {
			cacheUserFavoriteFilmsMutex.Lock()
			defer cacheUserFavoriteFilmsMutex.Unlock()

			if req.Method == http.MethodOptions || req.Method == http.MethodHead {
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

			userFavoriteFilms = unmarshal
			return req, nil
		},
	}
}
