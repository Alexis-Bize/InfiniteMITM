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
	"infinite-mitm/pkg/domains"
	"infinite-mitm/pkg/request"
	"infinite-mitm/pkg/smartcache"
	"io"
	"net/http"
	"regexp"

	"github.com/elazarl/goproxy"
	"github.com/gookit/event"
)

func HandleRootResponses(options traffic.TrafficOptions) ResponseHandlerStruct {
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

			if options.TrafficDisplay == traffic.TrafficOverrides && !isProxified && smartCache == nil {
				return resp
			} else if options.TrafficDisplay == traffic.TrafficSmartCache && smartCache == nil {
				return resp
			}

			var bodyBytes []byte
			var smartCachedItem *smartcache.SmartCacheItem

			if smartCache != nil && !isProxified {
				smartCachedItem = smartCache.Read(smartCache.CreateKey(resp.Request.URL.String()))
			}

			if smartCachedItem != nil {
				bodyBytes = smartCachedItem.Body
			} else {
				bodyBytes, _ = io.ReadAll(resp.Body)
			}

			resp.Body = io.NopCloser(bytes.NewBuffer(bodyBytes))

			if smartCache != nil {
				isSmartCachable := !isProxified && resp.StatusCode >= 200 && resp.StatusCode < 300

				if smartCachedItem == nil {
					if isSmartCachable {
						resp.Header.Set(request.MITMCacheHeaderKey, request.MITMCacheHeaderMissValue)
						smartCache.Write(smartCache.CreateKey(resp.Request.URL.String()), &smartcache.SmartCacheItem{
							Body: bodyBytes,
							Header: resp.Header,
						})
					}
				} else {
					resp.Header.Set(request.MITMCacheHeaderKey, request.MITMCacheHeaderHitValue)
				}
			}

			shouldDispatch := options.TrafficDisplay == traffic.TrafficAll || (
				options.TrafficDisplay == traffic.TrafficOverrides && isProxified ||
				options.TrafficDisplay == traffic.TrafficSmartCache && smartCache != nil)

			if shouldDispatch {
				details := events.StringifyResponseEventData(events.ProxyResponseEventData{
					ID: uuid,
					URL: resp.Request.URL.String(),
					Method: resp.Request.Method,
					Status: resp.StatusCode,
					Headers: request.HeadersToMap(resp.Header),
					Body: bodyBytes,
					Proxified: isProxified,
					SmartCached: !isProxified && smartCache != nil,
				})

				event.MustFire(events.ProxyResponseReceived, event.M{"details": details})
			}

			return resp
		},
	}
}
