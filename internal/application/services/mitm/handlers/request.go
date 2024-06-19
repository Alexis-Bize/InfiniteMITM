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
	context "infinite-mitm/internal/application/services/mitm/modules/context"
	traffic "infinite-mitm/internal/application/services/mitm/modules/traffic"
	"infinite-mitm/pkg/request"
	"infinite-mitm/pkg/smartcache"
	"io"
	"net/http"

	"github.com/elazarl/goproxy"
	"github.com/gookit/event"
)

func HandleRequest(options traffic.TrafficOptions, req *http.Request, resp *http.Response, ctx *goproxy.ProxyCtx) (*http.Request, *http.Response) {
	if req.Method == http.MethodOptions {
		return req, resp
	}

	customCtx := context.ContextHandler(ctx)
	uuid := customCtx.GetUserData("uuid").(string)
	if uuid == "" {
		return req, resp
	}

	proxified := customCtx.GetUserData("proxified").(map[string]bool)
	isProxified := proxified["req"]

	if isProxified {
		req.Header.Set(request.CacheControlHeaderKey, "no-store, no-cache, must-revalidate")
		req.Header.Set(request.PragmaHeaderKey, "no-cache")
	}

	var smartCache *smartcache.SmartCache
	cacheCtx := customCtx.GetUserData("cache")
	if cacheCtx != nil {
		smartCache = cacheCtx.(*smartcache.SmartCache)
	}

	if options.TrafficDisplay == traffic.TrafficOverrides && !isProxified && smartCache == nil {
		return req, resp
	} else if options.TrafficDisplay == traffic.TrafficSmartCache && smartCache == nil {
		return req, resp
	}

	var bodyBytes []byte
	if req.Body != nil {
		bodyBytes, _ = io.ReadAll(req.Body)
	}

	req.Body = io.NopCloser(bytes.NewBuffer(bodyBytes))
	var smartCachedItem *smartcache.SmartCacheItem

	if smartCache != nil && !isProxified {
		smartCachedItem = smartCache.Get(smartCache.CreateKey(
			request.StripPort(req.URL.String()),
			req.Header.Get("Accept"),
			req.Header.Get("Accept-Language"),
		))
	}

	details := events.StringifyRequestEventData(events.ProxyRequestEventData{
		ID: uuid,
		URL: req.URL.String(),
		Method: req.Method,
		Headers: request.HeadersToMap(req.Header),
		Body: bodyBytes,
		Proxified: isProxified,
		SmartCached: !isProxified && (smartCache != nil || smartCachedItem != nil),
	})

	go func() {
		shouldDispatch := options.TrafficDisplay == traffic.TrafficAll || (
			options.TrafficDisplay == traffic.TrafficOverrides && isProxified ||
			options.TrafficDisplay == traffic.TrafficSmartCache && smartCache != nil)

		if shouldDispatch {
			event.MustFire(events.ProxyRequestSent, event.M{"details": details})
		}
	}()

	if smartCachedItem != nil {
		resp := goproxy.NewResponse(
			req,
			smartCachedItem.Header.Get(request.ContentTypeHeaderKey),
			http.StatusOK,
			string(smartCachedItem.Body),
		)

		for k, v := range request.HeadersToMap(smartCachedItem.Header) {
			resp.Header.Set(k, v)
		}
	}

	return req, resp
}
