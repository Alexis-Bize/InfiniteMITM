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
	eventsService "infinite-mitm/internal/application/services/events"
	context "infinite-mitm/internal/application/services/mitm/modules/context"
	"infinite-mitm/pkg/mitm"
	"infinite-mitm/pkg/request"
	"infinite-mitm/pkg/smartcache"
	"io"
	"net/http"

	"github.com/elazarl/goproxy"
	"github.com/gookit/event"
)

func HandleRequest(options mitm.TrafficOptions, req *http.Request, resp *http.Response, ctx *goproxy.ProxyCtx) (*http.Request, *http.Response) {
	if req.Method == http.MethodOptions {
		return req, resp
	}

	customCtx := context.ContextHandler(ctx); if customCtx == nil {
		return req, resp
	}

	uuid := getUUID(customCtx)
	smartCache := getSmartCache(customCtx)
	if options.TrafficDisplay == mitm.TrafficSmartCache && smartCache == nil {
		return req, resp
	}

	isProxified := isRequestProxified(customCtx) || isResponseProxified(customCtx)
	if options.TrafficDisplay == mitm.TrafficOverrides && !isProxified && smartCache == nil {
		return req, resp
	}

	var smartCachedItem *smartcache.SmartCacheItem

	if !isProxified && smartCache != nil {
		smartCachedItem = smartCache.Get(smartCache.CreateKey(
			request.StripPort(req.URL.String()),
			req.Header.Get("Accept"),
			req.Header.Get("Accept-Language"),
		))

		if smartCachedItem != nil {
			resp = &http.Response{
				Request: req,
				StatusCode: http.StatusOK,
				Status: http.StatusText(http.StatusOK),
				Header: smartCachedItem.Header,
				Body: io.NopCloser(bytes.NewReader(smartCachedItem.Body)),
			}
		}
	}

	if isRequestProxified(customCtx) {
		req.Header.Set(request.CacheControlHeaderKey, "no-store, no-cache, must-revalidate")
		req.Header.Set(request.PragmaHeaderKey, "no-cache")
	}

	shouldDispatch := options.TrafficDisplay == mitm.TrafficAll || (
		options.TrafficDisplay == mitm.TrafficOverrides && isProxified ||
		options.TrafficDisplay == mitm.TrafficSmartCache && !isResponseProxified(customCtx) && smartCache != nil)

	if shouldDispatch {
		var bodyBytes []byte
		if req.Body != nil {
			bodyBytes, _ = io.ReadAll(req.Body)
		}

		headersMap := request.HeadersToMap(req.Header)
		details := eventsService.StringifyRequestEventData(eventsService.ProxyRequestEventData{
			ID: uuid,
			URL: req.URL.String(),
			Method: req.Method,
			Headers: headersMap,
			Body: bodyBytes,
			Proxified: isProxified,
			SmartCached: !isProxified && smartCache != nil,
		})

		event.MustFire(eventsService.ProxyRequestSent, event.M{"details": details})
		req.Body = io.NopCloser(bytes.NewBuffer(bodyBytes))
	}

	return req, resp
}
