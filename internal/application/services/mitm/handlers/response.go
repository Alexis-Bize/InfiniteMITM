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
	networkTable "infinite-mitm/internal/application/services/ui/network"
	domains "infinite-mitm/internal/modules/domains"
	"net/http"
	"regexp"

	"github.com/elazarl/goproxy"
)

func HandleHaloWaypointResponses() ResponseHandlerStruct {
	target := regexp.MustCompile(`(?i)` + regexp.QuoteMeta(domains.HaloWaypointSVCDomains.Root))

	return ResponseHandlerStruct{
		Match: goproxy.UrlMatches(target),
		Fn: func(resp *http.Response, ctx *goproxy.ProxyCtx) *http.Response {
			if resp.Request.Method != http.MethodOptions {
				networkTable.Push(resp.Request.Method, resp.StatusCode, resp.Request.URL.String())
			}

			return resp
		},
	}
}
