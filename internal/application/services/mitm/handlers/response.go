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

func HandleResponse(options mitm.TrafficOptions, resp *http.Response, ctx *goproxy.ProxyCtx) *http.Response  {
	if resp.Request.Method == http.MethodOptions {
		return resp
	}

	customCtx := context.ContextHandler(ctx)
	uuid := customCtx.GetUserData("uuid").(string)
	if uuid == "" {
		return resp
	}

	proxified := customCtx.GetUserData("proxified").(map[string]bool)
	isProxified := proxified["resp"]

	if isProxified {
		resp.Header.Set(request.MITMProxyHeaderKey, request.MITMProxyEnabledValue)
		resp.Header.Set(request.CacheControlHeaderKey, "no-store, no-cache, must-revalidate, max-age=0")
		resp.Header.Set(request.PragmaHeaderKey, "no-cache")
		resp.Header.Set(request.ExpiresHeaderKey, "0")
	}

	var smartCache *smartcache.SmartCache
	cacheCtx := customCtx.GetUserData("cache")
	if cacheCtx != nil {
		smartCache = cacheCtx.(*smartcache.SmartCache)
	}

	if options.TrafficDisplay == mitm.TrafficOverrides && !isProxified && smartCache == nil {
		return resp
	} else if options.TrafficDisplay == mitm.TrafficSmartCache && smartCache == nil {
		return resp
	}

	var bodyBytes []byte
	var smartCachedItem *smartcache.SmartCacheItem
	var smartCacheKey string

	if smartCache != nil && !isProxified {
		smartCacheKey = smartCache.CreateKey(
			request.StripPort(resp.Request.URL.String()),
			resp.Request.Header.Get("Accept"),
			resp.Request.Header.Get("Accept-Language"),
		)

		smartCachedItem = smartCache.Get(smartCacheKey)
	}

	if smartCachedItem != nil {
		bodyBytes = smartCachedItem.Body
	} else {
		bodyBytes, _ = io.ReadAll(resp.Body)
	}

	if smartCache != nil && !isProxified {
		if smartCachedItem == nil {
			isSmartCachable := resp.StatusCode >= 200 && resp.StatusCode < 300
			if isSmartCachable {
				resp.Header.Set(request.MITMCacheHeaderKey, request.MITMCacheHeaderMissValue)
				smartCache.Write(
					smartCacheKey,
					&smartcache.SmartCacheItem{
						Body: bodyBytes,
						Header: resp.Header,
					},
				)
			}
		} else {
			resp.Header.Set(request.MITMCacheHeaderKey, request.MITMCacheHeaderHitValue)
		}
	}

	shouldDispatch := options.TrafficDisplay == mitm.TrafficAll || (
		options.TrafficDisplay == mitm.TrafficOverrides && isProxified ||
		options.TrafficDisplay == mitm.TrafficSmartCache && smartCache != nil)

	if shouldDispatch {
		headersMap := request.HeadersToMap(resp.Header)
		go func() {
			details := eventsService.StringifyResponseEventData(eventsService.ProxyResponseEventData{
				ID: uuid,
				URL: resp.Request.URL.String(),
				Method: resp.Request.Method,
				Status: resp.StatusCode,
				Headers: headersMap,
				Body: bodyBytes,
				Proxified: isProxified,
				SmartCached: !isProxified && smartCache != nil,
			})

			event.MustFire(eventsService.ProxyResponseReceived, event.M{"details": details})
		}()
	}

	resp.Body = io.NopCloser(bytes.NewBuffer(bodyBytes))
	return resp
}
