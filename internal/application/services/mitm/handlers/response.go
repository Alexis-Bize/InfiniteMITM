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

package MITMApplicationMITMServiceHandlers

import (
	"bytes"
	events "infinite-mitm/internal/application/events"
	context "infinite-mitm/internal/application/services/mitm/helpers/context"
	domains "infinite-mitm/internal/modules/domains"
	request "infinite-mitm/pkg/modules/utilities/request"
	"io"
	"net/http"
	"regexp"

	"github.com/elazarl/goproxy"
	"github.com/gookit/event"
)

func HandleRootResponses() ResponseHandlerStruct {
	target := regexp.MustCompile(`(?i)` + regexp.QuoteMeta(domains.HaloWaypointSVCDomains.Root))

	return ResponseHandlerStruct{
		Match: goproxy.UrlMatches(target),
		Fn: func(resp *http.Response, ctx *goproxy.ProxyCtx) *http.Response {
			if resp.Request.Method == http.MethodOptions {
				return resp
			}

			customCtx := context.ContextHandler(ctx)
			uuid := customCtx.GetUserData("uuid").(string)
			if uuid == "" {
				return resp
			}

			var bodyBytes []byte
			if resp.Body != nil {
				bodyBytes, _ = io.ReadAll(resp.Body)
			}

			resp.Body = io.NopCloser(bytes.NewBuffer(bodyBytes))

			proxified := customCtx.GetUserData("proxified").(map[string]bool)
			details := events.StringifyResponseEventData(events.ProxyResponseEventData{
				ID: uuid,
				Status: resp.StatusCode,
				Headers: request.HeadersToMap(resp.Header),
				Body: bodyBytes,
				Proxified: proxified["resp"],
			})

			event.MustFire(events.ProxyResponseReceived, event.M{"details": details})
			return resp
		},
	}
}
