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

type RequestHandlerStruct struct {
	Match goproxy.ReqConditionFunc
	Fn    func(req *http.Request, ctx *goproxy.ProxyCtx) (*http.Request, *http.Response)
}

type ResponseHandlerStruct struct {
	Match goproxy.ReqConditionFunc
	Fn    func(resp *http.Response, ctx *goproxy.ProxyCtx) *http.Response
}

func HandleRootRequests() RequestHandlerStruct {
	target := regexp.MustCompile(`(?i)` + regexp.QuoteMeta(domains.HaloWaypointSVCDomains.Root))

	return RequestHandlerStruct{
		Match: goproxy.UrlMatches(target),
		Fn: func(req *http.Request, ctx *goproxy.ProxyCtx) (*http.Request, *http.Response) {
			customCtx := context.ContextHandler(ctx)
			uuid := customCtx.GetUserData("uuid").(string)
			if uuid == "" {
				return req, nil
			}

			var bodyBytes []byte
			if req.Body != nil {
				bodyBytes, _ = io.ReadAll(req.Body)
			}

			req.Body = io.NopCloser(bytes.NewBuffer(bodyBytes))

			proxified := customCtx.GetUserData("proxified").(map[string]bool)
			details := events.StringifyRequestEventData(events.ProxyRequestEventData{
				ID: uuid,
				URL: req.URL.String(),
				Method: req.Method,
				Headers: request.HeadersToMap(req.Header),
				Body: bodyBytes,
				Proxified: proxified["req"],
			})

			event.MustFire(events.ProxyRequestSent, event.M{"details": details})
			return req, nil
		},
	}
}
