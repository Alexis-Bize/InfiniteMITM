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
	"net/http"

	context "infinite-mitm/internal/application/services/mitm/modules/context"
	"infinite-mitm/pkg/smartcache"

	"github.com/elazarl/goproxy"
)

type RequestHandlerStruct struct {
	Match goproxy.ReqConditionFunc
	Fn    func(req *http.Request, ctx *goproxy.ProxyCtx) (*http.Request, *http.Response)
}

type ResponseHandlerStruct struct {
	Match goproxy.RespConditionFunc
	Fn    func(resp *http.Response, ctx *goproxy.ProxyCtx) *http.Response
}

func getUUID(customCtx *context.CustomProxyCtx) string {
	uuid := customCtx.GetUserData(context.IDKey).(string)
	return uuid
}

func getSmartCache(customCtx *context.CustomProxyCtx) *smartcache.SmartCache {
	var smartCache *smartcache.SmartCache

	cacheCtx := customCtx.GetUserData(context.CacheKey)
	if cacheCtx != nil {
		smartCache = cacheCtx.(*smartcache.SmartCache)
	}

	return smartCache
}

func isRequestProxified(customCtx *context.CustomProxyCtx) bool {
	proxified := customCtx.GetUserData(context.ProxyKey).(map[string]bool)
	return proxified["req"]
}

func isResponseProxified(customCtx *context.CustomProxyCtx) bool {
	proxified := customCtx.GetUserData(context.ProxyKey).(map[string]bool)
	return proxified["resp"]
}
