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

package InfiniteMITMApplicationServiceModule

import (
	InfiniteMITMDomainsModule "infinite-mitm/internal/modules"
	"infinite-mitm/pkg/modules/utilities/request"
	"net/http"
	"regexp"

	"gopkg.in/elazarl/goproxy.v1"
)

func PatchBookmarkedFilms() HandlerStruct {
	url := regexp.MustCompile(InfiniteMITMDomainsModule.HaloWaypointSVCDomains.Discovery + "/hi/films/[a-fA-F0-9]{8}-[a-fA-F0-9]{4}-[a-fA-F0-9]{4}-[a-fA-F0-9]{4}-[a-fA-F0-9]{12}")

	return HandlerStruct{
		Match: goproxy.UrlMatches(url),
		Fn: func(req *http.Request, ctx *goproxy.ProxyCtx) (*http.Request, *http.Response) {
			spartanToken := request.ExtractSpartanTokenAuthorization(req)

			if spartanToken == "" {
				return req, nil
			}

			return req, nil
		},
	}
}
