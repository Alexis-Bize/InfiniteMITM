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

package MITMApplicationMITMServiceContextModule

import "github.com/elazarl/goproxy"

type dataKey string

const (
	IDKey    dataKey = "uuid"
	ProxyKey dataKey = "proxified"
	CacheKey dataKey = "cache"
)

type CustomProxyCtx struct {
	*goproxy.ProxyCtx
	UserDataMap map[dataKey]interface{}
}

func ContextHandler(ctx *goproxy.ProxyCtx) *CustomProxyCtx {
	if customCtx, ok := ctx.UserData.(*CustomProxyCtx); ok {
		return customCtx
	}

	customCtx := &CustomProxyCtx{
		ProxyCtx:    ctx,
		UserDataMap: make(map[dataKey]interface{}),
	}

	ctx.UserData = customCtx
	return customCtx
}

func (c *CustomProxyCtx) SetUserData(key dataKey, value interface{}) {
	c.UserDataMap[key] = value
}

func (c *CustomProxyCtx) UnsetUserData(key dataKey) {
	delete(c.UserDataMap, key)
}

func (c *CustomProxyCtx) GetUserData(key dataKey) interface{} {
	return c.UserDataMap[key]
}

func (c *CustomProxyCtx) Reflect() interface{} {
	return c.UserDataMap
}
